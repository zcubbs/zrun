#!/bin/sh
set -e

echo "-------------------------------------------"
echo "Installing zrun"
go build -o ./bin/zrun
cp ./bin/zrun /usr/local/bin/zrun
zrun version

echo "-------------------------------------------"
echo "Running tests..."

zrun_hello() {
    zrun hello -v
}

main() {
    zrun_hello
}

main
