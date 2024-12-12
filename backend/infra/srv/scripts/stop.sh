#!/bin/bash
set -e

# Log all commands for debugging
exec 1> >(logger -s -t $(basename $0)) 2>&1

echo "Starting application stop script..."

# Check if the service is running
if systemctl is-active --quiet srv.service; then
    echo "Stopping srv service..."
    sudo systemctl stop srv.service
    
    # Wait for the service to fully stop (max 30 seconds)
    COUNTER=0
    while systemctl is-active --quiet srv.service && [ $COUNTER -lt 30 ]; do
        sleep 1
        let COUNTER=COUNTER+1
        echo "Waiting for service to stop... ($COUNTER seconds)"
    done
    
    # Check if service successfully stopped
    if systemctl is-active --quiet srv.service; then
        echo "ERROR: Failed to stop the service within timeout period"
        exit 1
    fi
else
    echo "Service was not running"
fi

echo "Application stop script completed successfully"
exit 0
