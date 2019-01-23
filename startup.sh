#!/bin/bash
set -e
set -x

export LOCAL_IP=$(ip route get 8.8.8.8 | grep -oE 'src ([0-9\.]+)' | cut -d ' ' -f 2)
if [ "$SERVER_NAME" == "" ]; then
    SERVER_NAME=$LOCAL_IP
fi

echo "Starting the almighty Metrics Generator Tabajara..."
tabajara
