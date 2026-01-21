#!/bin/bash

# Build script for Multi-Model Database Web Admin

echo "Building Multi-Model Database Web Admin..."

# Build the frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build the backend
echo "Building backend..."
cd backend
go mod tidy
go build -o webadmin
cd ..

echo "Build completed successfully!"
echo ""
echo "To run the application:"
echo "1. Start the database engine: cd .. && go run main.go"
echo "2. In another terminal, start the web admin: cd backend && ./webadmin"
echo ""
echo "Or use Docker Compose to run both together:"
echo "docker-compose up --build"