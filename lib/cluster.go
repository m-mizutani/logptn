package logptn

import (
// logging "log"
)

type Cluster struct {
	logs []*Log
	Base *Log
}

func CalcDistance(a, b *Log) float64 {
	if len(a.Chunk) != len(b.Chunk) {
		return 0
	}

	matched := 0
	for idx, ac := range a.Chunk {
		if ac.Data == b.Chunk[idx].Data {
			matched += 1
		}
	}

	return float64(matched) / float64(len(a.Chunk))
}

func AppendNewCluster(clusters []*Cluster, log *Log) []*Cluster {
	cluster := Cluster{}
	cluster.Base = log
	// logging.Println("new cluster: ", log)
	clusters = append(clusters, &cluster)
	return clusters
}

func Clustering(logs []*Log) []*Cluster {
	var clusters []*Cluster
	THREATHOLD := 0.8

	for _, log := range logs {
		if len(clusters) == 0 {
			clusters = AppendNewCluster(clusters, log)
		} else {
			var max float64
			var maxIdx int

			for idx, cluster := range clusters {
				d := CalcDistance(log, cluster.Base)
				// logging.Println("cmp -> ", d)
				if max < d {
					max, maxIdx = d, idx
				}
			}

			if max < THREATHOLD {
				clusters = AppendNewCluster(clusters, log)
			} else {
				// logging.Println("merged")
				clusters[maxIdx].logs = append(clusters[maxIdx].logs, log)
			}
		}
	}

	return clusters
}
