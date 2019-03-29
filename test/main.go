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
	tlen, _ := strconv.Atoi(os.Args[4])
	spec, ok := dLola.GetCheckedSpec(filename)
	if ok {
		mons := dLola.BuildMonitorTopo(spec, past_future, trigger, tlen)
		//dLola.Tickn(mons, 3)
		prefix := "[dLola_Monitor_Runtime]: "
		//fmt.Printf("%sMons: %s\n", prefix, dLola.PrintMons(mons))
		verdict := dLola.ConvergeCountTrigger(mons)
		fmt.Printf("%sVerdict: %s\n", prefix, verdict.Short())
	}
}

//topology & nmons are declared in the spec
//go run main.go inputMonitor.txt past trigger 2
