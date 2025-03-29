# mcp-server-k8s-go
Sample MCP server for Kubernetes, written in Go

## codespace setup

Run the following command to start minikube.

```bash
minikube start
```

Optionally, create a sample deployment to test and verify the MCP server.

```bash
kubectl create deployment kubernetes-bootcamp --image=gcr.io/google-samples/kubernetes-bootcamp:v1
```

Run the following command to build the MCP server:

```bash
go install  github.com/magefile/mage@latest
mage
```

## testing MCP server in shell

Run the following command to communicate with the MCP server in the [studio](https://spec.modelcontextprotocol.io/specification/2025-03-26/basic/transports/#stdio) transport, which allows interaction through standard input and output.

```bash
mcp-server-k8s-go < <(ls examples | xargs -I{} bash -c 'echo $(cat examples/{})')
```

## cline setup

Add the following configuration to set up MCP servers:

```json
{
    "mcpServers": {
        "local-kubernetes": {
            "command": "/go/bin/mcp-server-k8s-go",
            "env": {
                "KUBECONFIG": "/home/codespace/.kube/config"
            }
        }
    }
}
```
