#!/bin/bash
set -e

# Log all commands for debugging
exec 1> >(logger -s -t $(basename $0)) 2>&1

echo "Starting application start script..."

# Verify the application binary exists
if [ ! -f "/usr/local/bin/app" ]; then
    echo "ERROR: Application binary not found"
    exit 1
fi

# Ensure proper permissions
sudo chown ubuntu:ubuntu /usr/local/bin/app
sudo chmod 755 /usr/local/bin/app

# Start the service
echo "Starting srv service..."
sudo systemctl start srv.service

# Wait for the service to start (max 30 seconds)
COUNTER=0
while ! systemctl is-active --quiet srv.service && [ $COUNTER -lt 30 ]; do
    sleep 1
    let COUNTER=COUNTER+1
    echo "Waiting for service to start... ($COUNTER seconds)"
done

# Verify the service started successfully
if ! systemctl is-active --quiet srv.service; then
    echo "ERROR: Failed to start the service"
    # Print the last few lines of the service logs for debugging
    echo "Service logs:"
    sudo journalctl -u srv.service -n 50 --no-pager
    exit 1
fi

echo "Application start script completed successfully"
exit 0