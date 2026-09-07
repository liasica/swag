#!/usr/bin/env bash

# 一键安装 swag：
#   curl -sSfL https://raw.githubusercontent.com/liasica/swag/v2-adaptation/install.sh | bash
#
# 环境变量：
#   SWAG_VERSION  指定版本，默认取最新发布
#   SWAG_BIN_DIR  安装目录，默认 /usr/local/bin

set -e

REPO="liasica/swag"
BIN_DIR="${SWAG_BIN_DIR:-/usr/local/bin}"
VERSION="${SWAG_VERSION:-}"

os=$(uname -s)
case "$os" in
    Linux) os="Linux" ;;
    Darwin) os="Darwin" ;;
    *)
        echo "Unsupported OS: ${os}" >&2
        exit 1
        ;;
esac

arch=$(uname -m)
case "$arch" in
    x86_64 | amd64) arch="x86_64" ;;
    arm64 | aarch64) arch="arm64" ;;
    *)
        echo "Unsupported architecture: ${arch}" >&2
        exit 1
        ;;
esac

if [ -z "$VERSION" ]; then
    VERSION=$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" |
        sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' |
        head -n 1)
fi

if [ -z "$VERSION" ]; then
    echo "Failed to resolve the latest version" >&2
    exit 1
fi

archive="swag_${VERSION#v}_${os}_${arch}.tar.gz"
base_url="https://github.com/${REPO}/releases/download/${VERSION}"

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

echo "Downloading ${archive} (${VERSION})"
curl -sSfL "${base_url}/${archive}" -o "${tmpdir}/${archive}"

# 校验下载的压缩包，checksums.txt 缺失时跳过
if curl -sSfL "${base_url}/checksums.txt" -o "${tmpdir}/checksums.txt" 2> /dev/null; then
    if command -v sha256sum > /dev/null 2>&1; then
        sha_check="sha256sum -c -"
    elif command -v shasum > /dev/null 2>&1; then
        sha_check="shasum -a 256 -c -"
    fi

    if [ -n "${sha_check:-}" ]; then
        if ! grep " ${archive}\$" "${tmpdir}/checksums.txt" | (cd "$tmpdir" && $sha_check) > /dev/null; then
            echo "Checksum verification failed" >&2
            exit 1
        fi
        echo "Checksum verified"
    fi
fi

tar -xzf "${tmpdir}/${archive}" -C "$tmpdir" swag

if [ -w "$BIN_DIR" ]; then
    install -m 0755 "${tmpdir}/swag" "${BIN_DIR}/swag"
elif command -v sudo > /dev/null 2>&1; then
    echo "Installing to ${BIN_DIR} with sudo"
    sudo install -m 0755 "${tmpdir}/swag" "${BIN_DIR}/swag"
else
    echo "No write permission for ${BIN_DIR}, set SWAG_BIN_DIR to a writable directory" >&2
    exit 1
fi

echo "Installed to ${BIN_DIR}/swag: $("${BIN_DIR}/swag" --version)"
