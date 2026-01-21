// Multi-Model Database Admin Application
class DatabaseAdminApp {
    constructor() {
        this.currentView = 'dashboard';
        this.apiBaseUrl = '/api'; // Proxy to backend
        
        // Initialize the application
        this.init();
    }
    
    init() {
        // Set up event listeners
        this.setupEventListeners();
        
        // Load initial data
        this.loadDashboardData();
        
        console.log('Database Admin App initialized');
    }
    
    setupEventListeners() {
        // Navigation links
        document.getElementById('dashboard-link').addEventListener('click', () => this.switchView('dashboard'));
        document.getElementById('documents-link').addEventListener('click', () => this.switchView('documents'));
        document.getElementById('keyvalue-link').addEventListener('click', () => this.switchView('keyvalue'));
        document.getElementById('columns-link').addEventListener('click', () => this.switchView('columns'));
        document.getElementById('graph-link').addEventListener('click', () => this.switchView('graph'));
        document.getElementById('cluster-link').addEventListener('click', () => this.switchView('cluster'));
        
        // Refresh button
        document.getElementById('refresh-btn').addEventListener('click', () => this.refreshCurrentView());
        
        // Action buttons
        document.getElementById('add-document-btn').addEventListener('click', () => this.openAddDocumentModal());
        document.getElementById('add-kv-btn').addEventListener('click', () => this.openAddKVModal());
        document.getElementById('add-column-btn').addEventListener('click', () => this.openAddColumnModal());
        document.getElementById('add-node-btn').addEventListener('click', () => this.openAddNodeModal());
        document.getElementById('add-node-btn').addEventListener('click', () => this.openAddNodeModal()); // Cluster view also has this ID
        document.getElementById('add-edge-btn').addEventListener('click', () => this.openAddEdgeModal());
    }
    
    switchView(viewName) {
        // Hide all views
        const views = ['dashboard', 'documents', 'keyvalue', 'columns', 'graph', 'cluster'];
        views.forEach(view => {
            document.getElementById(`${view}-view`).style.display = 'none';
        });
        
        // Remove active class from all nav links
        views.forEach(view => {
            const link = document.getElementById(`${view}-link`);
            if (link) link.classList.remove('active');
        });
        
        // Show selected view
        document.getElementById(`${viewName}-view`).style.display = 'block';
        
        // Update active nav link
        const activeLink = document.getElementById(`${viewName}-link`);
        if (activeLink) activeLink.classList.add('active');
        
        // Update page title
        document.getElementById('page-title').textContent = 
            viewName.charAt(0).toUpperCase() + viewName.slice(1).replace(/([A-Z])/g, ' $1');
        
        // Load data for the specific view
        this.loadViewData(viewName);
        
        this.currentView = viewName;
    }
    
    loadViewData(viewName) {
        switch(viewName) {
            case 'dashboard':
                this.loadDashboardData();
                break;
            case 'documents':
                this.loadDocumentsData();
                break;
            case 'keyvalue':
                this.loadKeyValueData();
                break;
            case 'columns':
                this.loadColumnsData();
                break;
            case 'graph':
                this.loadGraphData();
                break;
            case 'cluster':
                this.loadClusterData();
                break;
        }
    }
    
    refreshCurrentView() {
        this.loadViewData(this.currentView);
    }
    
    async loadDashboardData() {
        try {
            // Simulate API call to get dashboard stats
            const stats = {
                totalDbs: 4,
                totalNodes: 3,
                totalCollections: 12,
                status: 'Active'
            };
            
            document.getElementById('total-dbs').textContent = stats.totalDbs;
            document.getElementById('total-nodes').textContent = stats.totalNodes;
            document.getElementById('total-collections').textContent = stats.totalCollections;
            document.getElementById('status').textContent = stats.status;
            
            console.log('Dashboard data loaded');
        } catch (error) {
            console.error('Error loading dashboard data:', error);
        }
    }
    
    async loadDocumentsData() {
        try {
            // Simulate API call to get documents
            const response = await fetch(`${this.apiBaseUrl}/documents/users`);
            const result = await response.json();
            
            if (result.success) {
                this.renderDocuments(result.data);
            } else {
                console.error('Failed to load documents:', result.error);
            }
        } catch (error) {
            console.error('Error loading documents:', error);
        }
    }
    
