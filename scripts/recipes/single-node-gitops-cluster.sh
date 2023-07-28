#!/bin/bash
set -e

echo "-------------------------------------------"
echo "Installing zrun"
go build -o ./bin/zrun
cp ./bin/zrun /usr/local/bin/zrun
zrun about

echo "-------------------------------------------"
echo "Bootstrapping environment"

zrun_info() {
    zrun info disk
    zrun info mem
}

k3s_install() {
    sudo zrun k3s install \
       --tls-san "127.0.0.1" \
       --disable "traefik" || { echo "k3s installation failed"; exit 1; }
}

k9s_install() {
    sudo zrun k9s install || { echo "k9s installation failed"; exit 1; }
    k9s version
}

helm_install() {
    sudo zrun helm install-helm || { echo "Helm installation failed"; exit 1; }
    helm version
}

cert_manager_install() {
    sudo zrun certmanager install -v \
        || { echo "Cert-Manager installation failed"; exit 1; }
}

traefik_install() {
    sudo zrun traefik install \
        --defaults \
        --insecure \
        --proxy \
        --forwardedHeaders \
        --ingressProvider "cert-manager-resolver" \
        || { echo "Traefik installation failed"; exit 1; }
}

helm_install_argocd_chart() {
    sudo zrun helm install-chart \
        --repo-name "argo-cd" \
        --repo-url "https://argoproj.github.io/argo-helm" \
        --chart-name "argo-cd" \
        --namespace "argo-cd" \
        --chart-version "5.41.1" || { echo "ArgoCD installation failed"; exit 1; }
}

wait_for_cluster() {
    echo "-------------------------------------------"
    echo "Waiting for Cluster to be ready..."
    until kubectl get nodes; do sleep 1; done
}

wait_for_argocd() {
    echo "-------------------------------------------"
    echo "Waiting for ArgoCD to be ready..."
    until kubectl -n argo-cd get pods | grep Running; do sleep 1; done
}

run_k9s() {
    sudo zrun k9s || { echo "Running k9s failed"; exit 1; }
}

main() {
    zrun_info
    k3s_install
    k9s_install
    helm_install
    wait_for_cluster
    cert_manager_install
    traefik_install
    helm_install_argocd_chart
    wait_for_argocd
    run_k9s
}

main
