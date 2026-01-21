package database

import (
	"encoding/json"
	"fmt"
	"sync"

	"multimodel-db-engine/internal/config"
)

// Document represents a document in the document store
type Document map[string]interface{}

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value interface{}
}

// ColumnFamily represents a column family in the column store
type ColumnFamily map[string]map[string]interface{}

// GraphNode represents a node in the graph store
type GraphNode struct {
	ID     string                 `json:"id"`
	Labels []string               `json:"labels"`
	Props  map[string]interface{} `json:"props"`
}

// GraphEdge represents an edge in the graph store
type GraphEdge struct {
	ID     string      `json:"id"`
	From   string      `json:"from"`
	To     string      `json:"to"`
	Type   string      `json:"type"`
	Props  interface{} `json:"props"`
}

// MultiModelDatabase represents the multi-model database engine
type MultiModelDatabase struct {
	config *config.Config
	
	// Document store
	documents map[string]Document
	docMutex  sync.RWMutex
	
	// Key-value store
	keyValues map[string]interface{}
	kvMutex   sync.RWMutex
	
	// Column store
	columnFamilies map[string]ColumnFamily
	colMutex       sync.RWMutex
	
	// Graph store
	graphNodes map[string]*GraphNode
	graphEdges map[string]*GraphEdge
	graphMutex sync.RWMutex
	
	// Distributed cluster components
	Cluster *Cluster  // Public field to access cluster from other packages
}

// NewMultiModelDatabase creates a new instance of the multi-model database
func NewMultiModelDatabase(cfg *config.Config) *MultiModelDatabase {
	db := &MultiModelDatabase{
		config:         cfg,
		documents:      make(map[string]Document),
		keyValues:      make(map[string]interface{}),
		columnFamilies: make(map[string]ColumnFamily),
		graphNodes:     make(map[string]*GraphNode),
		graphEdges:     make(map[string]*GraphEdge),
	}
	
	// Initialize cluster if enabled
	if cfg.ClusterEnabled {
		db.Cluster = NewCluster(cfg)
	}
	
	return db
}

// Document Store Operations
func (db *MultiModelDatabase) InsertDocument(collection, id string, doc Document) error {
	db.docMutex.Lock()
	defer db.docMutex.Unlock()
	
	if _, exists := db.documents[collection+"."+id]; exists {
		return fmt.Errorf("document with id %s already exists in collection %s", id, collection)
	}
	
	db.documents[collection+"."+id] = doc
	return nil
}

func (db *MultiModelDatabase) GetDocument(collection, id string) (Document, error) {
	db.docMutex.RLock()
	defer db.docMutex.RUnlock()
	
	doc, exists := db.documents[collection+"."+id]
	if !exists {
		return nil, fmt.Errorf("document with id %s not found in collection %s", id, collection)
	}
	
	return doc, nil
}

func (db *MultiModelDatabase) UpdateDocument(collection, id string, updates Document) error {
	db.docMutex.Lock()
	defer db.docMutex.Unlock()
	
	key := collection + "." + id
	doc, exists := db.documents[key]
	if !exists {
		return fmt.Errorf("document with id %s not found in collection %s", id, collection)
	}
	
	// Merge updates into existing document
	for k, v := range updates {
		doc[k] = v
	}
	
	db.documents[key] = doc
	return nil
}

func (db *MultiModelDatabase) DeleteDocument(collection, id string) error {
	db.docMutex.Lock()
	defer db.docMutex.Unlock()
	
	key := collection + "." + id
	if _, exists := db.documents[key]; !exists {
		return fmt.Errorf("document with id %s not found in collection %s", id, collection)
	}
	
	delete(db.documents, key)
	return nil
}

// Key-Value Store Operations
func (db *MultiModelDatabase) SetKeyValue(key string, value interface{}) error {
	db.kvMutex.Lock()
	defer db.kvMutex.Unlock()
	
	db.keyValues[key] = value
	return nil
}

func (db *MultiModelDatabase) GetKeyValue(key string) (interface{}, error) {
	db.kvMutex.RLock()
	defer db.kvMutex.RUnlock()
	
	value, exists := db.keyValues[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found", key)
	}
	
	return value, nil
}

