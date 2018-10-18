package main

import (
	//	"errors"
	"fmt"
	dLola "gitlab.software.imdea.org/luismiguel.danielsson/dLola"
	"os"
)

func main() {
	prefix := "[dLola_compiler]: "
	s := os.Args[1]
	fmt.Printf(prefix+"Parsing file %s\n", s)
	spec, err := dLola.GetSpec(s, prefix)
	if err != nil {
		fmt.Printf("There was an error while parsing: %s\n", err)
		return
	}
	//dLola.PrintSpec(spec, prefix)
	fmt.Printf(prefix + "Generating Pretty Print\n")
	fmt.Printf(dLola.PrettyPrintSpec(spec, prefix))
	dLola.CheckTypesSpec(spec, prefix)

	if analyzeWF(spec) {
		buildMonitor(spec)
	}
}

func analyzeWF(spec *dLola.Spec) bool {
	prefix := "[dLola_well-formedness_checker]: "
	g := dLola.SpecToGraph(spec)
	fmt.Printf("%s Dependency Graph: %v\n", prefix, g)
	fmt.Printf(prefix+"%v\n", g)

	r, err := dLola.GetReachableAdj(g)
	for _, e := range err {
		fmt.Printf("%sERROR: %s\n", prefix, e)
	}
	if len(err) != 0 {
		return false
	}
	fmt.Printf(prefix+"Reachability table: %v\n", r)
	simples, err := dLola.SimpleCyclesAdj(g, r)
	for _, e := range err {
		fmt.Printf("%sERROR: %s\n", prefix, e)
	}
	if len(err) != 0 {
		return false
	}
	//fmt.Printf(prefix+"Simple cycles from %v\n", simples)
	cpaths := dLola.CreateCycleMap(simples)
	fmt.Printf(prefix+"Clasified paths: %v\n", cpaths)

	wf_err := dLola.IsWF(cpaths)
	for _, e := range wf_err {
		fmt.Printf("%sWell-Formed ERROR: %s\n", prefix, e)
	}
	if len(wf_err) != 0 {
		return false
	}
	return true
}

func buildMonitor(spec *dLola.Spec) {
	prefix := "[dLola_Monitor_Builder]: "
	fmt.Printf("%sBuilding Monitor...\n", prefix)
	/*instantiateSpec(spec, 2, 2)
	instantiateSpec(spec, 0, 2)
	instantiateSpec(spec, 1, 2)*/
	depGraph := dLola.SpecToGraph(spec)
	routes := map[dLola.Id]dLola.Id{0: 0}
	delta := dLola.RoundrDelta(*spec, 1)
	dependencies := dLola.InterestedMonitors(delta, depGraph)
	eval := []dLola.StreamName{}
	mon := dLola.NewMonitor(0, 2, *spec, routes, delta, eval, depGraph, dependencies)
	mons := map[dLola.Id]*dLola.Monitor{0: &mon}
	dLola.Converge(mons)
	fmt.Printf("Converged mons:%s\n", dLola.PrintMons(mons))
}

func instantiateSpec(spec *dLola.Spec, tick, tlen int) {
	if tick >= 0 && tick < tlen { //othw may produce errors!!
		prefix := "[dLola_Monitor_Builder]: "
		fmt.Printf("%sInstantiating spec for tick %d with tlen %d\n", prefix, tick, tlen)
		for _, o := range spec.Output {
			iexpr := o.Expr.InstantiateExpr(tick, tlen)
			fmt.Printf("%sInstantiated Expression: %s = %s\n", prefix, o.Name, iexpr.Sprint())
			s := dLola.InstStreamFetchExpr{"hard", 0}
			v := dLola.InstIntLiteralExpr{2}
			sexpr := iexpr.Substitute(s, v)
			fmt.Printf("%sSubstituted Expression with pair: %s = %s\n %s = %s\n", prefix, s.Sprint(), v.Sprint(), o.Name, sexpr.Sprint())
			simpExpr, _ := sexpr.Simplify()
			fmt.Printf("%sSimplified Expression %s = %s\n", prefix, o.Name, simpExpr.Sprint())
		}
	}
}
