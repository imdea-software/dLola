package dLola

import (
	//	"fmt"
	"math"
)

func RoundrDelta(s Spec, nmons int) map[StreamName]Id {
	delta := make(map[StreamName]Id)
	var i Id = 0
	for _, o := range s.Output {
		delta[o.Name] = i
		i = (i + 1) % nmons
	}
	return delta
}

func GenerateRoutes(nmons, src int, topo string) map[Id]Id {
	switch topo {
	case "clique":
		return cliqueRoute(nmons, src)
	case "ring":
		return ringRoute(nmons, src)
	case "ringshort":
		return ringshortRoute(nmons, src)
	case "line":
		return lineRoute(nmons, src)
	case "star":
		return starRoute(nmons, src, getArms(nmons))
	default:
		return cliqueRoute(nmons, src)
	}
}

func cliqueRoute(nmons, src int) map[Id]Id {
	r := make(map[Id]Id)
	for i := 0; i < nmons; i++ {
		r[i] = i
	}
	return r
}
func ringRoute(nmons, src int) map[Id]Id {
	r := make(map[Id]Id)
	for i := 0; i < nmons; i++ {
		if i == src {
			r[i] = src
		} else {
			r[i] = (src + 1) % nmons
		}

	}
	return r
}
func ringshortRoute(nmons, src int) map[Id]Id {
	r := make(map[Id]Id)
	for i := 0; i < nmons; i++ {
		if i == src {
			r[i] = i
		} else {
			if (i-src > 0 && math.Abs(float64(i-src))/float64(nmons) <= 0.5) || (i-src < 0 && math.Abs(float64(i-src))/float64(nmons) > 0.5) {
				r[i] = (src + 1) % nmons
			} else {
				r[i] = (src - 1) % nmons
			}
		}
	}
	return r
}
func lineRoute(nmons, src int) map[Id]Id {
	r := make(map[Id]Id)
	for i := 0; i < nmons; i++ {
		if i == src {
			r[i] = src
		} else {
			if i > src {
				r[i] = (src + 1) % nmons
			} else {
				r[i] = (src - 1) % nmons
			}
		}

	}
	return r
}

/*PRE: nodes-1/arms :: Nat and 0 <= src <= nodes-1
this should be an int so call getStarPairs(5,_,4) getStarPairs(9,_,4) getStarPairs(7,_,3) getStarPairs(10,_,3)
POST: pairs of Star, the center will be 0 and the numbering will be clockwise starting at 12:00 from inside to outside in spiral
this way all nodes in a given arm will have the following characteristic property:
id mod arms = armId
example:
getStarPairs(9,_,4):
    5
    1
8 4 0 2 6
    3
    7
has 4 arms, 9 nodes, the center is 0 and every nodes meets the property:e.g. 1 mod 4 = 1 which is the id of the north arm
Note: in order to move inside the arm the difference between nodes is the number of arms: from 1 to 5 we add (because we go outwards) 4 which is the number of arms
getStarPairs(9,1,4) = '(0,0),(1,1),(2,0),(3,0),(4,0),(5,5),(6,0),(7,0),(8,0)'
getStarPairs(9,0,4) = '(0,0),(1,1),(2,2),(3,3),(4,4),(5,1),(6,2),(7,3),(8,4)'*/
func starRoute(nmons, src, arms int) map[Id]Id {
	//nodesPerArm = (n-1)/arms
	//0 will be the central node
	//srcArm = src % arms
	r := make(map[Id]Id)
	if src == 0 { //src is the center of Star
		r = getStarPairsCenter(nmons, arms)
	} else { //src in an arm
		for i := 0; i < nmons; i++ { //from 0 to n-1
			//print i
			if i == src {
				r[i] = i
			} else {
				if sameArm(src, i, arms) && i > src {
					r[i] = src + arms //outwards in the same arm
				} else {
					if !firstOrbit(src, arms) && (!sameArm(src, i, arms) || i < src) { //if i in other arm or in the same arm but closer to the center
						r[i] = src - arms //hop closer to the center
					} else { //case of src in first orbit and dst not in its same arm
						r[i] = 0 //to the center
					}
				}
			}
		}
	}
	return r
}
func getStarPairsCenter(nmons, arms int) map[Id]Id {
	r := make(map[Id]Id)
	for i := 0; i < nmons; i++ {
		if i == 0 {
			r[i] = i
		} else {
			if i%arms != 0 {
				r[i] = i % arms //outwards
			} else {
				r[i] = i%arms + arms //outwards by the 0 arm, in order to avoid generating pair like (4,0) which is not correct
			}
		}
	}
	return r
}
func sameArm(src, dst, arms int) bool {
	if (src == 0 || dst == 0) && src != dst {
		return false
	}
	return dst%arms == src%arms
}

