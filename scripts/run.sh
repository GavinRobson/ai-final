#!/bin/bash

OUT_NAME="server"

case "$(uname -s)" in
  Linux*)
    EXEC="./builds/${OUT_NAME}-linux-amd64"
    ;;
  Darwin*)
    EXEC="./builds/${OUT_NAME}-macos-amd64"
    ;;
  *)
    echo "Unsupported OS: $(usname -s)"
    exit 1
    ;;

esac

if [[ -f "$EXEC" ]]; then
  echo "Running $EXEC..."
  chmod +x "$EXEC"
  $EXEC
else
  echo "Executable not found: $EXEC"
  exit 1
fi
