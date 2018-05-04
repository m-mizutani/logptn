package logptn

import (
// logging "log"
)

// Cluster is an interface of a set of Log
type Cluster interface {
	Logs() []*Log
	Length() int
}

// ClusterBuilder is an interface of clustering worker.
type ClusterBuilder interface {
	Clustering(logs []*Log) []Cluster
}

// -------------------------------
// Cluster builder implementation: SimpleClusterBuilder & SimpleCluster
type simpleCluster struct {
	logs []*Log
	base *Log
}

func (x *simpleCluster) Logs() []*Log {
	return x.logs
}

func (x *simpleCluster) Length() int {
	return len(x.base.Chunk)
}

// simpleClusterBuilder is clustering builder.
// This builder's algorithm is very simple. Clustering policy is below
// 1. Clustering only logs having same length of chunk
// 2. Clustering if chunk matching ratio is over 0.7 (default)
type simpleClusterBuilder struct {
	threshold float64
	clusters  []*simpleCluster
}

// Constructor of simpleClusterBuilder
func NewSimpleClusterBuilder() *simpleClusterBuilder {
	builder := simpleClusterBuilder{threshold: 0.7}
	return &builder
}

func (x *simpleClusterBuilder) SetThreshold(threshold float64) {
	x.threshold = threshold
}

func (x *simpleClusterBuilder) Clustering(logs []*Log) []Cluster {
	calcDistance := func(a, b *Log) float64 {
		if len(a.Chunk) != len(b.Chunk) {
			return 0
		}

		matched := 0
		for idx, ac := range a.Chunk {
			if ac.Data == b.Chunk[idx].Data {
				matched++
			}
		}

		return float64(matched) / float64(len(a.Chunk))
	}

	THREATHOLD := 0.65

	for _, log := range logs {
		if len(x.clusters) == 0 {
			x.appendCluster(log)
		} else {
			var max float64
			var maxIdx int

			for idx, cluster := range x.clusters {
				d := calcDistance(log, cluster.base)
				// logging.Println("cmp -> ", d)
				if max < d {
					max, maxIdx = d, idx
				}
			}

			if max < THREATHOLD {
				x.appendCluster(log)
			} else {
				// logging.Println("merged")
				x.clusters[maxIdx].logs = append(x.clusters[maxIdx].logs, log)
			}
		}
	}

	res := make([]Cluster, len(x.clusters))
	for idx, cluster := range x.clusters {
		res[idx] = cluster
	}
	return res
}

func (x *simpleClusterBuilder) appendCluster(log *Log) {
	cluster := simpleCluster{}
	cluster.base = log
	cluster.logs = append(cluster.logs, log)

	x.clusters = append(x.clusters, &cluster)
}
