# zrun

[![lint](https://github.com/zcubbs/zrun/actions/workflows/lint.yml/badge.svg)](https://github.com/zcubbs/zrun/actions/workflows/lint.yml)
[![release](https://github.com/zcubbs/zrun/actions/workflows/release.yml/badge.svg)](https://github.com/zcubbs/zrun/actions/workflows/release.yml)
[![vulnerability-scan](https://github.com/zcubbs/zrun/actions/workflows/vulnerability-scan.yml/badge.svg)](https://github.com/zcubbs/zrun/actions/workflows/vulnerability-scan.yml)

---
<p align="center">
  _ _ . .  . _ .  . . _  _ .
</p>
<p align="center">
  <img width="750" src="_assets/zrun_alt.jpg">
</p>

---

**zrun** is a comprehensive command-line interface (CLI) that provides a range of functionalities from installing k3s, managing Helm Deployments & Argocd applications to Git operations and more...

## Install

```bash
curl -sfL https://raw.githubusercontent.com/zcubbs/zrun/main/scripts/install/install.sh | bash
```

## Available Commands

- `about` - Print the info about ZRun
- `argo` - ArgoCD Commands
- `awx` - AWX Management Commands
- `certmanager` - cert-manager Commands
- `completion` - Generate the autocompletion script for the specified shell
- `config` - List CLI configuration
- `git` - Git Helper Commands
- `hello` - Used for test purposes
- `helm` - Helm Application Deployment/Control Commands
- `help` - Help about any command
- `info` - A palette that contains system info commands
- `k3s` - K3s management commands
- `k8s` - Kubernetes Management Commands
- `k9s` - Run K9s
- `os` - OS Helper Commands
- `traefik` - traefik Commands
- `upgrade` - Used to upgrade ZRun to the latest version
- `vault` - Vault Helper Commands
- `version` - Print the version number of ZRun

### Global Flags

- `-h, --help` - Help for this command
- `-v, --verbose` - Verbose output

### Sub-Commands 

- `k3s` - K3s management commands
  - **Usage**: `k3s [flags]` or `k3s [command]`
  - **Available Sub-Commands**:
    - `install` - Install k3s
    - `uninstall` - Uninstall k3s

- `argo` - ArgoCD Commands
  - **Usage**: `argo [flags]` or `argo [command]`
  - **Available Sub-Commands**:
    - `add-application` - Manage ArgoCD Applications
    - `add-project` - Add project to ArgoCD
    - `add-repository` - Add repository to ArgoCD
    - `install` - Install argo-cd Chart
  - **Flags**:
    - `-n, --namespace string` - Namespace (default "argo-cd")

- `certmanager` - cert-manager Commands
  - **Usage**: `certmanager [flags]` or `certmanager [command]`
  - **Available Sub-Commands**:
    - `install` - Install cert-manager Chart
    - `issuer` - Setup cert-manager issuer

- `traefik` - traefik Commands
  - **Usage**: `traefik [flags]` or `traefik [command]`
  - **Available Sub-Commands**:
    - `ingressroute` - Create IngressRoute
    - `install` - Install traefik Chart

- `helm` - Helm Application Deployment/Control Commands
  - **Usage**: `helm [flags]` or `helm [command]`
  - **Available Sub-Commands**:
    - `install-chart` - List all helm releases
    - `install-helm` - Install helm CLI
    - `list` - List all helm releases
    - `uninstall-chart` - List all helm releases

## Contributing

If you want to contribute to this project, please read the [contributing guidelines](CONTRIBUTING.md).

### Pre-requisites

- [Go](https://go.dev/doc/install) 1.20 or later
- Docker: (depending on OS)
  - [Docker](https://docs.docker.com/get-docker/) 20.10 or later
  - [Docker Compose](https://docs.docker.com/compose/install/) 1.29 or later
  - [Rancher Desktop](https://rancherdesktop.io/) 0.6.0 or later
- [GNU Make](https://www.gnu.org/software/make/) 4.3 or later
- [GNU Bash](https://www.gnu.org/software/bash/) 5.1 or later
- [GoSec](https://github.com/securego/gosec)
  - `go install github.com/securego/gosec/v2/cmd/gosec@latest`
- [golangci-lint](https://golangci-lint.run/usage/install/) 1.42 or later 
  - `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- [Cobra CLI](https://github.com/spf13/cobra-cli) 1.2 or later
  - `go install github.com/spf13/cobra-cli@latest`
- Virtual Machines:
  - Windows
    - [VirtualBox](https://www.virtualbox.org/wiki/Downloads) 6.1 or later
    - [Vagrant](https://www.vagrantup.com/downloads) 2.2 or later
  - MacOS
    - [Lima](https://github.com/lima-vm/lima) 0.6.0 or later
    - [nerdctl](https://github.com/containerd/nerdctl) 0.12.0 or later

### Tooling

#### Lint
```bash
make lint
```

#### Test
```bash
make test
```

#### Cobra
Extending the command line:
```bash
cobra-cli add <command>
```

#### Vagrant
Spin up a local Ubuntu VM (with a builtin zrun binary for testing:

```bash
make vssh
```

---
## License

[Apache 2.0](LICENSE)
