#!/usr/bin/env bash
set -euo pipefail

# -------------------------------------------
# DEFAULT CONFIGURATION
# -------------------------------------------

# Set your static default URL here:
DEFAULT_BASE_URL="https://raw.githubusercontent.com/danielr1996/hardware_exporter/refs/heads/main/"

# This variable may be overridden by --url
BASE_URL="$DEFAULT_BASE_URL"

SERVICE_NAME="hardware-exporter.service"
BINARY_NAME="hardware_exporter"
INSTALL_PATH="/usr/local/bin/${BINARY_NAME}"
SERVICE_PATH="/etc/systemd/system/${SERVICE_NAME}"
SERVICE_USER="hardware_exporter"

usage() {
  echo "Usage: $0 [--url <custom_base_url>]"
  echo
  echo "If --url is omitted, default is:"
  echo "  $DEFAULT_BASE_URL"
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --url)
      BASE_URL="$2"
      shift 2
      ;;
    *)
      echo "Unknown argument: $1"
      usage
      ;;
  esac
done

# -------------------------------------------
# FUNCTIONS
# -------------------------------------------

download() {
  local url="$1"
  local out="$2"

  echo "Downloading $url → $out"
  curl -fsSL "$url" -o "$out"
}

ensure_user() {
  if ! id -u "$SERVICE_USER" >/dev/null 2>&1; then
    echo "Creating system user: $SERVICE_USER"
    sudo useradd -r -s /usr/sbin/nologin "$SERVICE_USER"
  fi
}

install_binary() {
  echo "Installing binary to ${INSTALL_PATH}"
  sudo mv "./${BINARY_NAME}" "$INSTALL_PATH"
  sudo chmod 755 "$INSTALL_PATH"
  sudo chown "$SERVICE_USER":"$SERVICE_USER" "$INSTALL_PATH"
}

install_service() {
  echo "Installing systemd service to ${SERVICE_PATH}"
  sudo mv "./${SERVICE_NAME}" "$SERVICE_PATH"
  sudo chmod 644 "$SERVICE_PATH"
}

restart_systemd() {
  echo "Reloading systemd daemon…"
  sudo systemctl daemon-reload

  echo "Enabling service ${SERVICE_NAME}…"
  sudo systemctl enable "${SERVICE_NAME}"

  echo "Starting service ${SERVICE_NAME}…"
  sudo systemctl restart "${SERVICE_NAME}"

  echo "Service status:"
  sudo systemctl status "${SERVICE_NAME}" --no-pager
}

# -------------------------------------------
# MAIN EXECUTION
# -------------------------------------------

echo "=== Inventory Exporter Installer ==="
echo "Using download URL: $BASE_URL"
echo

BIN_URL="${BASE_URL}/agent/dist/hardware_exporter"
SVC_URL="${BASE_URL}/agent/dist/hardware-exporter.service"

download "$BIN_URL" "${BINARY_NAME}"
download "$SVC_URL" "${SERVICE_NAME}"

ensure_user
install_binary
install_service
restart_systemd

echo
echo "Installation complete."
echo "Metrics available at: http://<host>:9105/metrics"
