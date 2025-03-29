package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	s := server.NewMCPServer(
		"mcp-server-k8s-go",
		"0.0.0",
		server.WithLogging(),
	)

	version := mcp.NewTool(
		"kubernetes_version",

		mcp.WithDescription("Get Kubernetes version for current context"),
	)

	s.AddTool(version, versionHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func versionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config *rest.Config

	wd, _ := os.Getwd()
	debug := fmt.Sprintf("os user=%s, os wd=%s", os.Getenv("USER"), wd)

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig (%s): %v", debug, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %v", err)
	}

	return mcp.NewToolResultText(serverVersion.String()), nil
}
