# Multi-Model Database Web Administration Interface

A web-based administration interface for managing the multi-model database engine. This interface provides a user-friendly way to interact with document, key-value, column-family, and graph stores through a unified dashboard.

## Features

- **Dashboard**: Overview of database status, statistics, and cluster health
- **Document Store Management**: Create, read, update, and delete documents in collections
- **Key-Value Store Management**: Manage key-value pairs with different data types
- **Column Store Management**: Work with column-family data structures
- **Graph Store Management**: Visualize and manage nodes and edges in graph databases
- **Cluster Management**: Monitor and manage cluster nodes and their status
- **Real-time Monitoring**: Live status updates and performance metrics

## Architecture

The web administration interface consists of:

1. **Frontend**: A responsive web application built with HTML, CSS, and JavaScript using Bootstrap for styling
2. **Backend API**: A Go-based API server that communicates with the database engine
3. **Proxy Layer**: Handles communication between the frontend and the database engine

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js (for local frontend development)

### Running with Docker Compose

The easiest way to run the entire stack is using Docker Compose:

```bash
# Navigate to the webapp directory
cd webapp

# Build and start the services
docker-compose up -d

# Access the web admin interface at http://localhost:3000
# Access the database engine API at http://localhost:8080
```

### Local Development

For frontend development:

```bash
cd webapp/frontend
npm install
npm run dev  # Starts a development server on localhost:5000
```

For backend development:

```bash
cd webapp/backend
go mod tidy
go run api.go
```

## API Endpoints

The web admin backend provides the following API endpoints:

- `GET /api/health` - Check database engine health
- `GET /api/cluster/status` - Get cluster status
- `GET /api/documents/collections` - List document collections
- `GET /api/documents/{collection}` - Get documents from a collection
- `POST /api/documents/{collection}/{id}` - Create a document
- `PUT /api/documents/{collection}/{id}` - Update a document
- `DELETE /api/documents/{collection}/{id}` - Delete a document
- `GET /api/kv/{key}` - Get a key-value pair
- `POST/PUT /api/kv/{key}` - Set a key-value pair
- `DELETE /api/kv/{key}` - Delete a key
- `GET /api/columns/{family}/{row}/{column}` - Get a column value
- `POST/PUT /api/columns/{family}/{row}/{column}` - Insert a column value
- `GET /api/graph/nodes/{id}` - Get a graph node
- `POST /api/graph/nodes` - Create a graph node
- `GET /api/graph/edges/{id}` - Get a graph edge
- `POST /api/graph/edges` - Create a graph edge

## Configuration

The application can be configured using environment variables:

- `DB_URL`: URL of the database engine (default: http://localhost:8080)
- `PORT`: Port to run the web admin on (default: 3000)

## Security Considerations

- Authentication and authorization should be implemented in production
- HTTPS should be used for secure communication
- Input validation should be implemented on both frontend and backend

## Deployment

The application is designed for cloud-native deployment with Docker containers. It can be deployed to Kubernetes, Docker Swarm, or any cloud platform that supports containerized applications.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

This project is licensed under the MIT License.