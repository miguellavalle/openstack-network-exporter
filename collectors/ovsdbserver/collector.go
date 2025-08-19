// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Miguel Lavalle

package ovsdbserver

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	"github.com/openstack-k8s-operators/openstack-network-exporter/appctl"
	"github.com/openstack-k8s-operators/openstack-network-exporter/collectors/lib"
	"github.com/openstack-k8s-operators/openstack-network-exporter/config"
	"github.com/openstack-k8s-operators/openstack-network-exporter/log"
	"github.com/prometheus/client_golang/prometheus"
)

type raftClusterInfo struct {
	database        string
	clusterUUID     string
	serverUUID      string
	role            string
	status          string
	vote            string
	term            int
	electionTimer   int
	logStart        int
	logNext         int
	logNotCommitted int
	logNotApplied   int
	inboundConns    int
	outboundConns   int
	isLeader        bool
}

var (
	// Regular expressions for parsing cluster status output
	clusterIdRe     = regexp.MustCompile(`^Cluster ID: (\w+) \(([^)]+)\)$`)
	serverIdRe      = regexp.MustCompile(`^Server ID: (\w+) \(([^)]+)\)$`)
	roleRe          = regexp.MustCompile(`^Role: (.+)$`)
	statusRe        = regexp.MustCompile(`^Status: (.+)$`)
	termRe          = regexp.MustCompile(`^Term: (\d+)$`)
	voteRe          = regexp.MustCompile(`^Vote: (.+)$`)
	electionTimerRe = regexp.MustCompile(`^Election timer: (\d+)$`)
	logRe           = regexp.MustCompile(`^Log: \[(\d+), (\d+)\]$`)
	notCommittedRe  = regexp.MustCompile(`^Entries not yet committed: (\d+)$`)
	notAppliedRe    = regexp.MustCompile(`^Entries not yet applied: (\d+)$`)
	connectionsRe   = regexp.MustCompile(`^Connections: (.+)$`)
	nameRe          = regexp.MustCompile(`^Name: (.+)$`)
)

func parseConnections(connStr string) (int, int) {
	// Parse connections like "<-200e ->200e <-cd93 ->cd93"
	inbound := 0
	outbound := 0

	parts := strings.Fields(connStr)
	for _, part := range parts {
		if strings.HasPrefix(part, "<-") {
			inbound++
		} else if strings.HasPrefix(part, "->") {
			outbound++
		}
	}

	return inbound, outbound
}

func parseClusterStatus(output string) (*raftClusterInfo, error) {
	info := &raftClusterInfo{}
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case nameRe.MatchString(line):
			match := nameRe.FindStringSubmatch(line)
			if len(match) == 2 {
				info.database = match[1]
			}

		case clusterIdRe.MatchString(line):
			match := clusterIdRe.FindStringSubmatch(line)
			if len(match) == 3 {
				info.clusterUUID = match[2] // Use full UUID
			}

		case serverIdRe.MatchString(line):
			match := serverIdRe.FindStringSubmatch(line)
			if len(match) == 3 {
				info.serverUUID = match[2] // Use full UUID
			}

		case roleRe.MatchString(line):
			match := roleRe.FindStringSubmatch(line)
			if len(match) == 2 {
				info.role = match[1]
				info.isLeader = (match[1] == "leader")
			}

		case statusRe.MatchString(line):
			match := statusRe.FindStringSubmatch(line)
			if len(match) == 2 {
				info.status = match[1]
			}

		case voteRe.MatchString(line):
			match := voteRe.FindStringSubmatch(line)
			if len(match) == 2 {
				info.vote = match[1]
			}

		case termRe.MatchString(line):
			match := termRe.FindStringSubmatch(line)
			if len(match) == 2 {
				if val, err := strconv.Atoi(match[1]); err == nil {
					info.term = val
				}
			}

		case electionTimerRe.MatchString(line):
			match := electionTimerRe.FindStringSubmatch(line)
			if len(match) == 2 {
				if val, err := strconv.Atoi(match[1]); err == nil {
					info.electionTimer = val
				}
			}

		case logRe.MatchString(line):
			match := logRe.FindStringSubmatch(line)
			if len(match) == 3 {
				if start, err := strconv.Atoi(match[1]); err == nil {
					info.logStart = start
				}
				if next, err := strconv.Atoi(match[2]); err == nil {
					info.logNext = next
				}
			}

		case notCommittedRe.MatchString(line):
			match := notCommittedRe.FindStringSubmatch(line)
			if len(match) == 2 {
				if val, err := strconv.Atoi(match[1]); err == nil {
					info.logNotCommitted = val
				}
			}

		case notAppliedRe.MatchString(line):
			match := notAppliedRe.FindStringSubmatch(line)
			if len(match) == 2 {
				if val, err := strconv.Atoi(match[1]); err == nil {
					info.logNotApplied = val
				}
			}

		case connectionsRe.MatchString(line):
			match := connectionsRe.FindStringSubmatch(line)
			if len(match) == 2 {
				info.inboundConns, info.outboundConns = parseConnections(match[1])
			}
		}
	}

	return info, scanner.Err()
}

