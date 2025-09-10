#!/bin/bash

# Первоначальная настройка Budget Bot (выполняется один раз с sudo)
set -e

APP_DIR="/var/www/ai-bot"
SERVICE_NAME="ai-bot"

echo "Настройка AI Bot..."

# Создать директории
sudo mkdir -p "$APP_DIR"/{bin,data,logs,configs,migrations}
sudo chown -R "$USER:$USER" "$APP_DIR"

# Создать пользовательский systemd сервис
mkdir -p ~/.config/systemd/user
tee ~/.config/systemd/user/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=AI Bot
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

# Включить пользовательский сервис
systemctl --user daemon-reload
systemctl --user enable $SERVICE_NAME
systemctl --user start $SERVICE_NAME

# Включить автозапуск пользовательских сервисов
sudo loginctl enable-linger $USER

echo "Настройка завершена!"
echo "Теперь можно использовать: make deploy"
