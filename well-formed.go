package dLola

import (
	"errors"
	"fmt"
)

/*Graph represented as a map*/
type DepGraphAdj map[string][]Adj
type Reachable map[string](map[string]struct{}) //will contain if there is a path (of any length) from the first to the second

type Adj struct {
	Src    string
	Weight int
	Dest   string
}

func (a Adj) Sprint() string {
	//fmt.Printf("printing Adj %d\n", a.Weight)
	return fmt.Sprintf("Adjacency{Src = %s, Weight= %d, Dest = %s}", a.Src, a.Weight, a.Dest)
}

func EqAdj(a, a2 Adj) bool {
	return a.Weight == a2.Weight && a.Dest == a2.Dest && a.Src == a2.Src
}

func EqAdjs(a, a2 []Adj) bool {
	res := false
	if len(a) == len(a2) {
		i := 0
		for res = true; res; i++ {
			res = res && EqAdj(a[i], a2[i])
		}
	}
	return res
}

func SprintAdjs(as []Adj) string {
	//fmt.Printf("printing Adjs \n")
	s := "Adjs =["
	for _, a := range as {
		s += a.Sprint() + ","
	}
	return s + "]"
}

type PathAdj struct {
	weight int
	path   []Adj
}

func SprintPathAdj(p PathAdj) string {
	return fmt.Sprintf("PathAdj{ weight = %d, path = %s}", p.weight, SprintAdjs(p.path))
}

/*func (p Path) Sprint() string {
	return fmt.Sprintf("%+v\n", p)
}*/

type ClasifiedPathsAdj struct {
	negs  []PathAdj
	zeros []PathAdj
	pos   []PathAdj
}

func NewClasifiedPathsAdj() ClasifiedPathsAdj {
	return ClasifiedPathsAdj{make([]PathAdj, 0), make([]PathAdj, 0), make([]PathAdj, 0)}
}

func (c ClasifiedPathsAdj) Sprint() string {
	return fmt.Sprintf("%v", c)
}

/*functions*/
func SpecToGraph(spec *Spec) DepGraphAdj {
	sToGVisitor := SpecToGraphVisitor{DepGraphAdj{}, ""}
	for _, v := range spec.Output {
		sToGVisitor.s = v.Name.Sprint()
		v.Expr.Accept(&sToGVisitor)
	}
	return sToGVisitor.g
}

func checkDepGraphAdj(g DepGraphAdj) []error {
	err := make([]error, 0)
	for node, adjs := range g {
		for _, a := range adjs {
			if node != a.Src {
				err = append(err, errors.New(fmt.Sprintf("Adj %s does not start at node %s", a.Sprint(), node)))
			}
		}
	}
	return err
}

func GetReachableAdj(g DepGraphAdj) (Reachable, []error) {
	err := checkDepGraphAdj(g)
	if len(err) != 0 {
		return Reachable{}, err
	}
	reach := make(Reachable, 0)
	for node, pending := range g {
		//fmt.Printf("Node:%s\n", node)
		reach[node] = make(map[string]struct{})
		i := 0
		for len(pending) != 0 /*&& i < 5*/ { // cap(pending)
			//fmt.Printf("New it %d pending:%v\n", i, pending)
			head := pending[0]
			pending = pending[1:] //drop head
			if _, ok := reach[node][head.Dest]; !ok /*|| head.Dest == node*/ {
				reach[node][head.Dest] = struct{}{} //mark that head.Dest is reachable from node
				//fmt.Printf("reach:%v\n", reach)
				l := g[head.Dest]          // take adjacencies of the destiny
				l = filter(l, reach[node]) //filter those adjacencies that lead to an already reachable node
				//fmt.Printf("appending :%v\n", l)
				pending = append(pending, l...) //as append is a variadic function(take an arbitrary #args, with this notation it accepts a slice)
			}
			i++
		}
	}
	return reach, []error{}
}

