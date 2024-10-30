#!/bin/bash
#
# curl -sSL https://raw.githubusercontent.com/rollchains/spawn/release/v0.50/install.sh | bash
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

URL="${BASE_URL}/spawn_${VERSION_NORMALIZED}_${OS}_${ARCH}.tar.gz"

TARGET_DIR="$(go env GOPATH)/bin"
mkdir -p "$TARGET_DIR"

mkdir -p /tmp/spawn-install
wget -O "/tmp/spawn-install/spawn.tar.gz" "$URL"
tar -xvf "/tmp/spawn-install/spawn.tar.gz" -C /tmp/spawn-install

mv "/tmp/spawn-install/spawn" "$TARGET_DIR"
chmod +x "$TARGET_DIR/spawn"


# see if the command spawn is avaliable in `which spawn`
if [ -x "$(command -v spawn)" ]; then
    echo "Spawn Installation complete. spawn is now available in $TARGET_DIR. Run the command 'spawn' from any terminal session."
else
    echo "Spawn is not available in your PATH"
    echo "To make spawn available from any terminal session, add the following line to your .bashrc or .zshrc:"
    echo 'export PATH="$PATH:$(go env GOPATH)/bin"'
fi


rm -rf /tmp/spawn-install
