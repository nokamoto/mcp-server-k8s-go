package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	list := mcp.NewTool(
		"kubernetes_list_resources",
		mcp.WithDescription(("List resources in Kubernetes for current context")),
		mcp.WithString(
			"resource",
			mcp.Required(),
			mcp.Description("The resource type to list (e.g. pods, services)"),
		),
		mcp.WithString(
			"group",
			mcp.Description("The API group of the resource (e.g. apps)"),
		),
		mcp.WithString(
			"version",
			mcp.Description("The API version of the resource (e.g. v1)"),
		),
		mcp.WithString(
			"label_selector",
			mcp.Description("A selector to restrict the list of returned objects by their labels. Defaults to everything."),
		),
		mcp.WithString(
			"field_selector",
			mcp.Description("A selector to restrict the list of returned objects by their fields. Defaults to everything."),
		),
		mcp.WithString(
			"namespace",
			mcp.Description("The namespace to list resources in. Defaults to all namespaces."),
		),
	)

	s.AddTool(version, versionHandler)
	s.AddTool(list, dynamicListHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func newConfig() (*rest.Config, error) {
	var config *rest.Config

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %v", err)
	}

	return config, nil
}

func versionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %v", err)
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

func dynamicListHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %v", err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	ps := func(name string) string {
		v, ok := request.Params.Arguments[name]
		if !ok {
			return ""
		}
		return v.(string)
	}

	schema := schema.GroupVersionResource{
		Group:    ps("group"),
		Version:  ps("version"),
		Resource: ps("resource"),
	}
	opts := metav1.ListOptions{
		LabelSelector: ps("label_selector"),
		FieldSelector: ps("field_selector"),
	}

	res, err := client.Resource(schema).Namespace(ps("namespace")).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %v", err)
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %v", err)
	}
	return mcp.NewToolResultText(string(bytes)), nil
}
