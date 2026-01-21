package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"multimodel-db-engine/internal/config"
)

// Node represents a node in the distributed cluster
type Node struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Status   string `json:"status"` // active, inactive, joining, leaving
	LastSeen int64  `json:"last_seen"`
}

// Cluster represents the distributed cluster component
type Cluster struct {
	selfNode    *Node
	nodes       map[string]*Node
	nodesMutex  sync.RWMutex
	config      *config.Config
	httpClient  *http.Client
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

// NewCluster creates a new cluster instance
func NewCluster(cfg *config.Config) *Cluster {
	ctx, cancel := context.WithCancel(context.Background())
	
	cluster := &Cluster{
		selfNode: &Node{
			ID:      generateNodeID(),
			Address: "localhost", // In production, get actual IP
			Port:    cfg.ClusterPort,
			Status:  "active",
		},
		nodes:      make(map[string]*Node),
		config:     cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		ctx:        ctx,
		cancelFunc: cancel,
	}
	
	// Add self to the cluster
	cluster.nodes[cluster.selfNode.ID] = cluster.selfNode
	
	// Start cluster maintenance routines
	go cluster.startHeartbeat()
	go cluster.startGossipProtocol()
	
	return cluster
}

// generateNodeID creates a unique identifier for a node
func generateNodeID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("node-%d-%d", time.Now().Unix(), rand.Intn(10000))
}

// Join attempts to join an existing cluster
func (c *Cluster) Join(seedAddress string) error {
	url := fmt.Sprintf("http://%s/cluster/join", seedAddress)
	
	reqBody, _ := json.Marshal(c.selfNode)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to join cluster: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("join request failed with status: %d", resp.StatusCode)
	}
	
	return nil
}

// AddNode adds a new node to the cluster
func (c *Cluster) AddNode(node *Node) {
	c.nodesMutex.Lock()
	defer c.nodesMutex.Unlock()
	
	node.LastSeen = time.Now().Unix()
	c.nodes[node.ID] = node
	
	log.Printf("Added node %s to cluster", node.ID)
}

// RemoveNode removes a node from the cluster
func (c *Cluster) RemoveNode(nodeID string) {
	c.nodesMutex.Lock()
	defer c.nodesMutex.Unlock()
	
	if node, exists := c.nodes[nodeID]; exists {
		node.Status = "inactive"
		log.Printf("Removed node %s from cluster", nodeID)
	}
}

// GetActiveNodes returns all active nodes in the cluster
func (c *Cluster) GetActiveNodes() []*Node {
	c.nodesMutex.RLock()
	defer c.nodesMutex.RUnlock()
	
	var activeNodes []*Node
	for _, node := range c.nodes {
		if node.Status == "active" {
			activeNodes = append(activeNodes, node)
		}
	}
	
	return activeNodes
}

// GetNode returns a specific node by ID
func (c *Cluster) GetNode(nodeID string) (*Node, bool) {
	c.nodesMutex.RLock()
	defer c.nodesMutex.RUnlock()
	
	node, exists := c.nodes[nodeID]
	return node, exists
}

// startHeartbeat starts the heartbeat mechanism to detect node failures
func (c *Cluster) startHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.performHeartbeat()
		}
	}
}

// performHeartbeat sends heartbeats to other nodes and updates their status
func (c *Cluster) performHeartbeat() {
	c.nodesMutex.RLock()
	nodes := make([]*Node, 0, len(c.nodes))
	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}
	c.nodesMutex.RUnlock()
	
	for _, node := range nodes {
		if node.ID == c.selfNode.ID {
			continue // Skip self
		}
		
		alive := c.pingNode(node)
		c.updateNodeStatus(node.ID, alive)
	}
}

// pingNode checks if a node is alive
func (c *Cluster) pingNode(node *Node) bool {
	url := fmt.Sprintf("http://%s:%s/health", node.Address, node.Port)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

// updateNodeStatus updates the status of a node based on heartbeat results
func (c *Cluster) updateNodeStatus(nodeID string, alive bool) {
	c.nodesMutex.Lock()
	defer c.nodesMutex.Unlock()
	
	if node, exists := c.nodes[nodeID]; exists {
		if alive {
			node.Status = "active"
			node.LastSeen = time.Now().Unix()
		} else {
			node.Status = "inactive"
		}
	}
}

// startGossipProtocol starts the gossip protocol for sharing cluster information
func (c *Cluster) startGossipProtocol() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.runGossipRound()
		}
	}
}

