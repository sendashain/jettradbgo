package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// DBClient represents a client to communicate with the database engine
type DBClient struct {
	BaseURL string
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DatabaseInfo holds information about the database
type DatabaseInfo struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	CollectionCount int    `json:"collection_count"`
	Size            int64  `json:"size"`
}

// Document represents a document in the document store
type Document map[string]interface{}

// NewDBClient creates a new client for the database engine
func NewDBClient(baseURL string) *DBClient {
	return &DBClient{BaseURL: baseURL}
}

// makeRequest makes HTTP requests to the database engine
func (c *DBClient) makeRequest(method, endpoint string, payload interface{}) (*Response, error) {
	var req *http.Request
	var err error

	if payload != nil {
		payloadBytes, _ := json.Marshal(payload)
		req, err = http.NewRequest(method, c.BaseURL+endpoint, strings.NewReader(string(payloadBytes)))
	} else {
		req, err = http.NewRequest(method, c.BaseURL+endpoint, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// GetHealth checks the health of the database engine
func (c *DBClient) GetHealth() (*Response, error) {
	return c.makeRequest("GET", "/health", nil)
}

// GetClusterStatus gets the cluster status
func (c *DBClient) GetClusterStatus() (*Response, error) {
	return c.makeRequest("GET", "/cluster/status", nil)
}

// GetCollections gets all collections in the document store
func (c *DBClient) GetCollections() ([]string, error) {
	// Since we don't have a direct API endpoint for this, we'll simulate by listing known collections
	// In a real implementation, this would come from the database engine
	return []string{"users", "products", "orders"}, nil
}

// GetDocuments gets documents from a collection
func (c *DBClient) GetDocuments(collection string) (*Response, error) {
	return c.makeRequest("GET", fmt.Sprintf("/docs/%s", collection), nil)
}

// CreateDocument creates a new document
func (c *DBClient) CreateDocument(collection string, id string, doc Document) (*Response, error) {
	return c.makeRequest("POST", fmt.Sprintf("/docs/%s/%s", collection, id), doc)
}

// UpdateDocument updates an existing document
func (c *DBClient) UpdateDocument(collection string, id string, updates Document) (*Response, error) {
	return c.makeRequest("PUT", fmt.Sprintf("/docs/%s/%s", collection, id), updates)
}

// DeleteDocument deletes a document
func (c *DBClient) DeleteDocument(collection string, id string) (*Response, error) {
	return c.makeRequest("DELETE", fmt.Sprintf("/docs/%s/%s", collection, id), nil)
}

// GetKeyValue gets a key-value pair
func (c *DBClient) GetKeyValue(key string) (*Response, error) {
	return c.makeRequest("GET", fmt.Sprintf("/kv/%s", key), nil)
}

// SetKeyValue sets a key-value pair
func (c *DBClient) SetKeyValue(key string, value interface{}) (*Response, error) {
	return c.makeRequest("POST", fmt.Sprintf("/kv/%s", key), value)
}

// DeleteKey deletes a key
func (c *DBClient) DeleteKey(key string) (*Response, error) {
	return c.makeRequest("DELETE", fmt.Sprintf("/kv/%s", key), nil)
}

// GetColumn gets a column value
func (c *DBClient) GetColumn(family, row, column string) (*Response, error) {
	return c.makeRequest("GET", fmt.Sprintf("/columns/%s/%s/%s", family, row, column), nil)
}

// InsertColumn inserts a column value
func (c *DBClient) InsertColumn(family, row, column string, value interface{}) (*Response, error) {
	return c.makeRequest("POST", fmt.Sprintf("/columns/%s/%s/%s", family, row, column), value)
}

// GetNode gets a graph node
func (c *DBClient) GetNode(id string) (*Response, error) {
	return c.makeRequest("GET", fmt.Sprintf("/graph/nodes/%s", id), nil)
}

// CreateNode creates a graph node
func (c *DBClient) CreateNode(id string, labels []string, props map[string]interface{}) (*Response, error) {
	data := map[string]interface{}{
		"id":    id,
		"labels": labels,
		"props": props,
	}
	return c.makeRequest("POST", "/graph/nodes", data)
}

// GetEdge gets a graph edge
func (c *DBClient) GetEdge(id string) (*Response, error) {
	return c.makeRequest("GET", fmt.Sprintf("/graph/edges/%s", id), nil)
}

// CreateEdge creates a graph edge
func (c *DBClient) CreateEdge(id, from, to, edgeType string, props interface{}) (*Response, error) {
	data := map[string]interface{}{
		"id":   id,
		"from": from,
		"to":   to,
		"type": edgeType,
		"props": props,
	}
	return c.makeRequest("POST", "/graph/edges", data)
}

// Global DB client instance
var dbClient *DBClient

func setupRoutes(router *mux.Router) {
	// Health check
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	
	// Cluster management
	router.HandleFunc("/api/cluster/status", clusterStatusHandler).Methods("GET")
	
	// Document store management
	router.HandleFunc("/api/documents/collections", getCollectionsHandler).Methods("GET")
	router.HandleFunc("/api/documents/{collection}", getDocumentsHandler).Methods("GET")
	router.HandleFunc("/api/documents/{collection}/{id}", createDocumentHandler).Methods("POST")
	router.HandleFunc("/api/documents/{collection}/{id}", updateDocumentHandler).Methods("PUT")
	router.HandleFunc("/api/documents/{collection}/{id}", deleteDocumentHandler).Methods("DELETE")
	
	// Key-value store management
	router.HandleFunc("/api/kv/{key}", getKeyHandler).Methods("GET")
	router.HandleFunc("/api/kv/{key}", setKeyHandler).Methods("POST", "PUT")
	router.HandleFunc("/api/kv/{key}", deleteKeyHandler).Methods("DELETE")
	
	// Column store management
	router.HandleFunc("/api/columns/{family}/{row}/{column}", getColumnHandler).Methods("GET")
	router.HandleFunc("/api/columns/{family}/{row}/{column}", insertColumnHandler).Methods("POST", "PUT")
	
	// Graph store management
	router.HandleFunc("/api/graph/nodes/{id}", getNodeHandler).Methods("GET")
	router.HandleFunc("/api/graph/nodes", createNodeHandler).Methods("POST")
	router.HandleFunc("/api/graph/edges/{id}", getEdgeHandler).Methods("GET")
	router.HandleFunc("/api/graph/edges", createEdgeHandler).Methods("POST")
	
	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./webapp/frontend/dist/")))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := dbClient.GetHealth()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func clusterStatusHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := dbClient.GetClusterStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	collections, err := dbClient.GetCollections()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := Response{
		Success: true,
		Data:    collections,
	}
	json.NewEncoder(w).Encode(response)
}

func getDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	
	resp, err := dbClient.GetDocuments(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	id := vars["id"]
	
	var doc Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.CreateDocument(collection, id, doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func updateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	id := vars["id"]
	
	var updates Document
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.UpdateDocument(collection, id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	id := vars["id"]
	
	resp, err := dbClient.DeleteDocument(collection, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func getKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	resp, err := dbClient.GetKeyValue(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func setKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	var value interface{}
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.SetKeyValue(key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func deleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	resp, err := dbClient.DeleteKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func getColumnHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	family := vars["family"]
	row := vars["row"]
	column := vars["column"]
	
	resp, err := dbClient.GetColumn(family, row, column)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func insertColumnHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	family := vars["family"]
	row := vars["row"]
	column := vars["column"]
	
	var value interface{}
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.InsertColumn(family, row, column, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func getNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	resp, err := dbClient.GetNode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func createNodeHandler(w http.ResponseWriter, r *http.Request) {
	var nodeData struct {
		ID     string                 `json:"id"`
		Labels []string               `json:"labels"`
		Props  map[string]interface{} `json:"props"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&nodeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.CreateNode(nodeData.ID, nodeData.Labels, nodeData.Props)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func getEdgeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	resp, err := dbClient.GetEdge(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func createEdgeHandler(w http.ResponseWriter, r *http.Request) {
	var edgeData struct {
		ID   string      `json:"id"`
		From string      `json:"from"`
		To   string      `json:"to"`
		Type string      `json:"type"`
		Props interface{} `json:"props"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&edgeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	resp, err := dbClient.CreateEdge(edgeData.ID, edgeData.From, edgeData.To, edgeData.Type, edgeData.Props)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "http://localhost:8080" // Default URL
	}
	
	dbClient = NewDBClient(dbURL)
	
	router := mux.NewRouter()
	setupRoutes(router)
	
	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	handler := c.Handler(router)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	fmt.Printf("Web Admin starting on port %s\n", port)
	fmt.Printf("Connecting to database engine at: %s\n", dbURL)
	
	log.Fatal(http.ListenAndServe(":"+port, handler))
}