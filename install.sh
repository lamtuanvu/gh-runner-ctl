#!/bin/sh
set -e

REPO="lamtuanvu/gh-runner-ctl"
BINARY="ghr"
INSTALL_DIR="/usr/local/bin"
VERSION=""

INSTALL_COMPLETIONS=""

usage() {
    echo "Usage: install.sh [-v VERSION] [-d INSTALL_DIR] [-c] [-h]"
    echo ""
    echo "Options:"
    echo "  -v VERSION      Install a specific version (e.g. v0.1.0). Default: latest"
    echo "  -d INSTALL_DIR  Installation directory. Default: /usr/local/bin"
    echo "  -c              Install shell completions after binary install"
    echo "  -h              Show this help message"
    exit 0
}

while getopts "v:d:ch" opt; do
    case "$opt" in
        v) VERSION="$OPTARG" ;;
        d) INSTALL_DIR="$OPTARG" ;;
        c) INSTALL_COMPLETIONS=1 ;;
        h) usage ;;
        *) usage ;;
    esac
done

detect_os() {
    os="$(uname -s)"
    case "$os" in
        Linux)  echo "linux" ;;
        Darwin) echo "darwin" ;;
        *)      echo "Unsupported OS: $os" >&2; exit 1 ;;
    esac
}

detect_arch() {
    arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac
}

get_latest_version() {
    url="https://api.github.com/repos/${REPO}/releases/latest"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" | grep '"tag_name"' | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$url" | grep '"tag_name"' | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/'
    else
        echo "Error: curl or wget is required" >&2
        exit 1
    fi
}

download() {
    url="$1"
    output="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$output" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$output" "$url"
    else
        echo "Error: curl or wget is required" >&2
        exit 1
    fi
}

verify_checksum() {
    archive="$1"
    checksums="$2"
    filename="$(basename "$archive")"

    expected="$(grep "$filename" "$checksums" | awk '{print $1}')"
    if [ -z "$expected" ]; then
        echo "Error: checksum not found for $filename" >&2
        exit 1
    fi

    if command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$archive" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$archive" | awk '{print $1}')"
    else
        echo "Warning: no sha256 tool found, skipping checksum verification" >&2
        return 0
    fi

    if [ "$expected" != "$actual" ]; then
        echo "Error: checksum mismatch" >&2
        echo "  expected: $expected" >&2
        echo "  actual:   $actual" >&2
        exit 1
    fi
}

main() {
    OS="$(detect_os)"
    ARCH="$(detect_arch)"

    if [ -z "$VERSION" ]; then
        echo "Fetching latest version..."
        VERSION="$(get_latest_version)"
        if [ -z "$VERSION" ]; then
            echo "Error: could not determine latest version" >&2
            exit 1
        fi
    fi

    # Strip leading v for filename
    VERSION_NUM="${VERSION#v}"

    echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

    ARCHIVE="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    CHECKSUMS="${BINARY}_${VERSION_NUM}_checksums.txt"
    BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"

    TMPDIR="$(mktemp -d)"
    trap 'rm -rf "$TMPDIR"' EXIT

    echo "Downloading ${ARCHIVE}..."
    download "${BASE_URL}/${ARCHIVE}" "${TMPDIR}/${ARCHIVE}"

    echo "Downloading checksums..."
    download "${BASE_URL}/${CHECKSUMS}" "${TMPDIR}/${CHECKSUMS}"

    echo "Verifying checksum..."
    verify_checksum "${TMPDIR}/${ARCHIVE}" "${TMPDIR}/${CHECKSUMS}"

    echo "Extracting..."
    tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"

    echo "Installing to ${INSTALL_DIR}..."
    if [ -w "$INSTALL_DIR" ]; then
        install -m 755 "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
        sudo install -m 755 "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi

    echo "Successfully installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
}

install_completions() {
    echo "Installing shell completions..."
    if "${INSTALL_DIR}/${BINARY}" completion install; then
        echo "Shell completions installed successfully."
    else
        echo "Warning: could not install shell completions." >&2
    fi
}

main

if [ "$INSTALL_COMPLETIONS" = "1" ]; then
    install_completions
elif [ -t 0 ]; then
    printf "Would you like to install shell completions? [y/N] "
    read -r answer
    case "$answer" in
        [yY]|[yY][eE][sS]) install_completions ;;
        *) echo "Skipping shell completions." ;;
    esac
fi
