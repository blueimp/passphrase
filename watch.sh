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
# Usage: ./watch.sh [chrome|safari|firefox]
#

SCRIPT='make -s lambda'

RELOAD_CHROME='tell application "Google Chrome"
  reload active tab of window 1
end tell'

RELOAD_SAFARI='tell application "Safari"
  set URL of document 1 to (URL of document 1)
end tell'

RELOAD_FIREFOX='activate application "Firefox"
tell application "System Events" to keystroke "r" using command down'

stop() {
  STATUS=$?
  make -s stop
  exit $STATUS
}

if command -v osascript > /dev/null 2>&1; then
  case "$1" in
    firefox)  OSASCRIPT=$RELOAD_FIREFOX;;
    safari)   OSASCRIPT=$RELOAD_SAFARI;;
    *)        OSASCRIPT=$RELOAD_CHROME;;
  esac
  SCRIPT="$SCRIPT && osascript -e '$OSASCRIPT'"
fi

cd "$(dirname "$0")" || exit 1

trap stop INT TERM
make -s start

while true; do
  find . -name '*.go' | entr -d -p sh -c "$SCRIPT"
done
