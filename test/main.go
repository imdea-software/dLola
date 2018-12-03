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
	spec, ok := dLola.GetCheckedSpec(filename)
	if ok {
		mons := dLola.BuildMonitorTopo(spec, past_future, trigger, topo, nmons, tlen)
		//dLola.Tickn(mons, 4)
		prefix := "[dLola_Monitor_Builder]: "
		verdict := dLola.ConvergeCountTrigger(mons)
		fmt.Printf("%sVerdict: %s\n", prefix, verdict.Short())
	}
}

/*func instantiateSpec(spec *dLola.Spec, tick, tlen int) {
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
*/
//go run main.go inputMonitor.txt past trigger clique 2 20