func firstOrbit(src, arms int) bool {
	return src < arms
}

func getArms(nmons int) int {
	arms := 0
	switch nmons {
	case 5:
		arms = 4
	case 9:
		arms = 4
	case 4:
		arms = 3
	case 7:
		arms = 3
	case 10:
		arms = 3
	default:

	}
	return arms
}

func InterestedMonitors(delta map[StreamName]Id, depGraph DepGraphAdj) map[StreamName][]Id {
	monitorDependencies := make(map[StreamName][]Id)
	for stream, mon := range delta {
		streamDeps := depGraph[stream]
		for _, streamDep := range streamDeps {
			monitorDependencies[streamDep.Dest] = append(monitorDependencies[streamDep.Dest], mon) //the Monitor assigned to compute the stream stream will need the value of streamDep.Dst TODO revise append to nil
		}
	}
	return monitorDependencies
}

func GenerateReqs(spec *Spec, past_future, trigger string, tlen int, delta map[StreamName]Id) map[Id][]Msg {
	var tick_req int
	switch past_future {
	case "past":
		tick_req = tlen - 1
	case "future":
		tick_req = 0
	default:
		tick_req = tlen - 1
	}
	var kind MsgType
	switch trigger {
	case "trigger":
		kind = Trigger
	default:
		kind = Req
	}
	reqs := make(map[Id][]Msg)
	//fmt.Printf("Generrating reqs\n")
	depGraph := SpecToGraph(spec)
	for _, o := range spec.Output {
		if RootStream(o.Name, depGraph) {
			stream := InstStreamFetchExpr{o.Name, tick_req}
			dst := delta[o.Name]
			m := Msg{kind, stream, nil, nil, nil, dst, dst} //src of the msgs will be themselves so they do not emit a response msg
			reqsi, ok := reqs[dst]
			if ok {
				//fmt.Printf("There were prev reqs\n")
				reqs[dst] = append(reqsi, m)
			} else {
				//fmt.Printf("There were NO prev reqs\n")
				reqs[dst] = []Msg{m}
			}
		}
	}
	return reqs
}

/*return if s does NOT appear in the right expression of any of the OTHER streams of the spec, it can appear in the right hand side of its own definition!!!*/
func RootStream(s StreamName, depGraph DepGraphAdj) bool {
	root := true
	for _, dep := range depGraph {
		for _, d := range dep {
			if s == d.Dest && s != d.Src { //is not root if it appears in the right hand side of another stream
				root = false
				break
			}
		}
	}
	return root
}

func BuildMonitors(tlen, nmons int, spec *Spec, reqs map[Id][]Msg, delta map[StreamName]Id, topo string) map[Id]*Monitor {
	mons := make(map[Id]*Monitor)
	depGraph := SpecToGraph(spec)
	dependencies := InterestedMonitors(delta, depGraph)
	//eval := map[StreamName]struct{}{}
	for i := 0; i < nmons; i++ {
		routes := GenerateRoutes(nmons, i, topo)
		channels := GenerateChannels(delta, spec, depGraph, i, tlen)
		mon := NewMonitor(i, tlen, *spec, reqs[i], routes, delta, depGraph, dependencies, channels)
		mons[i] = &mon
	}
	return mons
}

func GenerateChannels(delta map[StreamName]Id, spec *Spec, depGraph DepGraphAdj, id Id, tlen int) []chan Resolved {
	channels := make([]chan Resolved, 0)
	for stream, dependencies := range depGraph {
		if delta[stream] == id {
			for _, d := range dependencies {
				//fmt.Printf("%v\n", spec.Input)
				if inputDecl, ok := spec.Input[d.Dest]; ok {
					//fmt.Printf("found input %s for monitor %d\n", d.Dest, id)
					c := make(chan Resolved)
					generateInput(d.Dest, inputDecl.Type, c, tlen) //call to inputReader!!! TODO
					channels = append(channels, c)
				}
			}
		}
	}
	return channels
}
