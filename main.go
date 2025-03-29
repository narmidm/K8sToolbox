// K8sToolbox Golang Utility - Enhanced Implementation
// This utility provides Kubernetes-specific diagnostics, automated health checks, and connectivity tests.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// Version information - should be set during build
var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	Commit    = "unknown"
)

// Global variables
var (
	clientset *kubernetes.Clientset
	config    *rest.Config
	logger    *log.Logger

	// Prometheus metrics
	checksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8stoolbox_checks_total",
			Help: "Total number of health checks performed",
		},
		[]string{"namespace", "status"},
	)

	connectivityChecksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8stoolbox_connectivity_checks_total",
			Help: "Total number of connectivity checks performed",
		},
		[]string{"namespace", "status", "target"},
	)

	resourceUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8stoolbox_resource_usage",
			Help: "Resource usage metrics collected by K8sToolbox",
		},
		[]string{"namespace", "pod", "resource_type"},
	)
)

// Configuration holds all app configuration
type Configuration struct {
	// Web UI configuration
	EnableWebUI bool
	WebUIPort   int
	WebUIPath   string

	// Authentication configuration
	EnableAuth     bool
	AuthUsername   string
	AuthPassword   string
	AuthSecretName string

	// Monitoring configuration
	EnablePrometheus bool
	PrometheusPort   int

	// Logging configuration
	LogLevel  string
	LogFormat string

	// Kubernetes client configuration
	KubeConfig     string
	DefaultTimeout time.Duration
}

// PodHealthStatus represents the health status of a pod
type PodHealthStatus struct {
	Name   string
	Status string
	Issues []string
}

// APIResponse defines the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Global configuration
var config2 = Configuration{
	EnableWebUI:      getBoolEnv("ENABLE_WEB_UI", false),
	WebUIPort:        getIntEnv("WEB_UI_PORT", 8080),
	WebUIPath:        getEnv("WEB_UI_PATH", "/"),
	EnableAuth:       getBoolEnv("ENABLE_AUTH", true),
	AuthUsername:     getEnv("AUTH_USERNAME", "admin"),
	AuthPassword:     getEnv("AUTH_PASSWORD", ""),
	AuthSecretName:   getEnv("AUTH_SECRET_NAME", "k8stoolbox-auth"),
	EnablePrometheus: getBoolEnv("ENABLE_PROMETHEUS", false),
	PrometheusPort:   getIntEnv("PROMETHEUS_PORT", 9090),
	LogLevel:         getEnv("LOG_LEVEL", "info"),
	LogFormat:        getEnv("LOG_FORMAT", "text"),
	KubeConfig:       getEnv("KUBECONFIG", ""),
	DefaultTimeout:   time.Duration(getIntEnv("DEFAULT_TIMEOUT", 30)) * time.Second,
}

// StandaloneMode allows running without Kubernetes
var StandaloneMode = getBoolEnv("STANDALONE_MODE", false)

func init() {
	// Initialize logger
	logger = log.New(os.Stdout, "[K8sToolbox] ", log.LstdFlags|log.Lshortfile)

	// Register Prometheus metrics
	prometheus.MustRegister(checksTotal)
	prometheus.MustRegister(connectivityChecksTotal)
	prometheus.MustRegister(resourceUsage)
}

