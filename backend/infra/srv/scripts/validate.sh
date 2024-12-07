#!/bin/bash
set -e

# Default values
DEFAULT_HOST="localhost"
DEFAULT_PORT="8080"
DEFAULT_SERVICE="srv.service"

# Function to check if service is running
check_service() {
  local service_name=$1

  echo "Checking if service $service_name is running..."
  if ! systemctl is-active --quiet "$service_name"; then
    echo "ERROR: Service $service_name is not running"
    systemctl status "$service_name"
    return 1
  fi
  echo "Service $service_name is running"
  return 0
}

# Function to check web endpoint
check_web() {
  local host=$1
  local port=$2
  local url="http://${host}:${port}/"

  echo "Checking web endpoint at $url..."
  if ! curl -s -f --connect-timeout 5 "$url" >/dev/null; then
    echo "ERROR: Web service is not responding at $url"
    return 1
  fi
  echo "Web endpoint at $url is accessible"
  return 0
}

# Main function
main() {
  local host=${1:-$DEFAULT_HOST}
  local port=${2:-$DEFAULT_PORT}
  local service=${3:-$DEFAULT_SERVICE}

  echo "Starting validation with host=$host, port=$port, service=$service"

  # Run checks
  check_service "$service" || exit 1
  # check_web "$host" "$port" || exit 1

  echo "All validations passed successfully"
  exit 0
}

# Run main with command line arguments
main "$@"
