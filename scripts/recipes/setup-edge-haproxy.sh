#!/bin/sh
set -e

echo "-------------------------------------------"
echo "test setup edge haproxy"

update_from_git() {
    sudo zrun haproxy update-from-git \
        -r "$HAPROXY_GIT_REPO" \
        -u "$HAPROXY_GIT_USERNAME" \
        -p "$HAPROXY_GIT_PASSWORD" \
        -f "$HAPROXY_GIT_CONFIG_FILE" \
        || { echo "haproxy update-from-git failed"; exit 1; }
}

update_from_git

echo "-------------------------------------------"
echo "test setup edge haproxy git update cronjob"

# required ENV vars:
# export HAPROXY_GIT_REPO=""
# export HAPROXY_GIT_CONFIG_FILE=""
# export HAPROXY_GIT_USERNAME=""
# export HAPROXY_GIT_PASSWORD=""

setup_git_update_cronjob() {
    sudo zrun haproxy setup-git-cron \
        || { echo "haproxy setup-git-update-cronjob failed"; exit 1; }
}

setup_git_update_cronjob
