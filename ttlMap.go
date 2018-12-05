package dLola

import (
	//	"errors"
	//"fmt"
	"math"
	//	"strconv"
)

type TTLMap = map[StreamName]Time

func getTTLMap(depen DepGraphAdj, nid Id, delta map[StreamName]Id) TTLMap {
	ttlMap := make(TTLMap)
	for streamName, adjs := range depen {
		getStreamTtl(streamName, adjs, depen, delta, ttlMap)
	}
	return ttlMap
}

func getStreamTtl(streamName StreamName, adjs []Adj, depen DepGraphAdj, delta map[StreamName]Id, ttlMap TTLMap) Time {
	max := 1.0 //its the minimum, including inputs
	if _, in := ttlMap[streamName]; !in {
		for _, a := range adjs {
			ttl := a.Weight
			if a.Dest != streamName {
				ttl += getStreamTtl(a.Dest, depen[a.Dest], depen, delta, ttlMap)
			}
			max = math.Max(max, float64(ttl))
		}
		ttlMap[streamName] = int(max)
	}
	return int(max)
}
