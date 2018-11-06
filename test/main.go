package main

import (
	//	"errors"
	"fmt"
	dLola "gitlab.software.imdea.org/luismiguel.danielsson/dLola"
	"os"
	"strconv"
)

func main() {
	filename := os.Args[1]
	past_future := os.Args[2]
	trigger := os.Args[3]
	topo := os.Args[4]
	nmons, _ := strconv.Atoi(os.Args[5])
	tlen, _ := strconv.Atoi(os.Args[6])
	getSpec(filename, past_future, trigger, topo, nmons, tlen)
}

func getSpec(filename, past_future, trigger, topo string, nmons, tlen int) {
	prefix := "[dLola_compiler]: "
	fmt.Printf(prefix+"Parsing file %s\n", filename)
	spec, err := dLola.GetSpec(filename, prefix)
	if err != nil {
		fmt.Printf("There was an error while parsing: %s\n", err)
		return
	}
	//dLola.PrintSpec(spec, prefix)
	//fmt.Printf(prefix + "Generating Pretty Print\n")
	//fmt.Printf(dLola.PrettyPrintSpec(spec, prefix))
	dLola.CheckTypesSpec(spec, prefix)
	if analyzeWF(spec) {
		buildMonitor(spec, past_future, trigger, topo, tlen, nmons)
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

func buildMonitor(spec *dLola.Spec, past_future, trigger, topo string, tlen, nmons int) {
	prefix := "[dLola_Monitor_Builder]: "
	fmt.Printf("%sBuilding Monitor...\n", prefix)
	delta := dLola.RoundrDelta(*spec, nmons)
	fmt.Printf("Delta:%v\n", delta)
	req := dLola.GenerateReqs(spec, past_future, trigger, tlen, delta)
	fmt.Printf("Generated Reqs:%v\n", req)
	mons := dLola.BuildMonitors(tlen, nmons, spec, req, delta, topo)
	verdict := dLola.ConvergeCountTrigger(mons)
	//dLola.Tickn(mons, 4)
	fmt.Printf("Verdict: %s\n", verdict.Short())

	/*	f := dLola.RootStream("one", dLola.SpecToGraph(spec))
		fmt.Printf("One is root:%t\n", f)
		f2 := dLola.RootStream("two", dLola.SpecToGraph(spec))
		fmt.Printf("Two is root:%t\n", f2)
	*/
}

func instantiateSpec(spec *dLola.Spec, tick, tlen int) {
	if tick >= 0 && tick < tlen { //othw may produce errors!! for those streams with no shift
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
