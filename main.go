// K8sToolbox Golang Utility - Enhanced Implementation
// This utility provides Kubernetes-specific diagnostics, automated health checks, and connectivity tests.

package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Define flags for CLI commands
	healthCheckCmd := flag.NewFlagSet("healthcheck", flag.ExitOnError)
	namespace := healthCheckCmd.String("namespace", "default", "Namespace to check pod health")

	connectivityCheckCmd := flag.NewFlagSet("connectivity", flag.ExitOnError)
	namespaceConn := connectivityCheckCmd.String("namespace", "default", "Namespace of the pod")
	podName := connectivityCheckCmd.String("pod", "", "Name of the pod to test connectivity from")
	target := connectivityCheckCmd.String("target", "", "Target service or IP to check connectivity to")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'healthcheck' or 'connectivity' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "healthcheck":
		err := healthCheckCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		performHealthCheck(*namespace)
	case "connectivity":
		err := connectivityCheckCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *podName == "" || *target == "" {
			fmt.Println("Please specify both pod name and target for connectivity check")
			os.Exit(1)
		}
		testPodConnectivity(*namespaceConn, *podName, *target)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

// performHealthCheck performs a health check on all pods in the specified namespace
func performHealthCheck(namespace string) {
	// Load Kubernetes client configuration
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Get all pods in the specified namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		os.Exit(1)
	}

	// Perform health checks on each pod
	fmt.Printf("Performing health checks on namespace '%s'\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("Checking pod: %s\n", pod.Name)
		if pod.Status.Phase == "Running" {
			fmt.Printf("Pod %s is healthy\n", pod.Name)
		} else {
			fmt.Printf("Pod %s is not healthy (Status: %s)\n", pod.Name, pod.Status.Phase)
		}
	}
}

// testPodConnectivity tests network connectivity from a pod to a specified target
func testPodConnectivity(namespace, podName, target string) {
	// Load Kubernetes client configuration
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Set up the exec request
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", podName).
		Param("command", "ping").
		Param("command", "-c").
		Param("command", "3").
		Param("command", target).
		Param("stdin", "true").
		Param("stderr", "true").
		Param("stdout", "true").
		Param("tty", "false")

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Fatalf("Could not initialize command: %v", err)
	}

	// Create a context that can be used for managing cancellation and timeout
	ctx := context.Background()

	// Call StreamWithContext instead of the deprecated Stream
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		log.Fatalf("Could not execute command: %v", err)
	}
}
