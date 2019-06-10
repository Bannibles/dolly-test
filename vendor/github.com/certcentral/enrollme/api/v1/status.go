package v1

import (
	"time"
)

// ServerStatus provides response about current server
type ServerStatus struct {
	// NodeName is the human-readable name of the cluster member.
	NodeName string `protobuf:"bytes,1,name=nodename,proto3" json:"nodename"`
	// HostName is operating system's host name
	HostName string `protobuf:"bytes,2,name=hostname,proto3" json:"hostname"`
	// Port is the listening port
	Port string `protobuf:"bytes,3,name=port,proto3" json:"port"`
	// StartedAt is the time when the server started
	StartedAt time.Time `protobuf:"bytes,4,name=started_at,proto3" json:"started_at"`
	// Uptime is the total time elapsed since the server started
	Uptime time.Duration `protobuf:"bytes,5,name=uptime,proto3" json:"uptime"`
	// Version is the app build version
	Version string `protobuf:"bytes,6,name=version,proto3" json:"version"`
	// Peers is a list of all members associated with the cluster.
	Peers []*NodeInfo `protobuf:"bytes,7,rep,name=peers" json:"peers,omitempty"`
	// LeaderID is the node ID of the cluster leader
	LeaderID string `protobuf:"bytes,8,name=leader_id,proto3" json:"leader_id"`
	// NodeID is the node ID in the cluster.
	NodeID string `protobuf:"bytes,9,opt,name=node_id,proto3" json:"node_id,omitempty"`
}

// StatusResponse is response for /v1/status
type StatusResponse struct {
	Status *ServerStatus `json:"status"`
}
