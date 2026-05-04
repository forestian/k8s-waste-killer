# k8s-waste-killer

A CLI tool that finds over-provisioned Kubernetes pods by comparing resource requests to actual usage from metrics-server, then suggests right-sized requests.

## Requirements

- Go 1.22+
- A running Kubernetes cluster with [metrics-server](https://github.com/kubernetes-sigs/metrics-server) installed
- `~/.kube/config` or `$KUBECONFIG` set

## Install

```bash
git clone https://github.com/k8s-waste-killer
cd k8s-waste-killer
go mod tidy
go build -o k8s-waste-killer .
```

## Install from GitHub Releases

Download a prebuilt binary from the [GitHub Releases](https://github.com/forestian/k8s-waste-killer/releases) page.

**Linux / macOS**
```bash
tar -xzf k8s-waste-killer_<version>_linux_amd64.tar.gz
chmod +x k8s-waste-killer
./k8s-waste-killer --help
```

**Windows**
```
# Extract the archive and run:
k8s-waste-killer.exe --help
```

## Usage

```bash
# Scan all namespaces
./k8s-waste-killer scan

# Scan a specific namespace
./k8s-waste-killer scan --namespace default

# Show only the top 10 most wasteful pods
./k8s-waste-killer scan --top 10

# JSON output
./k8s-waste-killer scan --output json

# Combine flags
./k8s-waste-killer scan --namespace monitoring --top 5 --output json

# Use a custom kubeconfig
./k8s-waste-killer scan --kubeconfig /path/to/kubeconfig
```

## Output example

```
NAMESPACE   POD              CPU_REQ   CPU_USE   CPU_WASTE   CPU_RECOMMEND   MEM_REQ   MEM_USE   MEM_WASTE   MEM_RECOMMEND
default     api-server       500m      120m      76%         180m            512Mi     180Mi     64%         270Mi
monitoring  loki-gateway     300m      80m       73%         120m            256Mi     90Mi      65%         135Mi
```

## How waste is calculated

```
waste_rate       = (request - usage) / request * 100
recommended_req  = usage * 1.5
```

- If usage > request, waste is reported as 0% (you are under-provisioned).
- If a pod has no resource request set, it is listed as `none`.
- If metrics-server is unavailable, usage columns show `no metrics`.
- Only `Running` pods are included.

## JSON output example

```json
[
  {
    "namespace": "default",
    "pod": "api-server-abc123",
    "cpu_request": "500m",
    "cpu_usage": "120m",
    "cpu_waste": "76%",
    "cpu_recommended": "180m",
    "mem_request": "512Mi",
    "mem_usage": "180Mi",
    "mem_waste": "64%",
    "mem_recommended": "270Mi"
  }
]
```

## Development

```bash
# Run directly without building
go run . scan --namespace default --top 10

# Run tests (once added)
go test ./...
```
