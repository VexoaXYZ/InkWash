#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/server-data"
if [ ! -d "$SCRIPT_DIR/txData" ]; then
    echo "First run detected - starting txAdmin setup..."
    bash "$SCRIPT_DIR/server-core/run.sh"
else
    bash "$SCRIPT_DIR/server-core/run.sh" +exec server.cfg
fi