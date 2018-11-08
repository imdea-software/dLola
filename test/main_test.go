package dLola

import (
	//dLola "gitlab.software.imdea.org/luismiguel.danielsson/dLola"
	"testing"
)

func benchmarkMain(nmons, tlen int, b *testing.B) {
	spec, ok := dLola.GetCheckedSpec("inputMonitor.txt")
	if ok {
		mons := dLola.BuildMonitorTopo(spec, "past", "trigger", "clique", nmons, tlen)
		b.ResetTimer() //to skip setups
		for i := 0; i < b.N; i++ {
			//verdict :=
			dLola.ConvergeCountTrigger(mons)
		}
	}
}

func BenchmarkMain1100(b *testing.B) { benchmarkMain(1, 100, b) }

//func BenchmarkMain11000(b *testing.B) { benchmarkMain(1, 1000, b) }
//func BenchmarkMain12000(b *testing.B) { benchmarkMain(1, 2000, b) }

//func BenchmarkMain2100(b *testing.B)  { benchmarkMain(2, 100, b) }
//func BenchmarkMain21000(b *testing.B) { benchmarkMain(2, 1000, b) }
//func BenchmarkMain22000(b *testing.B) { benchmarkMain(2, 2000, b) }

//go test -bench=.

//go test -coverprofile cover.out -cpuprofile cpu.out -memprofile mem.out
//go tool pprof cpu.out
