#!/bin/sh
set -e

echo "-------------------------------------------"
echo "Installing zrun"
go build -o ./bin/zrun
cp ./bin/zrun /usr/local/bin/zrun
zrun version

echo "-------------------------------------------"
echo "Bootstrapping environment..."

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
}

helm_install() {
    sudo zrun helm install-helm || { echo "Helm installation failed"; exit 1; }
}

cert_manager_install() {
    sudo zrun certmanager install \
        || { echo "cert-manager installation failed"; exit 1; }
}

traefik_install() {
    sudo zrun traefik install \
        --defaults \
        --insecure \
        --proxy \
        --forwardedHeaders \
        --ingressProvider "cert-manager-resolver" \
        || { echo "traefik installation failed"; exit 1; }
}

argocd_install() {
    sudo zrun argo install \
        || { echo "argocd installation failed"; exit 1; }
}

argocd_add_project() {
    sudo zrun argocd add-project \
        --name "default" \
        --upsert || { echo "argocd project creation failed"; exit 1; }
}

# values "argo-git-repo-username" and "argo-git-repo-password" are stored in vault
# and are used to authenticate to the git repo
# to add a secret to vault, run the following command:
# sudo zrun vault add --key "argo-git-repo-username" --val "usernameXYZ"
# sudo zrun vault add --key "argo-git-repo-password" --val "password123"
argocd_add_repos() {
    sudo zrun argocd add-repo \
        --name "gitops" \
        --url "https://github.com/zcubbs/zrun-gitops-test-repo" \
        --type "git" \
        --use-vault \
        --username "argo-git-repo-username" \
        --password "argo-git-repo-password" \
        || { echo "argocd repo creation failed"; exit 1; }

}

wait_for_cluster() {
    echo "-------------------------------------------"
    echo "waiting for cluster to be ready..."
    until kubectl get nodes; do sleep 1; done
}

wait_for_argocd() {
    echo "-------------------------------------------"
    echo "waiting for argo-cd to be ready..."
    until kubectl -n argo-cd get pods | grep Running; do sleep 1; done
}

run_k9s() {
    sudo zrun k9s || { echo "Running k9s failed"; exit 1; }
}

main() {
#    zrun_info
    k3s_install
    k9s_install
    helm_install
#    wait_for_cluster
    cert_manager_install
    traefik_install
    argocd_install
#    wait_for_argocd
#    run_k9s
}

main
