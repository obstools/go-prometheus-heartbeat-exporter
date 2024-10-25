#!/bin/sh
set -e

RELEASES_URL="https://github.com/obstools/go-prometheus-heartbeat-exporter/releases"
BINARY_NAME="heartbeat"
ARCH_TYPE=".tar.gz"
TAR_FILE="$BINARY_NAME$ARCH_TYPE"

latest_release() {
  curl -sL -o /dev/null -w '%{url_effective}' "$RELEASES_URL/latest" | rev | cut -f 1 -d '/'| rev
}

remove_tmp_download() {
  rm -f "$TAR_FILE"
}

download() {
  test -z "$VERSION" && VERSION="$(latest_release)"
  test -z "$VERSION" && {
    echo "Unable to get heartbeat release." >&2
    exit 1
  }
  remove_tmp_download
  curl -s -L -o "$TAR_FILE" "$RELEASES_URL/download/$VERSION/heartbeat_$(uname -s)_$(uname -m)$ARCH_TYPE"
}

extract() {
  tar -zxf "$TAR_FILE" "$BINARY_NAME"
  remove_tmp_download
}

download
extract
