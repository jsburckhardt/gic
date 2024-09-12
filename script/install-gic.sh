#!/bin/bash

# Variables
REPO_OWNER="jsburckhardt"
REPO_NAME="gic"
BINARY_NAME="gic"

# Function to get the latest version from GitHub API
get_latest_version() {
    LATEST_URL="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
    curl -s "$LATEST_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Check if a version is passed as an argument
if [ -z "$1" ]; then
    # No version provided, get the latest version
    VERSION=$(get_latest_version)
    echo "No version provided, installing the latest version: $VERSION"
else
    VERSION=$1
    echo "Installing specified version: $VERSION"
fi

# Determine the OS and architecture following the naming template
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    ARCH="x86_64"
elif [ "$ARCH" == "i686" ]; then
    ARCH="i386"
elif [ "$ARCH" == "armv7l" ]; then
    ARCH="armv7"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

# Construct the download URL to match the naming template
DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${REPO_NAME}_$(echo "$OS" | sed 's/.*/\u&/')_$ARCH.tar.gz"

# Create a temporary directory for the download
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR" || exit

# Download the binary tarball
echo "Downloading $BINARY_NAME from $DOWNLOAD_URL"
curl -LO "$DOWNLOAD_URL"

# Extract the tarball
echo "Extracting the tarball"
tar -xzf "${REPO_NAME}_$(echo "$OS" | sed 's/.*/\u&/')_$ARCH.tar.gz"

# Move the binary to /usr/local/bin
echo "Installing $BINARY_NAME"
sudo mv $BINARY_NAME /usr/local/bin/

# Cleanup
cd - || exit
rm -rf "$TMP_DIR"

# Verify installation
echo "Verifying installation"
$BINARY_NAME --version
