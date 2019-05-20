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
	verbose := "no"
	if len(os.Args) > 5 {
		verbose = os.Args[5]
	}
	maxrounds := 0
	if len(os.Args) > 6 {
		maxrounds, _ = strconv.Atoi(os.Args[6])
	}

	if ok {
		mons := dLola.BuildMonitorTopo(spec, past_future, trigger, tlen)
		if maxrounds != 0 {
			dLola.Tickn(mons, maxrounds)
		} else {
			//prefix := "[dLola_Monitor_Runtime]: "
			//fmt.Printf("%sMons: %s\n", prefix, dLola.PrintMons(mons))
			verdict := dLola.ConvergeCountTrigger(mons)
			if verbose == "verbose" {
				fmt.Printf("%s\n" /*prefix,*/, verdict.String())
			} else {
				fmt.Printf("%s\n" /*prefix,*/, verdict.Short())
			}
		}
	}
}

//topology & nmons are declared in the spec
//go run main.go inputMonitor.txt past trigger 2
//go run main.go inputMonitor.txt past trigger 2 5
