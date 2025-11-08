#!/bin/bash

export SKIPAUTH=true
set -o errexit
set -o pipefail
set -o errtrace

err_report() {
    echo "Error running '$BASH_COMMAND' [rc=$?], line $1"
}

trap 'err_report $LINENO' ERR

if [ -n "$GRACEFUL_TIMEOUT" ]; then
    TIMEOUT_ARG="--graceful-timeout=$GRACEFUL_TIMEOUT"
fi

CMD="./cmd/adronme-code-server/main.go"

go run "$CMD" "$TIMEOUT_ARG" "$@"
