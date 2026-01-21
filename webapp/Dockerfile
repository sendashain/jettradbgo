# Multi-stage build for the web admin application

# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/src ./src
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.21-alpine AS backend-build
WORKDIR /app
COPY webapp/backend/go.mod webapp/backend/go.sum ./
RUN go mod download
COPY webapp/backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o webadmin .

# Stage 3: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the backend binary
COPY --from=backend-build /app/webadmin .

# Copy the frontend assets
COPY --from=frontend-build /app/dist ./webapp/frontend/dist

# Expose port
EXPOSE 3000

# Run the application
CMD ["./webadmin"]