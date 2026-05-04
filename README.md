# k8s-waste-killer

Stop paying for CPU and memory your pods never use.

> Part of the [Forestian Cloud Native Toolkit](https://github.com/forestian) — small CLI tools for Kubernetes, observability, GitOps, and platform engineering.

## Demo

GIF demo coming soon.

## Quick Demo

```bash
k8s-waste-killer scan --top 5
```

```
NAMESPACE    POD                  CPU_REQ   CPU_USE   CPU_WASTE   CPU_RECOMMEND   MEM_REQ   MEM_USE   MEM_WASTE   MEM_RECOMMEND
default      api-server           500m      120m      76%         180m            512Mi     180Mi     64%         270Mi
monitoring   loki-gateway         300m      80m       73%         120m            256Mi     90Mi      65%         135Mi
kube-system  metrics-server       100m      12m       88%         18m             128Mi     32Mi      75%         48Mi
monitoring   prometheus-0         1000m     310m      69%         465m            2048Mi    820Mi     59%         1230Mi
logging      fluent-bit-7xk9p     200m      25m       87%         37m             128Mi     28Mi      78%         42Mi
```

## Quick Start

**1. Download a prebuilt binary**

Go to the [Releases page](https://github.com/forestian/k8s-waste-killer/releases) and download the archive for your OS and architecture.

```bash
# Linux / macOS
tar -xzf k8s-waste-killer_<version>_linux_amd64.tar.gz
chmod +x k8s-waste-killer
./k8s-waste-killer scan
```

```
# Windows — extract the archive and run:
k8s-waste-killer.exe scan
```

**2. Or build from source (requires Go 1.22+)**

```bash
git clone https://github.com/forestian/k8s-waste-killer
cd k8s-waste-killer
go build -o k8s-waste-killer .
./k8s-waste-killer scan
```

**Requirements:** Kubernetes cluster with [metrics-server](https://github.com/kubernetes-sigs/metrics-server) installed · `~/.kube/config` or `$KUBECONFIG` set

## Use Cases

- Audit CPU and memory requests before a cloud cost review
- Find the worst-offending pods in a specific namespace
- Right-size requests after a new service stabilises in production
- Pipe `--output json` into scripts or CI dashboards
- Spot over-provisioned workloads before cluster autoscaling kicks in

## Usage

```bash
k8s-waste-killer scan                                       # all namespaces
k8s-waste-killer scan -n default                            # one namespace
k8s-waste-killer scan --top 10                              # top 10 by waste
k8s-waste-killer scan -o json                               # JSON output
k8s-waste-killer scan -n monitoring --top 5 -o json         # combine flags
k8s-waste-killer scan --kubeconfig /path/to/kubeconfig      # custom kubeconfig
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--namespace` | `-n` | all | Kubernetes namespace to scan |
| `--top` | | 0 (all) | Show only top N pods by waste |
| `--output` | `-o` | `table` | Output format: `table` or `json` |
| `--kubeconfig` | | `~/.kube/config` | Path to kubeconfig file |

## Example Output

**Table (default)**

```
NAMESPACE   POD              CPU_REQ   CPU_USE   CPU_WASTE   CPU_RECOMMEND   MEM_REQ   MEM_USE   MEM_WASTE   MEM_RECOMMEND
default     api-server       500m      120m      76%         180m            512Mi     180Mi     64%         270Mi
monitoring  loki-gateway     300m      80m       73%         120m            256Mi     90Mi      65%         135Mi
```

**JSON (`-o json`)**

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

## How Waste is Calculated

```
waste_rate      = (request - usage) / request × 100
recommended_req = usage × 1.5
```

- Usage > request → waste reported as 0% (under-provisioned).
- No resource request set → listed as `none`.
- metrics-server unavailable → usage columns show `no metrics`.
- Only `Running` pods are included.

## Limitations

- **Read-only** — does not modify any cluster resource.
- Requires [metrics-server](https://github.com/kubernetes-sigs/metrics-server); without it, usage shows `no metrics`.
- Usage is a point-in-time snapshot from metrics-server, not a historical average.
- Only `Running` pods are scanned; Pending and Completed pods are skipped.
- Recommendations (`usage × 1.5`) are suggestions — review before applying to production.

## Roadmap

- `--min-waste` flag to filter by minimum waste percentage
- Namespace-level aggregated summary
- GitHub Actions integration for CI cost reports
- Unit and integration tests
- VPA recommendation comparison

## Development

```bash
go run . scan --namespace default --top 10   # run without building
go test ./...
```

---

Part of the [Forestian Cloud Native Toolkit](https://github.com/forestian) — small CLI tools for Kubernetes, observability, GitOps, and platform engineering.
