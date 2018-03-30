#!/bin/sh

#
# Starts the given make target and a watch process for source file changes.
# Reloads the active Chrome/Safari/Firefox tab.
# Optionally executes a command on each source file change.
#
# Requires make for the command execution.
# Requires entr for the watch process.
#
# Usage: ./watch.sh target [chrome|safari|firefox] [-- cmd args...]
#

stop() {
  STATUS=$?
  kill "$PID"
  exit $STATUS
}

start() {
  make -s "$1" &
  PID=$!
}

cd "$(dirname "$0")" || exit 1

trap stop INT TERM
start "$1"
shift

while true; do
  find . -name '*.go' | entr -d -p ./reload-browser.sh "$@"
done