func main() {
	// Display version information
	logger.Printf("K8sToolbox %s (Build: %s, Commit: %s)\n", Version, BuildTime, Commit)

	// Create context that can be canceled on SIGTERM/SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Printf("Received signal: %v, shutting down gracefully", sig)
		cancel()
	}()

	// Initialize Kubernetes client
	if !StandaloneMode {
		if err := initKubernetesClient(); err != nil {
			logger.Fatalf("Failed to initialize Kubernetes client: %v", err)
		}
	} else {
		logger.Println("Running in standalone mode - Kubernetes client not initialized")
	}

	// Start web server if enabled
	if config2.EnableWebUI {
		go startWebServer(ctx)
	}

	// Start Prometheus server if enabled
	if config2.EnablePrometheus {
		go startPrometheusServer(ctx)
	}

	// Define flags for CLI commands
	healthCheckCmd := flag.NewFlagSet("healthcheck", flag.ExitOnError)
	namespace := healthCheckCmd.String("namespace", "default", "Namespace to check pod health")
	timeout := healthCheckCmd.Duration("timeout", 30*time.Second, "Timeout for operations")

	connectivityCheckCmd := flag.NewFlagSet("connectivity", flag.ExitOnError)
	namespaceConn := connectivityCheckCmd.String("namespace", "default", "Namespace of the pod")
	podName := connectivityCheckCmd.String("pod", "", "Name of the pod to test connectivity from")
	target := connectivityCheckCmd.String("target", "", "Target service or IP to check connectivity to")
	protocol := connectivityCheckCmd.String("protocol", "tcp", "Protocol to use (tcp/http/icmp)")
	port := connectivityCheckCmd.Int("port", 80, "Port to connect to for TCP/HTTP checks")

	resourceCheckCmd := flag.NewFlagSet("resources", flag.ExitOnError)
	namespaceRes := resourceCheckCmd.String("namespace", "default", "Namespace to check resources")
	threshold := resourceCheckCmd.Int("threshold", 80, "Resource usage threshold percentage for warnings")

	// New monitoring command
	monitorCmd := flag.NewFlagSet("monitor", flag.ExitOnError)
	monitorNamespace := monitorCmd.String("namespace", "default", "Namespace to monitor")
	interval := monitorCmd.Duration("interval", 30*time.Second, "Monitoring interval")
	output := monitorCmd.String("output", "stdout", "Output destination (stdout, prometheus, json)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Create context with timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, *timeout)
	defer timeoutCancel()

	switch os.Args[1] {
	case "healthcheck":
		err := healthCheckCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		performHealthCheck(timeoutCtx, *namespace)
	case "connectivity":
		err := connectivityCheckCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *podName == "" || *target == "" {
			logger.Println("Please specify both pod name and target for connectivity check")
			os.Exit(1)
		}
		testPodConnectivity(timeoutCtx, *namespaceConn, *podName, *target, *protocol, *port)
	case "resources":
		err := resourceCheckCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		checkResourceUsage(timeoutCtx, *namespaceRes, *threshold)
	case "monitor":
		err := monitorCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		startMonitoring(ctx, *monitorNamespace, *interval, *output)
	case "version":
		fmt.Printf("K8sToolbox %s\nBuild Time: %s\nCommit: %s\n", Version, BuildTime, Commit)
	case "server":
		// Just start the servers and wait
		logger.Println("Starting in server mode...")
		<-ctx.Done()
	default:
		logger.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	// Wait for any background goroutines to complete
	if config2.EnableWebUI || config2.EnablePrometheus {
		select {
		case <-ctx.Done():
			logger.Println("Shutting down...")
			// Add a small delay to allow servers to shut down gracefully
			time.Sleep(1 * time.Second)
		}
	}
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get boolean environment variables
func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}

// Helper function to get integer environment variables
func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// basicAuth implements HTTP basic authentication middleware
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for prometheus metrics if needed
		if r.URL.Path == "/metrics" && !config2.EnableAuth {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || username != config2.AuthUsername || password != config2.AuthPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="K8sToolbox"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// startPrometheusServer starts the Prometheus metrics server
func startPrometheusServer(ctx context.Context) {
	addr := fmt.Sprintf(":%d", config2.PrometheusPort)
	logger.Printf("Starting Prometheus metrics server on %s", addr)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("Prometheus server error: %v", err)
		}
	}()

	<-ctx.Done()
	// Give the server a grace period to shut down
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Error shutting down Prometheus server: %v", err)
	} else {
		logger.Println("Prometheus server shut down successfully")
	}
}

// startWebServer starts the Web UI server
func startWebServer(ctx context.Context) {
	addr := fmt.Sprintf(":%d", config2.WebUIPort)
	logger.Printf("Starting Web UI server on %s", addr)

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/v1/health", healthHandler)
	mux.HandleFunc("/api/v1/namespaces", namespacesHandler)
	mux.HandleFunc("/api/v1/pods", podsHandler)
	mux.HandleFunc("/api/v1/services", servicesHandler)
	mux.HandleFunc("/api/v1/nodes", nodesHandler)

	// Static content (in a real implementation this would serve actual HTML/JS/CSS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "K8sToolbox Web UI - Version %s", Version)
			return
		}
		http.NotFound(w, r)
	})

	var handler http.Handler = mux
	if config2.EnableAuth {
		handler = basicAuth(mux)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("Web server error: %v", err)
		}
	}()

	<-ctx.Done()
	// Give the server a grace period to shut down
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Error shutting down web server: %v", err)
	} else {
		logger.Println("Web server shut down successfully")
	}
}

