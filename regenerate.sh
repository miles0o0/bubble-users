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