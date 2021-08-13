#!/usr/bin/env bash

set -euxo pipefail

if [ ! -f "/dev/shm/cosign_private_key" ]; then
    echo "No Cosign private key found, skipping signatures"
    exit 0
fi

cosign sign-blob -key /dev/shm/cosign_private_key "$1" > "$2"
