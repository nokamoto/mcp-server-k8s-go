# mcp-server-k8s-go
Sample MCP server for Kubernetes, written in Go

## codespace & cline setup

Run the following command to build the MCP server:

```bash
mage
```

The binary will be installed at `/go/bin/mcp-server-k8s-go`.

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
