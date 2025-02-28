#!/bin/bash

# Load environment variables from .env file
export $(grep -v '^#' .env | xargs)

# Run the tests, skipping those that require external services
go test ./pkg/... ./middleware/...

# To run all tests (including those that require external services), uncomment the line below
# go test ./... 