// runGossipRound performs a round of gossip with random nodes
func (c *Cluster) runGossipRound() {
	activeNodes := c.GetActiveNodes()
	
	if len(activeNodes) <= 1 {
		return // No other nodes to gossip with
	}
	
	// Select a random subset of nodes to gossip with
	nodesToGossip := 2
	if len(activeNodes) < 2 {
		nodesToGossip = len(activeNodes)
	}
	
	for i := 0; i < nodesToGossip; i++ {
		randomIndex := rand.Intn(len(activeNodes))
		targetNode := activeNodes[randomIndex]
		
		if targetNode.ID == c.selfNode.ID {
			continue
		}
		
		c.exchangeGossip(targetNode)
	}
}

// exchangeGossip exchanges cluster membership information with another node
func (c *Cluster) exchangeGossip(node *Node) {
	url := fmt.Sprintf("http://%s:%s/cluster/gossip", node.Address, node.Port)
	
	localNodes := c.GetActiveNodes()
	reqBody, _ := json.Marshal(localNodes)
	
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Failed to gossip with node %s: %v", node.ID, err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Gossip request to node %s failed with status: %d", node.ID, resp.StatusCode)
		return
	}
	
	// Process response with other node's membership info
	// This would normally merge the membership lists
}

// GetPartitionForKey determines which node should handle a given key
func (c *Cluster) GetPartitionForKey(key string) *Node {
	activeNodes := c.GetActiveNodes()
	if len(activeNodes) == 0 {
		return nil
	}
	
	// Simple hash-based partitioning
	hash := 0
	for _, r := range key {
		hash = (hash << 5) - hash + int(r)
	}
	
	index := hash % len(activeNodes)
	if index < 0 {
		index = -index
	}
	
	return activeNodes[index]
}

// ReplicateData replicates data to other nodes based on replication factor
func (c *Cluster) ReplicateData(key string, value interface{}) error {
	replicationFactor := c.config.ReplicationFactor
	if replicationFactor <= 1 {
		return nil // No replication needed
	}
	
	activeNodes := c.GetActiveNodes()
	if len(activeNodes) < replicationFactor {
		return fmt.Errorf("not enough active nodes for replication factor %d", replicationFactor)
	}
	
	// Find nodes responsible for this key
	nodes := make([]*Node, 0, replicationFactor)
	nodes = append(nodes, c.GetPartitionForKey(key)) // Primary
	
	// Add additional replicas
	for i := 1; i < replicationFactor && i < len(activeNodes); i++ {
		// Add next node in ring
		primaryIdx := 0
		for j, node := range activeNodes {
			if node.ID == nodes[0].ID {
				primaryIdx = j
				break
			}
		}
		
		replicaIdx := (primaryIdx + i) % len(activeNodes)
		nodes = append(nodes, activeNodes[replicaIdx])
	}
	
	// Replicate to selected nodes
	for _, node := range nodes {
		if node.ID == c.selfNode.ID {
			continue // Skip self, we already have the data
		}
		
		err := c.replicateToNode(node, key, value)
		if err != nil {
			log.Printf("Failed to replicate key %s to node %s: %v", key, node.ID, err)
		}
	}
	
	return nil
}

// replicateToNode sends data to a specific node for replication
func (c *Cluster) replicateToNode(node *Node, key string, value interface{}) error {
	url := fmt.Sprintf("http://%s:%s/data/replicate", node.Address, node.Port)
	
	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	
	reqBody, _ := json.Marshal(data)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to replicate to node %s: %w", node.ID, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("replication request failed with status: %d", resp.StatusCode)
	}
	
	return nil
}

// Close shuts down the cluster component gracefully
func (c *Cluster) Close() {
	c.cancelFunc()
}