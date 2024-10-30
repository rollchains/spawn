#!/bin/bash
#
# curl -sSL https://raw.githubusercontent.com/rollchains/spawn/release/v0.50/install.sh | bash -s -- v0.50
#

VERSION=${1:-"v0.50.10"}
BASE_URL="https://github.com/rollchains/spawn/releases/download/$VERSION"

ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
echo "Downloading spawn $VERSION for $OS/$ARCH..."

VERSION_NORMALIZED=$(echo $VERSION | tr -d 'v')

URL="${BASE_URL}/${APP_NAME}_${VERSION_NORMALIZED}_${OS}_${ARCH}.tar.gz"

TARGET_DIR="$(go env GOPATH)/bin"
mkdir -p "$TARGET_DIR"

wget -O "/tmp/${APP_NAME}.tar.gz" "$URL"

mkdir -p /tmp/spawn-install
tar -xvf "/tmp/${APP_NAME}.tar.gz" -C /tmp/spawn-install

mv "/tmp/spawn-install/${APP_NAME}" "$TARGET_DIR"
chmod +x "$TARGET_DIR/spawn"

echo "Installation complete. spawn is now available in $TARGET_DIR"
echo "To make spawn available from any terminal session, add the following line to your .bashrc or .zshrc:"
echo "export PATH=\"\$PATH:$TARGET_DIR\""

rm -rf /tmp/spawn-install
