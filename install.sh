#!/usr/bin/env sh
set -e

REPO="caioreix/agency-cli"
BINARY="agency-cli"
RELEASES_URL="https://github.com/${REPO}/releases/latest/download"

# Detect OS
case "$(uname -s)" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    echo "For Windows, see: https://github.com/${REPO}#installation" >&2
    exit 1
    ;;
esac

# Detect architecture
case "$(uname -m)" in
  x86_64)          ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

DOWNLOAD_URL="${RELEASES_URL}/${BINARY}-${OS}-${ARCH}"
TMP_FILE="$(mktemp)"

echo "Downloading ${BINARY} (${OS}/${ARCH})..."
curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_FILE}"
chmod +x "${TMP_FILE}"

# Determine install destination
INSTALL_DIR=""
if [ -w "/usr/local/bin" ] || sudo -n true 2>/dev/null; then
  INSTALL_DIR="/usr/local/bin"
  SUDO="sudo"
else
  INSTALL_DIR="${HOME}/.local/bin"
  SUDO=""
  mkdir -p "${INSTALL_DIR}"
fi

${SUDO} mv "${TMP_FILE}" "${INSTALL_DIR}/${BINARY}"

echo ""
echo "${BINARY} installed to ${INSTALL_DIR}/${BINARY}"

# Warn if $HOME/.local/bin is not in PATH
if [ "${INSTALL_DIR}" = "${HOME}/.local/bin" ]; then
  case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) ;;
    *)
      echo ""
      echo "WARNING: ${INSTALL_DIR} is not in your PATH."
      echo "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
      echo ""
      echo '  export PATH="$HOME/.local/bin:$PATH"'
      echo ""
      echo "Then restart your shell or run: source ~/.bashrc"
      ;;
  esac
fi

echo ""
${INSTALL_DIR}/${BINARY} --version