// API Handlers
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Success: true,
		Message: "K8sToolbox is running",
		Data: map[string]string{
			"version":    Version,
			"buildTime":  BuildTime,
			"commitHash": Commit,
			"status":     "healthy",
			"mode":       conditionalString(StandaloneMode, "standalone", "kubernetes"),
		},
	}
	json.NewEncoder(w).Encode(response)
}

func namespacesHandler(w http.ResponseWriter, r *http.Request) {
	var namespaceNames []string

	if StandaloneMode {
		// In standalone mode, return some dummy namespaces
		namespaceNames = []string{"default", "kube-system", "demo"}
	} else {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, ns := range namespaces.Items {
			namespaceNames = append(namespaceNames, ns.Name)
		}
	}

	response := APIResponse{
		Success: true,
		Data:    namespaceNames,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func podsHandler(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// Define the PodInfo type here for clarity
	type PodInfo struct {
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
		Status    string            `json:"status"`
		Ready     bool              `json:"ready"`
		Labels    map[string]string `json:"labels"`
	}

	var podInfos []PodInfo

	if StandaloneMode {
		// In standalone mode, return some dummy pods
		podInfos = []PodInfo{
			{
				Name:      "example-pod-1",
				Namespace: namespace,
				Status:    "Running",
				Ready:     true,
				Labels: map[string]string{
					"app":  "example",
					"tier": "frontend",
				},
			},
			{
				Name:      "example-pod-2",
				Namespace: namespace,
				Status:    "Running",
				Ready:     true,
				Labels: map[string]string{
					"app":  "example",
					"tier": "backend",
				},
			},
			{
				Name:      "example-pod-3",
				Namespace: namespace,
				Status:    "Pending",
				Ready:     false,
				Labels: map[string]string{
					"app":  "example",
					"tier": "database",
				},
			},
		}
	} else {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, pod := range pods.Items {
			ready := true
			for _, cs := range pod.Status.ContainerStatuses {
				if !cs.Ready {
					ready = false
					break
				}
			}

			podInfos = append(podInfos, PodInfo{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Status:    string(pod.Status.Phase),
				Ready:     ready,
				Labels:    pod.Labels,
			})
		}
	}

	response := APIResponse{
		Success: true,
		Data:    podInfos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// Define the ServiceInfo type here for clarity
	type ServiceInfo struct {
		Name       string            `json:"name"`
		Namespace  string            `json:"namespace"`
		Type       string            `json:"type"`
		ClusterIP  string            `json:"clusterIP"`
		ExternalIP []string          `json:"externalIP,omitempty"`
		Ports      []string          `json:"ports"`
		Labels     map[string]string `json:"labels"`
	}

	var serviceInfos []ServiceInfo

	if StandaloneMode {
		// In standalone mode, return some dummy services
		serviceInfos = []ServiceInfo{
			{
				Name:      "example-service-1",
				Namespace: namespace,
				Type:      "ClusterIP",
				ClusterIP: "10.96.0.10",
				Ports:     []string{"80/TCP", "443/TCP"},
				Labels: map[string]string{
					"app":  "example",
					"tier": "frontend",
				},
			},
			{
				Name:       "example-service-2",
				Namespace:  namespace,
				Type:       "LoadBalancer",
				ClusterIP:  "10.96.0.11",
				ExternalIP: []string{"192.168.1.100"},
				Ports:      []string{"8080/TCP"},
				Labels: map[string]string{
					"app":  "example",
					"tier": "backend",
				},
			},
		}
	} else {
		services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, svc := range services.Items {
			var ports []string
			for _, port := range svc.Spec.Ports {
				ports = append(ports, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
			}

			// Safely handle external IPs
			var externalIPs []string
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				for _, ingress := range svc.Status.LoadBalancer.Ingress {
					if ingress.IP != "" {
						externalIPs = append(externalIPs, ingress.IP)
					}
				}
			}

			serviceInfos = append(serviceInfos, ServiceInfo{
				Name:       svc.Name,
				Namespace:  svc.Namespace,
				Type:       string(svc.Spec.Type),
				ClusterIP:  svc.Spec.ClusterIP,
				ExternalIP: externalIPs,
				Ports:      ports,
				Labels:     svc.Labels,
			})
		}
	}

	response := APIResponse{
		Success: true,
		Data:    serviceInfos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	// Define the NodeInfo type here for clarity
	type NodeInfo struct {
		Name             string            `json:"name"`
		Status           string            `json:"status"`
		Addresses        []string          `json:"addresses"`
		KubeletVersion   string            `json:"kubeletVersion"`
		KernelVersion    string            `json:"kernelVersion"`
		OSImage          string            `json:"osImage"`
		ContainerRuntime string            `json:"containerRuntime"`
		Labels           map[string]string `json:"labels"`
	}

	var nodeInfos []NodeInfo

	if StandaloneMode {
		// In standalone mode, return some dummy nodes
		nodeInfos = []NodeInfo{
			{
				Name:             "example-node-1",
				Status:           "Ready",
				Addresses:        []string{"InternalIP: 192.168.1.10", "Hostname: example-node-1"},
				KubeletVersion:   "v1.28.3",
				KernelVersion:    "5.15.0-86-generic",
				OSImage:          "Ubuntu 22.04.3 LTS",
				ContainerRuntime: "containerd://1.7.4",
				Labels: map[string]string{
					"kubernetes.io/hostname":                "example-node-1",
					"node-role.kubernetes.io/control-plane": "",
				},
			},
			{
				Name:             "example-node-2",
				Status:           "Ready",
				Addresses:        []string{"InternalIP: 192.168.1.11", "Hostname: example-node-2"},
				KubeletVersion:   "v1.28.3",
				KernelVersion:    "5.15.0-86-generic",
				OSImage:          "Ubuntu 22.04.3 LTS",
				ContainerRuntime: "containerd://1.7.4",
				Labels: map[string]string{
					"kubernetes.io/hostname":         "example-node-2",
					"node-role.kubernetes.io/worker": "",
				},
			},
		}
	} else {
		nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, node := range nodes.Items {
			// Determine node status
			status := "Unknown"
			for _, cond := range node.Status.Conditions {
				if cond.Type == "Ready" {
					if cond.Status == "True" {
						status = "Ready"
					} else {
						status = "NotReady"
					}
					break
				}
			}

			// Collect addresses
			var addresses []string
			for _, addr := range node.Status.Addresses {
				addresses = append(addresses, fmt.Sprintf("%s: %s", addr.Type, addr.Address))
			}

			nodeInfos = append(nodeInfos, NodeInfo{
				Name:             node.Name,
				Status:           status,
				Addresses:        addresses,
				KubeletVersion:   node.Status.NodeInfo.KubeletVersion,
				KernelVersion:    node.Status.NodeInfo.KernelVersion,
				OSImage:          node.Status.NodeInfo.OSImage,
				ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
				Labels:           node.Labels,
			})
		}
	}

	response := APIResponse{
		Success: true,
		Data:    nodeInfos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Utility function to return error responses
func errorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := APIResponse{
		Success: false,
		Error:   errorMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// printUsage prints the usage instructions
func printUsage() {
	fmt.Println("K8sToolbox - Kubernetes Management and Troubleshooting Utility")
	fmt.Println("\nUsage:")
	fmt.Println("  k8stoolbox <command> [options]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  healthcheck    Performs health checks on pods in a namespace")
	fmt.Println("  connectivity   Tests network connectivity from a pod to a target")
	fmt.Println("  resources      Checks resource usage in a namespace")
	fmt.Println("  monitor        Continuously monitors resources with the specified interval")
	fmt.Println("  server         Starts the web UI and/or Prometheus metrics server")
	fmt.Println("  version        Shows version information")
	fmt.Println("\nUse 'k8stoolbox <command> --help' for more information about a command.")
	fmt.Println("\nEnvironment variables:")
	fmt.Println("  ENABLE_WEB_UI       Enable web UI (true/false)")
	fmt.Println("  WEB_UI_PORT         Web UI port (default: 8080)")
	fmt.Println("  ENABLE_AUTH         Enable authentication (true/false)")
	fmt.Println("  AUTH_USERNAME       Username for basic authentication")
	fmt.Println("  AUTH_PASSWORD       Password for basic authentication")
	fmt.Println("  ENABLE_PROMETHEUS   Enable Prometheus metrics endpoint (true/false)")
	fmt.Println("  PROMETHEUS_PORT     Prometheus metrics port (default: 9090)")
}

// initKubernetesClient initializes the Kubernetes client
func initKubernetesClient() error {
	var err error

	// Try to use in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		kubeconfig := config2.KubeConfig
		if kubeconfig == "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}

		logger.Printf("Using kubeconfig: %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes config: %v", err)
		}
	}

	// Set reasonable timeouts
	config.Timeout = 30 * time.Second

	// Creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %v", err)
	}

	return nil
}

// startMonitoring begins continuous monitoring of cluster resources
func startMonitoring(ctx context.Context, namespace string, interval time.Duration, outputFormat string) {
	logger.Printf("Starting monitoring of namespace '%s' with interval %v", namespace, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Println("Monitoring stopped")
			return
		case <-ticker.C:
			// Perform health check
			healthResult := performHealthCheckWithResults(ctx, namespace)

			// Output results based on format
			switch outputFormat {
			case "json":
				jsonData, err := json.Marshal(healthResult)
				if err == nil {
					fmt.Println(string(jsonData))
				}
			case "prometheus":
				// Update Prometheus metrics (already done in performHealthCheckWithResults)
				logger.Println("Metrics updated in Prometheus")
			default:
				// Default stdout output
				fmt.Printf("Health check at %s: %d healthy pods, %d unhealthy pods\n",
					time.Now().Format(time.RFC3339),
					healthResult.HealthyPods,
					healthResult.UnhealthyPods)
			}

			// Collect resource usage
			checkResourceUsage(ctx, namespace, 80)
		}
	}
}

// HealthCheckResult represents the result of a health check operation
type HealthCheckResult struct {
	Namespace     string            `json:"namespace"`
	HealthyPods   int               `json:"healthyPods"`
	UnhealthyPods int               `json:"unhealthyPods"`
	PodDetails    []PodHealthStatus `json:"podDetails"`
	Timestamp     time.Time         `json:"timestamp"`
}

// performHealthCheckWithResults performs a health check and returns structured results
func performHealthCheckWithResults(ctx context.Context, namespace string) HealthCheckResult {
	result := HealthCheckResult{
		Namespace:  namespace,
		Timestamp:  time.Now(),
		PodDetails: []PodHealthStatus{},
	}

	// Get all pods in the specified namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Printf("Error listing pods: %v\n", err)
		return result
	}

	if len(pods.Items) == 0 {
		logger.Printf("No pods found in namespace '%s'\n", namespace)
		return result
	}

	// Process each pod
	for _, pod := range pods.Items {
		podStatus := PodHealthStatus{
			Name:   pod.Name,
			Status: string(pod.Status.Phase),
			Issues: []string{},
		}

		// Check readiness and liveness probes
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				podStatus.Issues = append(podStatus.Issues,
					fmt.Sprintf("Container %s is not ready", containerStatus.Name))
			}

			if containerStatus.RestartCount > 5 {
				podStatus.Issues = append(podStatus.Issues,
					fmt.Sprintf("Container %s has restarted %d times",
						containerStatus.Name, containerStatus.RestartCount))
			}
		}

		// Check pod conditions
		for _, condition := range pod.Status.Conditions {
			if condition.Status != "True" && condition.Type != "PodScheduled" {
				podStatus.Issues = append(podStatus.Issues,
					fmt.Sprintf("Condition %s is %s: %s",
						condition.Type, condition.Status, condition.Message))
			}
		}

		// Update counts and details
		result.PodDetails = append(result.PodDetails, podStatus)

		if len(podStatus.Issues) > 0 || podStatus.Status != "Running" {
			result.UnhealthyPods++

			// Update Prometheus metrics
			checksTotal.WithLabelValues(namespace, "unhealthy").Inc()
		} else {
			result.HealthyPods++

			// Update Prometheus metrics
			checksTotal.WithLabelValues(namespace, "healthy").Inc()
		}
	}

	return result
}

