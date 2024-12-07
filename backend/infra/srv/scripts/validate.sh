#!/bin/bash
set -e

SERVICE_NAME="srv.service"

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

main() {
  local service=${3:-$SERVICE_NAME}

  check_service "$service" || exit 1

  echo "All validations passed successfully"
  exit 0
}

# Run main with command line arguments
main "$@"
