#!/bin/bash
set -e

# Log all commands for debugging
exec 1> >(logger -s -t $(basename $0)) 2>&1

echo "Starting service validation..."

# Function to check if a port is accepting connections
check_port() {
    local port=$1
    local timeout=$2
    nc -z -w$timeout localhost $port
    return $?
}

# Function to check HTTP endpoint (adjust URL as needed)
check_http() {
    local url=$1
    local timeout=$2
    curl -s -f -m $timeout "$url" > /dev/null
    return $?
}

# Function to check logs using journalctl
check_logs() {
    local service=$1
    local time_window=$2
    
    # Check for errors in recent logs
    ERROR_COUNT=$(journalctl -u $service --since "$time_window ago" -p err | wc -l)
    
    if [ $ERROR_COUNT -gt 0 ]; then
        echo "WARNING: Found $ERROR_COUNT errors in the last $time_window"
        echo "Recent errors:"
        journalctl -u $service --since "$time_window ago" -p err --no-pager
        return 1
    fi
    return 0
}


# Check if service is running
if ! systemctl is-active --quiet srv.service; then
    echo "ERROR: Service is not running"
    echo "Service status:"
    systemctl status srv.service
    exit 1
fi

# Check process resources
echo "Checking process resources..."
SERVICE_PID=$(systemctl show --property MainPID --value srv.service)
if [ -z "$SERVICE_PID" ] || [ "$SERVICE_PID" -eq "0" ]; then
    echo "ERROR: Cannot find service PID"
    exit 1
fi

# Check memory usage
MEM_USAGE=$(ps -o pmem= -p $SERVICE_PID)
if [ $(echo "$MEM_USAGE > 90" | bc -l) -eq 1 ]; then
    echo "WARNING: High memory usage detected: $MEM_USAGE%"
fi

# Check port availability (adjust port as needed)
echo "Checking service port..."
if ! check_port 8080 5; then
    echo "ERROR: Service port is not responding"
    exit 1
fi

# Check HTTP endpoint (adjust URL as needed)
echo "Checking service health endpoint..."
if ! check_http "http://localhost:8080/" 5; then
    echo "ERROR: Health check failed"
    exit 1
fi

# Check logs for the last 2 minutes
echo "Checking service logs..."
if ! check_logs srv.service "2 minutes"; then
    echo "WARNING: Found errors in recent logs"
    # You might want to exit 1 here depending on your requirements
fi
