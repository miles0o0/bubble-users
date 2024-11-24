#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")" || exit

# Check if gqlgen is installed
if ! command -v gqlgen &> /dev/null; then
    echo "gqlgen could not be found. Installing..."
    go install github.com/99designs/gqlgen@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Run gqlgen to regenerate the code
echo "Regenerating GraphQL code..."
gqlgen generate

# Check if the generation was successful
if [ $? -eq 0 ]; then
    echo "GraphQL code generation completed successfully."
else
    echo "Failed to generate GraphQL code."
fi

# Build the Docker image using the cache
docker build --cache-from=user-service:latest -t user-service .

# Run the Docker Compose services
docker compose --env-file .env up -d

# Wait for the container to be ready on port 8080
echo "Waiting for the container to be ready on port 8080..."

while ! nc -z localhost 8080; do
  sleep 1
done

migrate -path ./database/migrations -database "postgres://admin:admin@localhost:5432/userdb?sslmode=disable" up

# run tests
go test ./tests

# Hold the script and wait for user input to take down the services
read -p "Press any key to take down the Docker Compose services..."

# Take down the Docker Compose services
docker compose down