/*will filter all those adjacencies of a node that lead to an already reachable node*/
func filter(l []Adj, reachable map[string]struct{}) []Adj {
	//fmt.Printf("Candidates for appending :%v\n", l)
	for _, a := range l {
		if _, ok := reachable[a.Dest]; ok { //if present we drop the element
			l = remove(l, a)
		}
	}
	return l
}

func visitedDest(path []Adj, node string) bool {
	//fmt.Printf("path %v, node %s", path, node)
	is := false
	for i := 0; i < len(path) && !is; i++ {
		if /*path[i].Src == node ||*/ path[i].Dest == node {
			is = true
		}
	}
	return is
}

func remove(l []Adj, a Adj) []Adj {
	r := make([]Adj, 0)
	for _, v := range l {
		if !EqAdj(v, a) {
			r = append(r, v)
		}
	}
	return r
}

func elem(as []Adj, a Adj, f func(Adj, Adj) bool) bool { //will usually be called with f = EqAdj
	//fmt.Printf("path %v, node %s", path, node)
	is := false
	for i := 0; i < len(as) && !is; i++ {
		if f(a, as[i]) {
			is = true
		}
	}
	return is
}

/*in contrast to Self_ref this function will return EVERY SIMPLE (cannot repeat nodes) loop in g from src and back*/
func SimpleCyclesAdj(g DepGraphAdj, r Reachable) ([]PathAdj, []error) {
	err := checkDepGraphAdj(g)
	if len(err) != 0 {
		return []PathAdj{}, err
	}
	loops := make([]PathAdj, 0)
	for src, _ := range g {
		if _, ok := r[src][src]; ok { //only if the node is reachable from itself we look for simple cycles
			loops = append(loops, visitNodeAdj(g, src, src, PathAdj{0, []Adj{}})...) //visitedNodes is a Set of strings)
		}
	}
	return loops, []error{}
}

//IMPORTANT MAPS & SLICES in GO are pointers, so as they are modified down in the recursion when going up, they are changed
/*Expands the adjacencies of cur, and then searches for loops on them updating the path*/
func visitNodeAdj(g DepGraphAdj, src, cur_node string, path PathAdj) []PathAdj {
	//fmt.Printf("I'm in node :%s\n", cur_node)
	pending := g[cur_node] // :: Adj
	//fmt.Printf("pending:%s\n", SprintAdjs(pending))
	loops := make([]PathAdj, 0)
	for _, c := range pending {
		//fmt.Printf("current node: %s path: %v\n", cur_node, path)
		loops = append(loops, decideAdj(g, src, c, path)...)
	}
	//fmt.Printf("returning loops of node %s:%v\n", cur_node, loops)
	return loops
}

//cur will always be an adjacency of src to a child in the exploration path
//decides if the adj takes to src, then we found a loop, othw it takes it and will continue exploring the node cur.Dest
func decideAdj(g DepGraphAdj, src string, cur Adj, path PathAdj) []PathAdj {
	//fmt.Printf("Deciding adj :%s\n", cur.Sprint())
	loops := make([]PathAdj, 0)
	if shouldVisit(cur, path) { //only if not already visited we visit, IMPORTANT: every node may appear in the path just once!!!
		//fmt.Printf("Decide to visit adj :%s\n", cur.Sprint())
		cpath := PathAdj{path.weight + cur.Weight, append(path.path, cur)} //adjacency traversed, NOTE IT IS A NEW PATH
		//visited[cur.Dest] = struct{}{} //add it to the set so the Adj is not visited again
		if cur.Dest == src {
			//fmt.Printf("Found LOOP!! :%v\n", SprintPathAdj(cpath))
			loops = []PathAdj{cpath} //we found a loop from src [to other nodes] to src, so we add the path of the loop
		} else {
			loops = visitNodeAdj(g, src, cur.Dest, cpath) //IMPORTANT cpath is passed as value, othw the backtracking of the recursion will produce errors in the path of the cycles!!!
		}

	} else {
		//fmt.Printf("Decide not to visit:%s\n", cur.Sprint())
	}

	return loops
}