// performHealthCheck performs a health check on all pods in the specified namespace
func performHealthCheck(ctx context.Context, namespace string) {
	result := performHealthCheckWithResults(ctx, namespace)

	// Log the results
	logger.Printf("Performing health checks on namespace '%s'\n", namespace)

	for _, podStatus := range result.PodDetails {
		if len(podStatus.Issues) > 0 || podStatus.Status != "Running" {
			logger.Printf("⚠️ Pod %s is not healthy (Status: %s)\n", podStatus.Name, podStatus.Status)
			for _, issue := range podStatus.Issues {
				logger.Printf("  - %s\n", issue)
			}
		} else {
			logger.Printf("✅ Pod %s is healthy\n", podStatus.Name)
		}
	}

	logger.Printf("Health check summary: %d healthy pods, %d unhealthy pods\n",
		result.HealthyPods, result.UnhealthyPods)
}

// testPodConnectivity tests network connectivity from a pod to a specified target
func testPodConnectivity(ctx context.Context, namespace, podName, target, protocol string, port int) {
	if clientset == nil || config == nil {
		logger.Fatalf("Error: Kubernetes clientset or config is not initialized")
	}

	// Validate protocol
	protocol = strings.ToLower(protocol)
	if protocol != "tcp" && protocol != "http" && protocol != "icmp" {
		logger.Fatalf("Invalid protocol: %s. Must be one of: tcp, http, icmp", protocol)
	}

	var command []string
	switch protocol {
	case "tcp":
		command = []string{"nc", "-zv", "-w", "5", target, fmt.Sprintf("%d", port)}
	case "http":
		command = []string{"curl", "-sSf", "-m", "10", "-o", "/dev/null", fmt.Sprintf("http://%s:%d", target, port)}
	case "icmp":
		command = []string{"ping", "-c", "3", target}
	}

	logger.Printf("Testing %s connectivity from pod %s to %s\n", protocol, podName, target)

	// Set up the exec request
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", ""). // Leave empty to use the first container
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	// Add command params
	for _, cmd := range command {
		req.Param("command", cmd)
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		logger.Fatalf("Could not initialize command: %v", err)
		connectivityChecksTotal.WithLabelValues(namespace, "error", target).Inc()
	}

	// Call StreamWithContext
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		logger.Fatalf("Connectivity test failed: %v", err)
		connectivityChecksTotal.WithLabelValues(namespace, "failed", target).Inc()
	} else {
		logger.Printf("Connectivity test succeeded\n")
		connectivityChecksTotal.WithLabelValues(namespace, "success", target).Inc()
	}
}

