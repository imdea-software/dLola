package dLola

import (
	//	"errors"
	"fmt"
)

//adjacency matrix of the Dependency Graph
//src dst weight of the adjacency
type DepGraph map[string](map[string]int)

type Path struct {
	weight int
	path   []string
}

/*func (p Path) Sprint() string {
	return fmt.Sprintf("%+v\n", p)
}*/

type ClasifiedPaths struct {
	negs  []Path
	zeros []Path
	pos   []Path
}

func NewClasifiedPaths() ClasifiedPaths {
	return ClasifiedPaths{make([]Path, 0), make([]Path, 0), make([]Path, 0)}
}

func (c ClasifiedPaths) Sprint() string {
	return fmt.Sprintf("%v", c)
}

func HasSelfRef(g DepGraph, src string) bool {
	pending := Expand(g, src)
	for cap(pending) != 0 {
		//fmt.Printf("pending:%s\n", pending)
		head := pending[0]
		pending = pending[1:] //drop head
		if head == src {
			return true
		} else {
			l := Expand(g, head)
			pending = append(pending, l...) //as append is a variadic function(take an arbitrary #args, with this notation it accepts a slice)
		}
	} // if we reach this point, there were no self-references
	return false
}

/*in contrast to Self_ref this function will return EVERY loop in g from src and back*/
func SelfRefLoops(g DepGraph, src string) []Path {
	return visitNode(g, src, src, Path{0, []string{src}}, map[string]struct{}{}) //visitedNodes is a Set of strings
}

/*Expands the adjacencies of cur, and then searches for loops on them updating the path*/
func visitNode(g DepGraph, src, cur string, path Path, visitedNodes map[string]struct{}) []Path {
	visitedNodes[cur] = struct{}{} //add it to the set so the node is not visited again
	pending := Expand(g, cur)
	//	fmt.Printf("pending:%s\n", pending)
	loops := make([]Path, 0)
	for _, c := range pending {
		cpath := Path{path.weight + g[cur][c], append(path.path, c)} //adjacency traversed, NOTE IT IS A NEW PATH
		loops = append(loops, selfRefLoopsAux(g, src, c, cpath, visitedNodes)...)
	}
	return loops
}

//cur will always be an adjacency of src, a child in the exploration path
func selfRefLoopsAux(g DepGraph, src string, cur string, path Path, visitedNodes map[string]struct{}) []Path {
	loops := make([]Path, 0)
	if cur == src {
		loops = append(loops, path) //we found a loop from src [to other nodes] to src, so we add the path of the loop
	} else {
		if _, ok := visitedNodes[cur]; !ok { //only if not already visited we visit, IMPORTANT: every adjacency is explored just once!!!
			loops = visitNode(g, src, cur, path, visitedNodes)
		}
	}
	return loops
}

func Expand(g DepGraph, src string) []string {
	adjacent := g[src]
	res := make([]string, 0)
	for key, _ := range adjacent {
		res = append(res, key)
	}
	return res
}

func ClasifyPaths(paths []Path) ClasifiedPaths {
	cpaths := NewClasifiedPaths()
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
}
