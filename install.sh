#!/bin/bash
set -e

NEXUS_DIR="/opt/nexus"
REPO_URL="https://github.com/oladayo21/nexus.git"
BRANCH="${NEXUS_BRANCH:-poc/vps-installation}"

echo "Installing Nexus..."

# 1. Check root
if [ "$EUID" -ne 0 ]; then
  echo "Error: Run as root"
  exit 1
fi

# 2. Install prerequisites - only if missing
echo "Checking prerequisites..."
apt-get update -qq

if ! [ -x "$(command -v curl)" ]; then
  echo "Installing curl..."
  apt-get install -y -qq curl
fi

if ! [ -x "$(command -v openssl)" ]; then
  echo "Installing openssl..."
  apt-get install -y -qq openssl
fi

if ! [ -x "$(command -v git)" ]; then
  echo "Installing git..."
  apt-get install -y -qq git
fi

# 3. Install Docker (if not present)
if ! [ -x "$(command -v docker)" ]; then
  echo "Installing Docker..."
  curl -sSL https://get.docker.com | sh
fi

# 4. Clone or update repo
if [ -d "$NEXUS_DIR" ]; then
  echo "Updating existing installation..."
  cd "$NEXUS_DIR"
  git fetch origin
  git checkout "$BRANCH"
  git pull
else
  echo "Cloning repository..."
  git clone -b "$BRANCH" "$REPO_URL" "$NEXUS_DIR"
  cd "$NEXUS_DIR"
fi

# 5. Generate .env with random POSTGRES_PASSWORD
if [ ! -f .env ]; then
  echo "Generating .env..."
  POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)
  cat > .env <<EOF
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
EOF
fi

# 6. Build and start services
echo "Building and starting Nexus (this may take a few minutes)..."
cd deploy
cp ../.env .env
docker compose up -d --build

# 7. Print access URL
IP=$(hostname -I | awk '{print $1}')
echo ""
echo "=========================================="
echo "Nexus installed successfully!"
echo "Access at: http://$IP"
echo "=========================================="