func (db *MultiModelDatabase) DeleteKey(key string) error {
	db.kvMutex.Lock()
	defer db.kvMutex.Unlock()
	
	if _, exists := db.keyValues[key]; !exists {
		return fmt.Errorf("key %s not found", key)
	}
	
	delete(db.keyValues, key)
	return nil
}

// Column Store Operations
func (db *MultiModelDatabase) InsertColumn(columnFamily, rowKey, columnName string, value interface{}) error {
	db.colMutex.Lock()
	defer db.colMutex.Unlock()
	
	cf, exists := db.columnFamilies[columnFamily]
	if !exists {
		cf = make(map[string]map[string]interface{})
		db.columnFamilies[columnFamily] = cf
	}
	
	row, exists := cf[rowKey]
	if !exists {
		row = make(map[string]interface{})
		cf[rowKey] = row
	}
	
	row[columnName] = value
	return nil
}

func (db *MultiModelDatabase) GetColumn(columnFamily, rowKey, columnName string) (interface{}, error) {
	db.colMutex.RLock()
	defer db.colMutex.RUnlock()
	
	cf, exists := db.columnFamilies[columnFamily]
	if !exists {
		return nil, fmt.Errorf("column family %s not found", columnFamily)
	}
	
	row, exists := cf[rowKey]
	if !exists {
		return nil, fmt.Errorf("row %s not found in column family %s", rowKey, columnFamily)
	}
	
	value, exists := row[columnName]
	if !exists {
		return nil, fmt.Errorf("column %s not found in row %s of column family %s", columnName, rowKey, columnFamily)
	}
	
	return value, nil
}

// Graph Store Operations
func (db *MultiModelDatabase) CreateNode(id string, labels []string, props map[string]interface{}) error {
	db.graphMutex.Lock()
	defer db.graphMutex.Unlock()
	
	if _, exists := db.graphNodes[id]; exists {
		return fmt.Errorf("node with id %s already exists", id)
	}
	
	node := &GraphNode{
		ID:     id,
		Labels: labels,
		Props:  props,
	}
	
	db.graphNodes[id] = node
	return nil
}

func (db *MultiModelDatabase) GetNode(id string) (*GraphNode, error) {
	db.graphMutex.RLock()
	defer db.graphMutex.RUnlock()
	
	node, exists := db.graphNodes[id]
	if !exists {
		return nil, fmt.Errorf("node with id %s not found", id)
	}
	
	return node, nil
}

func (db *MultiModelDatabase) CreateEdge(id, from, to, edgeType string, props interface{}) error {
	db.graphMutex.Lock()
	defer db.graphMutex.Unlock()
	
	if _, exists := db.graphEdges[id]; exists {
		return fmt.Errorf("edge with id %s already exists", id)
	}
	
	// Check if nodes exist
	if _, exists := db.graphNodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}
	if _, exists := db.graphNodes[to]; !exists {
		return fmt.Errorf("target node %s does not exist", to)
	}
	
	edge := &GraphEdge{
		ID:   id,
		From: from,
		To:   to,
		Type: edgeType,
		Props: props,
	}
	
	db.graphEdges[id] = edge
	return nil
}

func (db *MultiModelDatabase) GetEdge(id string) (*GraphEdge, error) {
	db.graphMutex.RLock()
	defer db.graphMutex.RUnlock()
	
	edge, exists := db.graphEdges[id]
	if !exists {
		return nil, fmt.Errorf("edge with id %s not found", id)
	}
	
	return edge, nil
}

// Query methods for each model
func (db *MultiModelDatabase) QueryDocuments(collection string, filter map[string]interface{}) ([]Document, error) {
	db.docMutex.RLock()
	defer db.docMutex.RUnlock()
	
	var results []Document
	
	for key, doc := range db.documents {
		if collection == "" || len(collection) <= len(key) && key[:len(collection)] == collection {
			// Apply filters
			matches := true
			for field, expectedValue := range filter {
				if actualValue, exists := doc[field]; !exists || actualValue != expectedValue {
					matches = false
					break
				}
			}
			
			if matches {
				results = append(results, doc)
			}
		}
	}
	
	return results, nil
}

// JSON serialization helper
func (db *MultiModelDatabase) ToJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}