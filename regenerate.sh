#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")" || exit

# Check if gqlgen is installed
if ! command -v gqlgen &> /dev/null; then
    echo "gqlgen could not be found. Installing..."
    go install github.com/99designs/gqlgen@latest
fi

# Ensure all required packages are downloaded
echo "Fetching missing dependencies..."
go get github.com/99designs/gqlgen@v0.17.55
go get github.com/99designs/gqlgen/codegen/config@v0.17.55
go get github.com/99designs/gqlgen/internal/imports@v0.17.55
go get golang.org/x/tools/go/packages@v0.26.0
go get golang.org/x/tools/go/ast/astutil@v0.26.0
go get golang.org/x/tools/imports@v0.26.0
go get github.com/urfave/cli/v2@v2.27.5

# Run gqlgen to regenerate the code
echo "Regenerating GraphQL code..."
gqlgen generate

# Check if the generation was successful
if [ $? -eq 0 ]; then
    echo "GraphQL code generation completed successfully."
else
    echo "Failed to generate GraphQL code."
fi