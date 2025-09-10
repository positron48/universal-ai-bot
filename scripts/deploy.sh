#!/bin/bash

# Simple deployment script for Telegram Bot Skeleton
set -e

GITHUB_REPO="positron48/universal-ai-bot"  # Replace with your repository
APP_NAME="ai-bot"
APP_DIR="/var/www/ai-bot"

# Get latest release
get_latest_release() {
    curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | jq -r '.tag_name'
}

# Download and install binary
deploy() {
    echo "Deploying AI Bot..."
    
    # Create directories
    mkdir -p "$APP_DIR"/{bin,data,logs,configs}
    
    # Get version
    VERSION=$(get_latest_release)
    echo "Version: $VERSION"
    
    # Stop service before update
    if systemctl --user is-active ai-bot >/dev/null 2>&1; then
        echo "Stopping service..."
        systemctl --user stop ai-bot
    fi
    
    # Download binary
    BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/$APP_NAME-linux-amd64"
    curl -L -o "$APP_DIR/bin/$APP_NAME" "$BINARY_URL"
    chmod +x "$APP_DIR/bin/$APP_NAME"
    
    # Create .env if not exists
    if [ ! -f "$APP_DIR/.env" ]; then
        curl -s "https://raw.githubusercontent.com/$GITHUB_REPO/master/env.example" > "$APP_DIR/.env"
        echo "Created .env file. Edit it: nano $APP_DIR/.env"
    fi
    
    # Download configs
    echo "Downloading configs..."
    curl -s "https://raw.githubusercontent.com/$GITHUB_REPO/master/configs/config.yaml" > "$APP_DIR/configs/config.yaml" 2>/dev/null || echo "Config not found"
    
    # Restart service (if already configured)
    if systemctl --user is-enabled ai-bot >/dev/null 2>&1; then
        systemctl --user restart ai-bot
    else
        echo "Service not configured. Run: make setup"
        exit 1
    fi
    
    echo "Done! Status: systemctl --user status ai-bot"
}

# Update
update() {
    echo "Updating AI Bot..."
    systemctl --user stop ai-bot
    deploy
}

# Status
status() {
    systemctl --user status ai-bot
}

# Logs
logs() {
    journalctl --user -u ai-bot -f
}

case "${1:-deploy}" in
    deploy) deploy ;;
    update) update ;;
    status) status ;;
    logs) logs ;;
    *) echo "Usage: $0 [deploy|update|status|logs]" ;;
esac
