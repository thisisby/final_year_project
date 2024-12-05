#!/bin/sh

# Run migration
echo "Running migrations..."
./bin/migrate -up

# Start the main application
echo "Starting API server..."
exec "$@"
