package dLola

import (
	//dLola "gitlab.software.imdea.org/luismiguel.danielsson/dLola"
	"testing"
)

var verdict Verdict //so optimizations do not get rid of the call to ConvergeCountTrigger

func benchmarkMain(nmons, tlen int, topo string, b *testing.B) {
	spec, ok := GetCheckedSpec("test/inputMonitor.txt")
	if ok {
		mons := BuildMonitorTopo(spec, "past", "trigger", topo, nmons, tlen)
		b.ResetTimer() //to skip setups
		for i := 0; i < b.N; i++ {
			verdict = ConvergeCountTrigger(mons)
		}
	}
}

func BenchmarkClique1100(b *testing.B) { benchmarkMain(1, 100, "clique", b) }

//func BenchmarkClique11000(b *testing.B) { benchmarkMain(1, 1000, b) }
//func BenchmarkClique12000(b *testing.B) { benchmarkMain(1, 2000, b) }

func BenchmarkClique2100(b *testing.B) { benchmarkMain(2, 100, "clique", b) }

//func BenchmarkClique21000(b *testing.B) { benchmarkMain(2, 1000, b) }
//func BenchmarkClique22000(b *testing.B) { benchmarkMain(2, 2000, b) }

func BenchmarkClique3100(b *testing.B) { benchmarkMain(3, 100, "clique", b) }

//go test -bench=.

//go test -coverprofile cover.out -cpuprofile cpu.out -memprofile mem.out
//go tool pprof cpu.out
