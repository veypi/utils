package utils

import "testing"

func BenchmarkRandSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandSeq(32)
	}
}

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Rand(32)
	}
}
