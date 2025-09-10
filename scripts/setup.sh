#!/bin/bash

# Universal AI Bot Setup Script (run once with sudo)
set -e

# Configuration - can be overridden by environment variables
GITHUB_REPO="${GITHUB_REPO:-positron48/universal-ai-bot}"
APP_DIR="${PWD}"  # Use current directory

echo "🚀 Setting up Universal AI Bot..."

# Create directories in current location
mkdir -p "$APP_DIR"/{bin,logs}

echo "📁 Creating directories..."

# Download necessary files from GitHub
echo "📥 Downloading configuration files..."

# Download .env.example if .env doesn't exist
if [ ! -f "$APP_DIR/.env" ]; then
    curl -s "https://raw.githubusercontent.com/$GITHUB_REPO/master/env.example" > "$APP_DIR/.env"
    echo "✅ Created .env file from example"
    echo "📝 Edit the file: nano $APP_DIR/.env"
fi

# Load SERVICE_NAME from .env file
if [ -f "$APP_DIR/.env" ]; then
    export $(grep -v '^#' "$APP_DIR/.env" | xargs)
fi
SERVICE_NAME="${SERVICE_NAME:-ai-bot}"

# Download deployment Makefile
curl -s "https://raw.githubusercontent.com/$GITHUB_REPO/master/Makefile.deploy" > "$APP_DIR/Makefile" 2>/dev/null || {
    # If Makefile.deploy doesn't exist, create a minimal Makefile
    cat > "$APP_DIR/Makefile" << 'EOF'
# Deployment Makefile for Universal AI Bot
GITHUB_REPO ?= positron48/universal-ai-bot
APP_NAME ?= universal-ai-bot

# Load SERVICE_NAME from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif
SERVICE_NAME ?= ai-bot

# Get latest release
get_latest_release() {
    curl -s "https://api.github.com/repos/$(GITHUB_REPO)/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install binary
deploy:
	@echo "Deploying AI Bot..."
	@mkdir -p bin
	@VERSION=$$(get_latest_release); \
	echo "Version: $$VERSION"; \
	if systemctl --user is-active $(SERVICE_NAME) >/dev/null 2>&1; then \
		echo "Stopping service..."; \
		systemctl --user stop $(SERVICE_NAME); \
	fi; \
	curl -L -o "bin/$(APP_NAME)" "https://github.com/$(GITHUB_REPO)/releases/download/$$VERSION/$(APP_NAME)-linux_amd64"; \
	chmod +x "bin/$(APP_NAME)"; \
	if systemctl --user is-enabled $(SERVICE_NAME) >/dev/null 2>&1; then \
		systemctl --user restart $(SERVICE_NAME); \
	else \
		echo "Service not configured. Run setup first."; \
		exit 1; \
	fi; \
	echo "Done! Status: systemctl --user status $(SERVICE_NAME)"

update: deploy

status:
	@systemctl --user status $(SERVICE_NAME)

logs:
	@journalctl --user -u $(SERVICE_NAME) -f

restart:
	@systemctl --user restart $(SERVICE_NAME)

stop:
	@systemctl --user stop $(SERVICE_NAME)

start:
	@systemctl --user start $(SERVICE_NAME)
EOF
    echo "✅ Created deployment Makefile"
}

# Download deployment script
curl -s "https://raw.githubusercontent.com/$GITHUB_REPO/master/scripts/deploy.sh" > "$APP_DIR/deploy.sh"
chmod +x "$APP_DIR/deploy.sh"
echo "✅ Downloaded deployment script"

# Create user systemd service
echo "⚙️ Setting up systemd service..."
mkdir -p ~/.config/systemd/user
tee ~/.config/systemd/user/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=Universal AI Bot
After=network.target

[Service]
Type=simple
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/bin/$SERVICE_NAME
Restart=always
EnvironmentFile=$APP_DIR/.env

[Install]
WantedBy=default.target
EOF

# Enable user service
systemctl --user daemon-reload
systemctl --user enable $SERVICE_NAME

# Enable user service autostart
sudo loginctl enable-linger $USER

echo ""
echo "🎉 Setup completed!"
echo ""
echo "📋 Next steps:"
echo "1. Edit configuration: nano $APP_DIR/.env"
echo "2. Run deployment: make deploy"
echo "3. Check status: make status"
echo "4. View logs: make logs"
echo ""
echo "💡 Available commands:"
echo "   make deploy  - Deploy/update bot"
echo "   make status  - Check status"
echo "   make logs    - View logs"
echo "   make restart - Restart service"
