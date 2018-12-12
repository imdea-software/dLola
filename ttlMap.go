package dLola

import (
	//	"errors"
	//"fmt"
	"math"
	//	"strconv"
)

type TTLMap = map[StreamName]Time
type RevDepGraph = DepGraphAdj

func getTTLMap(depen DepGraphAdj, delta map[StreamName]Id, dists map[Id]map[Id]int) TTLMap {
	ttlMap := make(TTLMap)
	revDepen := reverseDepGraph(depen)
	//fmt.Printf("revDepen %v\n", revDepen)
	for streamName, _ := range revDepen {
		getStreamTtl(streamName, revDepen, delta, ttlMap, dists)
	}
	//fmt.Printf("getTTL done %v\n", ttlMap)
	return ttlMap
}

/*consider that for local negative shifts the value need to be kept in R for the time specified in the shift
also consider that the distance between the requesting monitor (Lazy) and the resolver monitor needs to be added
for positive shifts the positive value decrements the value of the distance
the time that an element should be kept in R is always non negative
This function is suitable for:
Eager: Resolver(revDep is local so won't use distance) Receiver(revDep is local so won't use distance)
Lazy: Resolver (revDep is remote so will use distance or local and won't use it) & Receiver(revDep is local so won't use distance)
Eager Resolver(revDep is remote so will send msg and forget) will use 1 or 0
*/
func getStreamTtl(streamName StreamName, reverseDepen RevDepGraph, delta map[StreamName]Id, ttlMap TTLMap, dists map[Id]map[Id]int) {
	//fmt.Printf("getStreamTTL: %s\n", streamName)
	max := 0.0 //its the minimum, including inputs
	if _, in := ttlMap[streamName]; !in {
		if revDepends, ok := reverseDepen[streamName]; ok {
			for _, revDep := range revDepends { //revDep.Src == streamName
				//fmt.Printf("revDep %s\n", revDep.Sprint())
				ttl := -revDep.Weight //consider both positives & negatives shifts (positives need to be considered because of distances between computing monitors)
				if revDep.Dest != streamName {
					ttl += dists[delta[streamName]][delta[revDep.Dest]]
				}
				max = math.Max(max, float64(ttl))
			}
		}
		ttlMap[streamName] = int(max)
	}
}

/*example output a := b[-1|0] + 2; output c := b[1] ; define b ...
depGraph {a : {a -1 b}, c: {c 1 b}}
reverseDepGraph {b : [a, c]}
returns all the streams that need the value used as key for the map
*/
func reverseDepGraph(depGraph DepGraphAdj) RevDepGraph {
	superDependencies := make(RevDepGraph)
	for _, streamDeps := range depGraph {
		for _, streamDep := range streamDeps {
			superDependencies[streamDep.Dest] = append(superDependencies[streamDep.Dest], reverseAdj(streamDep)) //superDependencies[b] = append(superDependencies[b], a)
		}
	}
	return superDependencies
}

//DOES NOT CHANGE WEIGHT, just swaps src and dest
func reverseAdj(a Adj) Adj {
	return Adj{a.Dest, a.Weight, a.Src}
}
