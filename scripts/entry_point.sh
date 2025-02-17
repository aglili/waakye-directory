#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
timeout=30
while ! nc -z postgres 5432; do
  sleep 1
  timeout=$((timeout - 1))
  if [ $timeout -eq 0 ]; then
    echo "PostgreSQL is not ready, exiting..."
    exit 1
  fi
done
echo "PostgreSQL is ready!"

# Run migrations
echo "Running database migrations..."
migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
if [ $? -ne 0 ]; then
  echo "Failed to run migrations, exiting..."
  exit 1
fi

# Debug: Show current directory and contents
echo "Current directory: $(pwd)"
echo "Contents of /app:"
ls -la /app

# Start the application
echo "Starting the application..."
exec ./main  # Changed to use relative path