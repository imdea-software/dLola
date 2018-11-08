package main

import (
	"testing"
)

func BenchmarkMain(b *testing.B) {
	spec, ok := getSpec("inputMonitor.txt")
	if ok {
		b.ResetTimer() //to skip setups
		for i := 0; i < b.N; i++ {
			buildMonitor(spec, "past", "trigger", "clique", 2, b.N)
		}
	}
}

//go test -coverprofile cover.out -cpuprofile cpu.out -memprofile mem.out
//go tool pprof cpu.out
