#!/bin/sh

#
# Starts a local API Gateway and a watch process for source file changes.
# Recompiles the lambda binary on changes.
# On MacOS also automatically reloads the active Chrome/Safari/Firefox tab.
#
# Requires make for the command executions.
# Requires AWS SAM for the local API Gateway.
# Requires Go for cross-compiling the lambda binary.
# Requires entr for the watch process.
#
# Usage: ./start.sh [chrome|safari|firefox]
#

stop() {
  STATUS=$?
  kill "$PID"
  exit $STATUS
}

start() {
  cd lambda || exit 1
  # Fake AWS credentials as fix for AWS SAM Local issue #134:
  # See also https://github.com/awslabs/aws-sam-local/issues/134
  AWS_ACCESS_KEY_ID=0 AWS_SECRET_ACCESS_KEY=0 sam local start-api &
  PID=$!
  cd ..
}

cd "$(dirname "$0")" || exit 1

trap stop INT TERM
start

set -- ./reload-browser.sh "$@" -- make -s lambda

while true; do
  find . -name '*.go' | entr -d -p "$@"
done