func shouldVisit(cur Adj, path PathAdj) bool {
	return !visitedDest(path.path, cur.Dest)
}

func CreateCycleMap(cycles []PathAdj) map[string]ClasifiedPathsAdj {
	res := make(map[string]ClasifiedPathsAdj)
	for _, p := range cycles {
		//fmt.Printf("Path %v, res = %v\n", p, res)
		src := p.path[0].Src
		c, ok := res[src]
		if !ok { //there were not previous clasifiedPathsAdj
			//fmt.Printf("there were no previous cycles for this stream: %s\n", src)
			cpaths := NewClasifiedPathsAdj()
			res[src] = appendPath(&cpaths, p)
		} else {
			//fmt.Printf("there were previous cycles for this stream: %s\n", src)
			res[src] = appendPath(&c, p)
		}
	}
	//fmt.Printf("Returning Path %v\n", res)
	return res
}

/*Appends p to cpaths depending on its weight*/
func appendPath(cpaths *ClasifiedPathsAdj, p PathAdj) ClasifiedPathsAdj {
	//fmt.Printf("Appending Path %v to %v\n", p, cpaths)
	if p.weight < 0 {
		cpaths.negs = append(cpaths.negs, p)
	} else {
		if p.weight == 0 {
			cpaths.zeros = append(cpaths.zeros, p)
		} else {
			cpaths.pos = append(cpaths.pos, p)
		}
	}
	//fmt.Printf("Updated Path %v\n", cpaths)
	return *cpaths
}

/*func ClasifyPathsAdj(paths []PathAdj) ClasifiedPathsAdj {
	cpaths := NewClasifiedPathsAdj()
	for _, p := range paths {
		if p.weight < 0 {
			cpaths.negs = append(cpaths.negs, p)
		} else {
			if p.weight == 0 {
				cpaths.zeros = append(cpaths.zeros, p)
			} else {
				cpaths.pos = append(cpaths.pos, p)
			}
		}
	}
	return cpaths
}*/

func IsWF(m map[string]ClasifiedPathsAdj) []error {
	err := make([]error, 0)
	for stream, cycles := range m {
		if len(cycles.zeros) != 0 {
			err = append(err, errors.New(fmt.Sprintf("Stream %s has 0 weight cycles: %v", stream, cycles.zeros)))
		}
		if len(cycles.pos) != 0 && len(cycles.negs) != 0 {
			err = append(err, errors.New(fmt.Sprintf("Stream %s has both positive and negative weight simple cycles:Positives: %v Negatives: %v", stream, cycles.pos, cycles.negs)))
		}
		err = append(err, checkComplexCycles(stream, cycles, m)...)
	}
	return err
}

func checkComplexCycles(stream string, cycles ClasifiedPathsAdj, m map[string]ClasifiedPathsAdj) []error {
	err := make([]error, 0)
	for _, n := range cycles.negs {
		err = append(err, complexCycle(stream, n, m)...)
	}
	for _, z := range cycles.zeros {
		err = append(err, complexCycle(stream, z, m)...)
	}
	for _, p := range cycles.pos {
		err = append(err, complexCycle(stream, p, m)...)
	}
	return err
}

func complexCycle(stream string, c PathAdj, m map[string]ClasifiedPathsAdj) []error {
	err := make([]error, 0)
	for _, adj := range c.path {
		node := adj.Dest
		if node != stream {
			if c.weight < 0 && len(m[node].pos) != 0 {
				err = append(err, errors.New(fmt.Sprintf("Stream %s has a negative simple cycle %v which includes a node %s that has positive simple cycles: %v", stream, c, node, m[node].pos)))
			} else {
				if c.weight > 0 && len(m[node].negs) != 0 {
					err = append(err, errors.New(fmt.Sprintf("Stream %s has a positive simple cycle %v which includes a node %s that has negative simple cycles: %v", stream, c, node, m[node].negs)))
				}

			}
		}
	}
	return err
}
