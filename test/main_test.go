package main

import (
	"testing"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getSpec("inputMonitor.txt", "past", "trigger", "clique", 2, b.N)
		//b.ResetTimer() //to skip setups
	}
}