    renderDocuments(documents) {
        const tbody = document.getElementById('documents-table-body');
        tbody.innerHTML = '';
        
        if (!documents || documents.length === 0) {
            tbody.innerHTML = '<tr><td colspan="3" class="text-center">No documents found</td></tr>';
            return;
        }
        
        documents.forEach(doc => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${doc.id || 'N/A'}</td>
                <td>${JSON.stringify(doc, null, 2)}</td>
                <td>
                    <button class="btn btn-sm btn-outline-primary me-1" onclick="app.editDocument('${doc.id}')">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="app.deleteDocument('${doc.id}')">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });
    }
    
    async loadKeyValueData() {
        try {
            // Simulate API call to get key-value pairs
            console.log('Loading key-value data...');
        } catch (error) {
            console.error('Error loading key-value data:', error);
        }
    }
    
    async loadColumnsData() {
        try {
            // Simulate API call to get column data
            console.log('Loading column data...');
        } catch (error) {
            console.error('Error loading column data:', error);
        }
    }
    
    async loadGraphData() {
        try {
            // Simulate API call to get graph data
            console.log('Loading graph data...');
        } catch (error) {
            console.error('Error loading graph data:', error);
        }
    }
    
    async loadClusterData() {
        try {
            const response = await fetch(`${this.apiBaseUrl}/cluster/status`);
            const result = await response.json();
            
            if (result.success && result.data) {
                this.renderClusterNodes(result.data.nodes || []);
            } else {
                console.error('Failed to load cluster data:', result.error);
            }
        } catch (error) {
            console.error('Error loading cluster data:', error);
        }
    }
    
    renderClusterNodes(nodes) {
        const tbody = document.getElementById('cluster-table-body');
        tbody.innerHTML = '';
        
        if (!nodes || nodes.length === 0) {
            tbody.innerHTML = '<tr><td colspan="7" class="text-center">No cluster nodes found</td></tr>';
            return;
        }
        
        nodes.forEach(node => {
            const statusClass = node.status === 'active' ? 'status-active' : 'status-inactive';
            const statusText = node.status.charAt(0).toUpperCase() + node.status.slice(1);
            
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${node.id || 'N/A'}</td>
                <td>${node.address || 'N/A'}</td>
                <td>${node.port || 'N/A'}</td>
                <td><span class="status-indicator ${statusClass}"></span> ${statusText}</td>
                <td>${node.lastSeen || 'N/A'}</td>
                <td>${node.role || 'Node'}</td>
                <td>
                    <button class="btn btn-sm btn-outline-primary me-1" onclick="app.viewNode('${node.id}')">
                        <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-warning me-1" onclick="app.toggleNode('${node.id}')">
                        <i class="fas fa-power-off"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="app.removeNode('${node.id}')">
                        <i class="fas fa-times"></i>
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });
    }
    
    openAddDocumentModal() {
        alert('Add Document functionality would open a modal here');
    }
    
    openAddKVModal() {
        alert('Add Key-Value functionality would open a modal here');
    }
    
    openAddColumnModal() {
        alert('Add Column functionality would open a modal here');
    }
    
    openAddNodeModal() {
        alert('Add Node functionality would open a modal here');
    }
    
    openAddEdgeModal() {
        alert('Add Edge functionality would open a modal here');
    }
    
    editDocument(id) {
        alert(`Edit document with ID: ${id}`);
    }
    
    deleteDocument(id) {
        if (confirm(`Are you sure you want to delete document with ID: ${id}?`)) {
            console.log(`Deleting document with ID: ${id}`);
        }
    }
    
    viewNode(id) {
        alert(`Viewing details for node: ${id}`);
    }
    
    toggleNode(id) {
        console.log(`Toggling node: ${id}`);
    }
    
    removeNode(id) {
        if (confirm(`Are you sure you want to remove node: ${id}?`)) {
            console.log(`Removing node: ${id}`);
        }
    }
    
    // API helper methods
    async apiCall(endpoint, method = 'GET', data = null) {
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
            }
        };
        
        if (data) {
            options.body = JSON.stringify(data);
        }
        
        try {
            const response = await fetch(`${this.apiBaseUrl}${endpoint}`, options);
            return await response.json();
        } catch (error) {
            console.error('API call failed:', error);
            throw error;
        }
    }
}

// Initialize the app when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new DatabaseAdminApp();
});

// Export the app instance globally so it can be accessed by inline event handlers
window.DatabaseAdminApp = DatabaseAdminApp;