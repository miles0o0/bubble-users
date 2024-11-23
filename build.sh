#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")" || exit

./regenerate.sh

# Build the Docker image using the cache
docker build -t bubble-user .

# Run the Docker Compose services
docker compose up -d

# Wait for the container to be ready on port 8080
echo "Waiting for the container to be ready on port 8080..."
while ! nc -z localhost 8080; do
  sleep 1
done

# run tests
go test ./tests

# take it down
docker compose down