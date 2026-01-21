# Multi-Model Database Engine for Cloud Native Applications

A distributed multi-model database engine built in Go that supports document, key-value, column-family, and graph data models. Designed for cloud-native deployments with built-in clustering and replication capabilities.

## Features

- **Multi-Model Support**: Document, Key-Value, Column-Family, and Graph models
- **Distributed Architecture**: Built-in clustering with gossip protocol
- **Replication**: Configurable replication factor for high availability
- **Cloud-Native**: Container-friendly with horizontal scaling
- **HTTP API**: RESTful interface for all operations
- **Consistency Levels**: Configurable consistency settings
- **Web Admin Interface**: Comprehensive web-based management dashboard

## Data Models Supported

### Document Model
Store flexible JSON-like documents in collections

### Key-Value Model
Simple key-value storage with rich value types

### Column-Family Model
Wide-column storage similar to Cassandra

### Graph Model
Nodes and edges with properties for graph relationships

## API Endpoints

### Health Check
```
GET /health
```

### Document Store
```
POST   /docs/{collection}/{id}     # Create document
GET    /docs/{collection}/{id}     # Get document
PUT    /docs/{collection}/{id}     # Update document
DELETE /docs/{collection}/{id}     # Delete document
GET    /docs/{collection}          # Query documents
```

### Key-Value Store
```
POST/PUT /kv/{key}     # Set key-value
GET      /kv/{key}     # Get value
DELETE   /kv/{key}     # Delete key
```

### Column Store
```
POST/PUT /columns/{family}/{row}/{column}     # Insert column value
GET      /columns/{family}/{row}/{column}     # Get column value
```

### Graph Store
```
POST /graph/nodes     # Create node
GET  /graph/nodes/{id} # Get node
POST /graph/edges     # Create edge
GET  /graph/edges/{id} # Get edge
```

### Cluster Management
```
GET /cluster/status     # Get cluster status
POST /cluster/nodes     # Add node to cluster
```

## Configuration

The database engine can be configured using environment variables:

- `DB_PORT`: Port for the HTTP API (default: 8080)
- `DB_DATA_DIR`: Directory for data storage (default: ./data)
- `CLUSTER_ENABLED`: Enable clustering (default: false)
- `CLUSTER_PORT`: Port for cluster communication (default: 9090)
- `REPLICATION_FACTOR`: Number of replicas (default: 1)
- `CONSISTENCY_LEVEL`: Consistency level (default: quorum)

## Building and Running

```bash
# Install dependencies
go mod tidy

# Run the database engine
go run main.go
```

Or build and run the binary:
```bash
go build -o multimodel-db
./multimodel-db
```

## Docker Deployment

Create a Dockerfile:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o multimodel-db

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/multimodel-db .
EXPOSE 8080
CMD ["./multimodel-db"]
```

## Kubernetes Deployment

For Kubernetes deployment, you can create a StatefulSet to ensure stable network identities for clustering:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: multimodel-db
spec:
  serviceName: multimodel-db
  replicas: 3
  selector:
    matchLabels:
      app: multimodel-db
  template:
    metadata:
      labels:
        app: multimodel-db
    spec:
      containers:
      - name: multimodel-db
        image: multimodel-db:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        env:
        - name: CLUSTER_ENABLED
          value: "true"
        - name: DB_PORT
          value: "8080"
        - name: CLUSTER_PORT
          value: "9090"
```

## Web Administration Interface

A complete web-based administration interface is available to manage the database engine. The interface provides:

- Dashboard with system health and statistics
- Document store management
- Key-value store management  
- Column store management
- Graph store management
- Cluster monitoring and management
- Real-time performance metrics

To run the web admin interface along with the database engine, use the docker-compose file in the `webapp` directory:

```bash
cd webapp
docker-compose up -d
```

Then access the web admin at `http://localhost:3000` and the database API at `http://localhost:8080`.

## Architecture

The database engine consists of several key components:

1. **Core Engine**: Manages all four data models in-memory
2. **Cluster Component**: Handles node discovery, gossip protocol, and data distribution
3. **API Layer**: HTTP REST interface for client interactions
4. **Storage Layer**: Pluggable storage backends (currently in-memory with persistence planned)

## Use Cases

- Content management systems requiring document flexibility
- Session storage with key-value performance
- Time-series data with column-family efficiency  
- Social networks and recommendation engines with graph relationships
- Microservices architectures requiring polyglot persistence

## Performance Considerations

- All operations are currently in-memory for speed
- Planned persistence layer will use BoltDB for local storage
- Clustering provides horizontal scaling
- Consistent hashing for load distribution