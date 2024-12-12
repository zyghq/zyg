#!/bin/bash
set -e

# Create necessary directories if they don't exist
sudo mkdir -p /usr/local/bin

# Set ownership to ubuntu user
sudo chown ubuntu:ubuntu /usr/local/bin/app || true

# Set proper permissions
sudo chmod 755 /usr/local/bin/app || true

# Ensure systemd directory exists and has correct permissions
sudo mkdir -p /etc/systemd/system
sudo chmod 755 /etc/systemd/system

# Set ownership to ubuntu user
sudo chown ubuntu:ubuntu /etc/systemd/system/srv.service || true

# Stop the service if it's running (ignore if it fails)
sudo systemctl stop srv.service || true
