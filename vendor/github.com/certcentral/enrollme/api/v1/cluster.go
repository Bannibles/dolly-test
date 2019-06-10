package v1

import (
	"time"
)

// NodeInfo provides information about cluster node
type NodeInfo struct {
	// ID for this node.
	ID string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name is the human-readable name of the node. If the node is not started, the name will be an empty string.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// ListenPeerURLs is the list of URLs the member exposes to the cluster for communication,
	// the node accepts incoming requests from its peers on the specified scheme://IP:port combinations.
	ListenPeerURLs []string `protobuf:"bytes,3,rep,name=listen_peer_urls" json:"listen_peer_urls,omitempty"`
}

// GetID returns node's ID
func (m *NodeInfo) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// GetName returns node's Name
func (m *NodeInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// GetListenPeerURLs return node's peer Urls
func (m *NodeInfo) GetListenPeerURLs() []string {
	if m != nil {
		return m.ListenPeerURLs
	}
	return nil
}

// NodeStatus provides response about current node
type NodeStatus struct {
	// NodeID is the node ID in the cluster.
	NodeID string `protobuf:"bytes,1,opt,name=node_id,proto3" json:"node_id,omitempty"`
	// NodeName is the human-readable name of the cluster node.
	NodeName string `protobuf:"bytes,2,name=nodename,proto3" json:"nodename"`
	// HostName is operating system's host name
	HostName string `protobuf:"bytes,3,name=hostname,proto3" json:"hostname"`
	// Port is the listening port
	Port string `protobuf:"bytes,4,name=port,proto3" json:"port"`
	// StartedAt is the time when the server started
	StartedAt time.Time `protobuf:"bytes,5,name=started_at,proto3" json:"started_at"`
	// Version is the app build version
	Version string `protobuf:"bytes,6,name=version,proto3" json:"version"`
	// LeaseID specifies lease from original hearbeat.
	// A node should obtain it on start and use for consecutive call
	LeaseID int64 `protobuf:"bytes,7,name=lease,proto3" json:"lease"`
	// TTL is the time to live for the lease
	TTL int64 `json:"ttl"`

	// HeartbeatAt should be removed after the client is upgraded
	HeartbeatAt *time.Time `json:"heartbeat_at,omitempty"`
	// IsLeader should be removed after the client is upgraded
	IsLeader *bool `json:"leader,omitempty"`
}

// ClusterResponse is response for /v1/cluster
type ClusterResponse struct {
	// Cluster is a list of nodes
	Cluster []*NodeStatus `json:"cluster"`
	// LeaderID is the leader ID for this cluster
	LeaderID string `json:"leader_id,omitempty"`
	// ServiceName specifies the name of the service
	ServiceName string `protobuf:"bytes,3,name=service,proto3" json:"service"`
}

// HeartbeatRequest specifies POST request to heartbeat a node status
type HeartbeatRequest struct {
	// Status specifies a status of the node in the cluster
	Status *NodeStatus `protobuf:"bytes,1,name=status,proto3" json:"status"`
	// LeaseID specifies lease from original hearbeat.
	// A node should obtain it on start and use for consecutive call
	LeaseID int64 `protobuf:"bytes,2,name=lease,proto3" json:"lease"`
	// ServiceName specifies the name of the service
	ServiceName string `protobuf:"bytes,3,name=service,proto3" json:"service"`
}

// HeartbeatResponse is response for HeartbeatRequest
type HeartbeatResponse struct {
	// Cluster is a list of nodes
	Cluster []*NodeStatus `protobuf:"bytes,1,rep,name=cluster" json:"cluster"`
	// LeaderID is the leader ID for this cluster
	LeaderID string `protobuf:"bytes,2,name=leader_id" json:"leader_id,omitempty"`
	// LeaseID specifies lease from original hearbeat.
	// A node should obtain it on start and use for consecutive call
	LeaseID int64 `protobuf:"bytes,3,name=lease,proto3" json:"lease"`
	// TTL is the time to live for the lease
	TTL int64 `protobuf:"bytes,4,name=ttl,proto3" json:"ttl"`
	// ServiceName specifies the name of the service
	ServiceName string `protobuf:"bytes,5,name=service,proto3" json:"service"`
}

// NodesResponse returns peers info in the cluster
type NodesResponse struct {
	// Peers is a list of all peers in the cluster.
	Peers []*NodeInfo `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
}

// NodeAddRequest is the structure to represent the request to add a member
type NodeAddRequest struct {
	// Name is the human-readable name of the node.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// ListenPeerURLs is the list of URLs the node exposes to the cluster for communication.
	ListenPeerURLs []string `protobuf:"bytes,2,rep,name=listen_peer_urls" json:"listen_peer_urls,omitempty"`
	// Token is the cluster token that should match the running cluster
	Token string `protobuf:"bytes,3,rep,name=token" json:"token,omitempty"`
}

// NodeAddResponse is the structure of response when adding a node is requested
type NodeAddResponse struct {
	// Node is the node information for the added node.
	Node *NodeInfo `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	// Peers is a list of all peers after adding the new node.
	Peers []*NodeInfo `protobuf:"bytes,2,rep,name=peers" json:"peers,omitempty"`
}

// NodeRemoveRequest is the structure to represent the request to remove a node
type NodeRemoveRequest struct{}

// NodeRemoveResponse is the structure of response when removing a node is requested
type NodeRemoveResponse struct {
	// Peers is a list of all peers after removing the node.
	Peers []*NodeInfo `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
}
