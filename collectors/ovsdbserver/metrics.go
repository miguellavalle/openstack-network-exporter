package ovsdbserver

import (
	"github.com/openstack-k8s-operators/openstack-network-exporter/collectors/lib"
	"github.com/openstack-k8s-operators/openstack-network-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

var clusterElectionTimer = lib.Metric{
	Name:        "ovn_raft_cluster_election_timer",
	Description: "A metric with the value of the election timer labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterId = lib.Metric{
	Name:        "ovn_raft_cluster_id",
	Description: "A metric with a constant '1' value labeled by database name and cluster uuid",
	Labels:      []string{"database", "cluster_uuid"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterServerId = lib.Metric{
	Name:        "ovn_raft_cluster_server_id",
	Description: "A metric with a constant '1' value labeled by database name, cluster uuid and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterServerRole = lib.Metric{
	Name:        "ovn_raft_cluster_server_role",
	Description: "A metric with a constant '1' value labeled by database name, cluster uuid, server uuid and role",
	Labels:      []string{"database", "cluster_uuid", "server_uuid", "role"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterServerStatus = lib.Metric{
	Name:        "ovn_raft_cluster_server_status",
	Description: "A metric with a constant '1' value labeled by database name, cluster uuid, server uuid and status",
	Labels:      []string{"database", "cluster_uuid", "server_uuid", "status"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterServerVote = lib.Metric{
	Name:        "ovn_raft_cluster_server_vote",
	Description: "A metric with a constant '1' value labeled by database name, cluster uuid, server uuid and vote",
	Labels:      []string{"database", "cluster_uuid", "server_uuid", "vote"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterTerm = lib.Metric{
	Name:        "ovn_raft_cluster_term",
	Description: "A metric with the value of the cluster term labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterLeader = lib.Metric{
	Name:        "ovn_raft_cluster_leader",
	Description: "A metric with value 1.0 if the server is the cluster leader for the given database or 0.0 if it is not, labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.GaugeValue,
	Set:         config.METRICS_BASE,
}

var clusterInboundConnectionsTotal = lib.Metric{
	Name:        "ovn_raft_cluster_inbound_connections_total",
	Description: "A metric with the value of total number of inbound connections to the server labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var clusterOutboundConnectionsTotal = lib.Metric{
	Name:        "ovn_raft_cluster_outbound_connections_total",
	Description: "A metric with the value of the total number of outbound connections from the server labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var logEntryIndex = lib.Metric{
	Name:        "ovn_raft_log_entry_index",
	Description: "A metric with the value of log entry index currently exposed to clients, labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var clusterLogIndexNext = lib.Metric{
	Name:        "ovn_raft_cluster_log_index_next",
	Description: "A metric with the value of the next log entry index labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var clusterLogNotCommitted = lib.Metric{
	Name:        "ovn_raft_cluster_log_not_committed",
	Description: "A metric with the value of the number of log entries not committed labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var clusterLogNotApplied = lib.Metric{
	Name:        "ovn_raft_cluster_log_not_applied",
	Description: "A metric with the value of the number of log entries not applied labeled by database name, cluster uuid, and server uuid",
	Labels:      []string{"database", "cluster_uuid", "server_uuid"},
	ValueType:   prometheus.CounterValue,
	Set:         config.METRICS_COUNTERS,
}

var metrics = []lib.Metric{
	clusterElectionTimer,
	clusterId,
	clusterServerId,
	clusterServerRole,
	clusterServerStatus,
	clusterServerVote,
	clusterTerm,
	clusterLeader,
	clusterInboundConnectionsTotal,
	clusterOutboundConnectionsTotal,
	logEntryIndex,
	clusterLogIndexNext,
	clusterLogNotCommitted,
	clusterLogNotApplied,
}