func collectRaftMetrics(info *raftClusterInfo, ch chan<- prometheus.Metric) {
	baseLabels := []string{info.database, info.clusterUUID, info.serverUUID}

	// Cluster election timer
	if config.MetricSets().Has(clusterElectionTimer.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterElectionTimer.Desc(), clusterElectionTimer.ValueType,
			float64(info.electionTimer), baseLabels...)
	}

	// Cluster ID (constant 1.0)
	if config.MetricSets().Has(clusterId.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterId.Desc(), clusterId.ValueType,
			1.0, info.database, info.clusterUUID)
	}

	// Cluster Server ID (constant 1.0)
	if config.MetricSets().Has(clusterServerId.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterServerId.Desc(), clusterServerId.ValueType,
			1.0, baseLabels...)
	}

	// Cluster Server Role (constant 1.0)
	if config.MetricSets().Has(clusterServerRole.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterServerRole.Desc(), clusterServerRole.ValueType,
			1.0, append(baseLabels, info.role)...)
	}

	// Cluster Server Status (constant 1.0)
	if config.MetricSets().Has(clusterServerStatus.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterServerStatus.Desc(), clusterServerStatus.ValueType,
			1.0, append(baseLabels, info.status)...)
	}

	// Cluster Server Vote (constant 1.0)
	if config.MetricSets().Has(clusterServerVote.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterServerVote.Desc(), clusterServerVote.ValueType,
			1.0, append(baseLabels, info.vote)...)
	}

	// Cluster Term
	if config.MetricSets().Has(clusterTerm.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterTerm.Desc(), clusterTerm.ValueType,
			float64(info.term), baseLabels...)
	}

	// Cluster Leader (1.0 if leader, 0.0 if not)
	if config.MetricSets().Has(clusterLeader.Set) {
		var leaderValue float64
		if info.isLeader {
			leaderValue = 1.0
		} else {
			leaderValue = 0.0
		}
		ch <- prometheus.MustNewConstMetric(
			clusterLeader.Desc(), clusterLeader.ValueType,
			leaderValue, baseLabels...)
	}

	// Inbound connections
	if config.MetricSets().Has(clusterInboundConnectionsTotal.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterInboundConnectionsTotal.Desc(), clusterInboundConnectionsTotal.ValueType,
			float64(info.inboundConns), baseLabels...)
	}

	// Outbound connections
	if config.MetricSets().Has(clusterOutboundConnectionsTotal.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterOutboundConnectionsTotal.Desc(), clusterOutboundConnectionsTotal.ValueType,
			float64(info.outboundConns), baseLabels...)
	}

	// Log entry index
	if config.MetricSets().Has(logEntryIndex.Set) {
		ch <- prometheus.MustNewConstMetric(
			logEntryIndex.Desc(), logEntryIndex.ValueType,
			float64(info.logStart), baseLabels...)
	}

	// Log index next
	if config.MetricSets().Has(clusterLogIndexNext.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterLogIndexNext.Desc(), clusterLogIndexNext.ValueType,
			float64(info.logNext), baseLabels...)
	}

	// Log not committed
	if config.MetricSets().Has(clusterLogNotCommitted.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterLogNotCommitted.Desc(), clusterLogNotCommitted.ValueType,
			float64(info.logNotCommitted), baseLabels...)
	}

	// Log not applied
	if config.MetricSets().Has(clusterLogNotApplied.Set) {
		ch <- prometheus.MustNewConstMetric(
			clusterLogNotApplied.Desc(), clusterLogNotApplied.ValueType,
			float64(info.logNotApplied), baseLabels...)
	}
}

type Collector struct{}

func (Collector) Name() string {
	return "ovsdbserver"
}

func (Collector) Metrics() []lib.Metric {
	return metrics
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	lib.DescribeEnabledMetrics(c, ch)
}

func (Collector) Collect(ch chan<- prometheus.Metric) {
	output := appctl.OvsDbServer("cluster/status")
	if output == "" {
		log.Debugf("No OVN Raft cluster status output available")
		return
	}

	info, err := parseClusterStatus(output)
	if err != nil {
		log.Errf("Failed to parse OVN Raft cluster status: %s", err)
		return
	}

	collectRaftMetrics(info, ch)
}
