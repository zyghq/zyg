#!/bin/bash
set -e

# Log all commands for debugging
exec 1> >(logger -s -t $(basename $0)) 2>&1

echo "Starting after-install script execution..."

# Reload systemd to recognize new or modified service file
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

# Enable the service to start on boot
echo "Enabling srv service..."
sudo systemctl enable srv.service

# Verify application binary exists and is executable
if [ ! -x "/usr/local/bin/app" ]; then
    echo "ERROR: Application binary is missing or not executable"
    exit 1
fi

echo "After-install script completed successfully"
exit 0
