package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"multimodel-db-engine/internal/database"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SetupRoutes configures all API routes
func SetupRoutes(router *mux.Router, db *database.MultiModelDatabase) {
	// Health check endpoint
	router.HandleFunc("/health", healthHandler).Methods("GET")
	
	// Document store endpoints
	router.HandleFunc("/docs/{collection}/{id}", createDocumentHandler(db)).Methods("POST")
	router.HandleFunc("/docs/{collection}/{id}", getDocumentHandler(db)).Methods("GET")
	router.HandleFunc("/docs/{collection}/{id}", updateDocumentHandler(db)).Methods("PUT")
	router.HandleFunc("/docs/{collection}/{id}", deleteDocumentHandler(db)).Methods("DELETE")
	router.HandleFunc("/docs/{collection}", queryDocumentsHandler(db)).Methods("GET")
	
	// Key-value store endpoints
	router.HandleFunc("/kv/{key}", setKeyValueHandler(db)).Methods("POST", "PUT")
	router.HandleFunc("/kv/{key}", getKeyValueHandler(db)).Methods("GET")
	router.HandleFunc("/kv/{key}", deleteKeyHandler(db)).Methods("DELETE")
	
	// Column store endpoints
	router.HandleFunc("/columns/{family}/{row}/{column}", insertColumnHandler(db)).Methods("POST", "PUT")
	router.HandleFunc("/columns/{family}/{row}/{column}", getColumnHandler(db)).Methods("GET")
	
	// Graph store endpoints
	router.HandleFunc("/graph/nodes", createNodeHandler(db)).Methods("POST")
	router.HandleFunc("/graph/nodes/{id}", getNodeHandler(db)).Methods("GET")
	router.HandleFunc("/graph/edges", createEdgeHandler(db)).Methods("POST")
	router.HandleFunc("/graph/edges/{id}", getEdgeHandler(db)).Methods("GET")
	
	// Cluster endpoints
	router.HandleFunc("/cluster/status", clusterStatusHandler(db)).Methods("GET")
	router.HandleFunc("/cluster/nodes", addNodeHandler(db)).Methods("POST")
	
	// Catch-all for undefined routes
	router.PathPrefix("/").HandlerFunc(notFoundHandler)
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helper function to read JSON body
func readJSONBody(r *http.Request, dst interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	
	return json.Unmarshal(body, dst)
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Multi-Model Database Engine is running",
	})
}

// Document Store Handlers
func createDocumentHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collection := vars["collection"]
		id := vars["id"]
		
		var doc database.Document
		if err := readJSONBody(r, &doc); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.InsertDocument(collection, id, doc); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusCreated, Response{
			Success: true,
			Message: "Document created successfully",
			Data:    doc,
		})
	}
}

func getDocumentHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collection := vars["collection"]
		id := vars["id"]
		
		doc, err := db.GetDocument(collection, id)
		if err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    doc,
		})
	}
}

func updateDocumentHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collection := vars["collection"]
		id := vars["id"]
		
		var updates database.Document
		if err := readJSONBody(r, &updates); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.UpdateDocument(collection, id, updates); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Document updated successfully",
		})
	}
}

func deleteDocumentHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collection := vars["collection"]
		id := vars["id"]
		
		if err := db.DeleteDocument(collection, id); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Document deleted successfully",
		})
	}
}

func queryDocumentsHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collection := vars["collection"]
		
		// Parse query parameters as filters
		filters := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				// For simplicity, take the first value
				filters[key] = values[0]
			}
		}
		
		docs, err := db.QueryDocuments(collection, filters)
		if err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    docs,
		})
	}
}

// Key-Value Store Handlers
func setKeyValueHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		
		var value interface{}
		if err := readJSONBody(r, &value); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.SetKeyValue(key, value); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Key-value pair set successfully",
		})
	}
}

func getKeyValueHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		
		value, err := db.GetKeyValue(key)
		if err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    value,
		})
	}
}

func deleteKeyHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		
		if err := db.DeleteKey(key); err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Key deleted successfully",
		})
	}
}

// Column Store Handlers
func insertColumnHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		family := vars["family"]
		row := vars["row"]
		column := vars["column"]
		
		var value interface{}
		if err := readJSONBody(r, &value); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.InsertColumn(family, row, column, value); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Column value inserted successfully",
		})
	}
}

func getColumnHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		family := vars["family"]
		row := vars["row"]
		column := vars["column"]
		
		value, err := db.GetColumn(family, row, column)
		if err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    value,
		})
	}
}

// Graph Store Handlers
func createNodeHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var nodeData struct {
			ID     string                 `json:"id"`
			Labels []string               `json:"labels"`
			Props  map[string]interface{} `json:"props"`
		}
		
		if err := readJSONBody(r, &nodeData); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.CreateNode(nodeData.ID, nodeData.Labels, nodeData.Props); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusCreated, Response{
			Success: true,
			Message: "Node created successfully",
		})
	}
}

func getNodeHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		
		node, err := db.GetNode(id)
		if err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    node,
		})
	}
}

func createEdgeHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var edgeData struct {
			ID   string      `json:"id"`
			From string      `json:"from"`
			To   string      `json:"to"`
			Type string      `json:"type"`
			Props interface{} `json:"props"`
		}
		
		if err := readJSONBody(r, &edgeData); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		if err := db.CreateEdge(edgeData.ID, edgeData.From, edgeData.To, edgeData.Type, edgeData.Props); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusCreated, Response{
			Success: true,
			Message: "Edge created successfully",
		})
	}
}

func getEdgeHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		
		edge, err := db.GetEdge(id)
		if err != nil {
			sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    edge,
		})
	}
}

// Cluster Handlers
func clusterStatusHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.Cluster == nil {
			sendJSONResponse(w, http.StatusOK, Response{
				Success: true,
				Message: "Clustering is disabled",
				Data:    map[string]interface{}{"enabled": false},
			})
			return
		}
		
		nodes := db.Cluster.GetActiveNodes()
		nodeInfo := make([]map[string]interface{}, len(nodes))
		for i, node := range nodes {
			nodeInfo[i] = map[string]interface{}{
				"id":       node.ID,
				"address":  node.Address,
				"port":     node.Port,
				"status":   node.Status,
				"lastSeen": node.LastSeen,
			}
		}
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    map[string]interface{}{"nodes": nodeInfo, "enabled": true},
		})
	}
}

func addNodeHandler(db *database.MultiModelDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.Cluster == nil {
			sendJSONResponse(w, http.StatusServiceUnavailable, Response{
				Success: false,
				Error:   "Clustering is not enabled",
			})
			return
		}
		
		var nodeData struct {
			Address string `json:"address"`
			Port    string `json:"port"`
		}
		
		if err := readJSONBody(r, &nodeData); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid JSON in request body",
			})
			return
		}
		
		node := &database.Node{
			ID:      fmt.Sprintf("node-%d", len(db.Cluster.GetActiveNodes())),
			Address: nodeData.Address,
			Port:    nodeData.Port,
			Status:  "active",
		}
		
		db.Cluster.AddNode(node)
		
		sendJSONResponse(w, http.StatusOK, Response{
			Success: true,
			Message: "Node added to cluster",
			Data:    node,
		})
	}
}

// Not Found Handler
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Endpoint not found",
	})
}