// checkResourceUsage checks the resource usage in a namespace
func checkResourceUsage(ctx context.Context, namespace string, threshold int) {
	logger.Printf("Checking resource usage in namespace: %s\n", namespace)

	// Get pods in the namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error listing pods: %v", err)
	}

	if len(pods.Items) == 0 {
		logger.Printf("No pods found in namespace '%s'\n", namespace)
		return
	}

	// Simple resource reporting for now - in a real implementation we'd use metrics-server
	// or prometheus for actual resource usage
	logger.Printf("Resource allocation in namespace '%s' (showing requested resources):\n", namespace)

	// Print headers
	fmt.Printf("%-40s %-10s %-10s %-10s %-10s\n", "POD", "CPU REQ", "CPU LIM", "MEM REQ", "MEM LIM")
	fmt.Println(strings.Repeat("-", 80))

	for _, pod := range pods.Items {
		// Calculate total requests and limits for the pod
		var cpuReq, cpuLim, memReq, memLim string

		cpuReq = "0"
		cpuLim = "0"
		memReq = "0"
		memLim = "0"

		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				if cpu, ok := container.Resources.Requests["cpu"]; ok {
					cpuReq = cpu.String()
					// Update Prometheus metrics
					cpuValue := cpu.AsApproximateFloat64()
					resourceUsage.WithLabelValues(namespace, pod.Name, "cpu_request").Set(cpuValue)
				}
				if mem, ok := container.Resources.Requests["memory"]; ok {
					memReq = mem.String()
					// Update Prometheus metrics
					memValue := mem.AsApproximateFloat64()
					resourceUsage.WithLabelValues(namespace, pod.Name, "memory_request").Set(memValue)
				}
			}

			if container.Resources.Limits != nil {
				if cpu, ok := container.Resources.Limits["cpu"]; ok {
					cpuLim = cpu.String()
					// Update Prometheus metrics
					cpuValue := cpu.AsApproximateFloat64()
					resourceUsage.WithLabelValues(namespace, pod.Name, "cpu_limit").Set(cpuValue)
				}
				if mem, ok := container.Resources.Limits["memory"]; ok {
					memLim = mem.String()
					// Update Prometheus metrics
					memValue := mem.AsApproximateFloat64()
					resourceUsage.WithLabelValues(namespace, pod.Name, "memory_limit").Set(memValue)
				}
			}
		}

		fmt.Printf("%-40s %-10s %-10s %-10s %-10s\n",
			pod.Name, cpuReq, cpuLim, memReq, memLim)
	}
}

// conditionalString returns the first string if condition is true, otherwise the second
func conditionalString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}
