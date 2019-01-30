package dLola

import (
	//dLola "gitlab.software.imdea.org/luismiguel.danielsson/dLola"
	// "fmt"
	"testing"
)

var verdict Verdict //so optimizations do not get rid of the call to ConvergeCountTrigger

func benchmarkMain(tlen int, b *testing.B) {
	spec, ok := GetCheckedSpec("test/generated/lotAcc.txt")
	if ok {
		mons := BuildMonitorTopo(spec, "past", "trigger", tlen)
		//fmt.Printf("Mons: %s\n", PrintMons(mons))
		b.ResetTimer() //to skip setups
		for i := 0; i < b.N; i++ {
			//fmt.Printf("testing again %d\n", tlen)
			verdict = ConvergeCountTrigger(mons)
		}
	}
}

func benchmarkMainN(tlen, ticks int, b *testing.B) {
	spec, ok := GetCheckedSpec("test/generated/lotAcc.txt")
	if ok {
		mons := BuildMonitorTopo(spec, "past", "trigger", tlen) //this must be inside the loop, the test environment must do some cleanup, othw the inputGenerators are kept alive from an iteration to the next
		b.ResetTimer()                                          //to skip setups
		for i := 0; i < b.N; i++ {
			//fmt.Printf("testing again %d\n", i)
			ConvergeCountTrigger(mons)
			//Tickn(mons, ticks)
		}
	}
}

//func BenchmarkClique10000(b *testing.B) { benchmarkMain(10000, b) }
func BenchmarkClique5000(b *testing.B) { benchmarkMain(5000, b) }
func BenchmarkClique2000(b *testing.B) { benchmarkMain(2000, b) }
func BenchmarkClique1000(b *testing.B) { benchmarkMain(1000, b) }
func BenchmarkClique100(b *testing.B)  { benchmarkMain(100, b) }
func BenchmarkClique10(b *testing.B)   { benchmarkMain(10, b) }

//func BenchmarkClique250(b *testing.B)   { benchmarkMain(250, b) }
//func BenchmarkClique300(b *testing.B)   { benchmarkMain(300, b) }
//func BenchmarkClique400(b *testing.B)   { benchmarkMain(400, b) }
//func BenchmarkClique500(b *testing.B)  { benchmarkMain(500, b) }

//func BenchmarkClique50000(b *testing.B) { benchmarkMain(50000, b) }
//func BenchmarkClique100000(b *testing.B) { benchmarkMain(100000, b) }

//go test -bench=. -benchmem

//go test -coverprofile cover.out -cpuprofile cpu.out -memprofile mem.out
//go tool pprof cpu.out
