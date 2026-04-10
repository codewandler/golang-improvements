package absorption_test

// Benchmark: boolean absorption pattern performance.
//
// Run BEFORE fix:  go test -bench=. -benchtime=3s -count=6 | tee bench_before.txt
// Run AFTER fix:   go test -bench=. -benchtime=3s -count=6 | tee bench_after.txt
// Compare:         benchstat bench_before.txt bench_after.txt
//
// Tested: go1.26.1 linux/amd64

import "testing"

// --- Functions under test (noinline to isolate the pattern) ---

//go:noinline
func andAbsorb64(x, y uint64) uint64 { return x & (x | y) }

//go:noinline
func orAbsorb64(x, y uint64) uint64 { return x | (x & y) }

//go:noinline
func andAbsorb32(x, y uint32) uint32 { return x & (x | y) }

//go:noinline
func orAbsorb32(x, y uint32) uint32 { return x | (x & y) }

// Optimal baseline — what the compiler should generate after the fix
//go:noinline
func identity64(x, y uint64) uint64 { return x }

//go:noinline
func identity32(x, y uint32) uint32 { return x }

// --- More realistic: absorption inside a hot loop ---

//go:noinline
func loopAbsorb64(a, b []uint64) {
	for i := range a {
		a[i] = a[i] & (a[i] | b[i])
	}
}

//go:noinline
func loopIdentity64(a, b []uint64) {
	for i := range a {
		_ = b[i] // keep same memory access pattern
		a[i] = a[i]
	}
}

// --- Benchmarks ---

func BenchmarkAndAbsorb64(b *testing.B) {
	x, y := uint64(0xDEADBEEF), uint64(0xCAFEBABE)
	for b.Loop() {
		_ = andAbsorb64(x, y)
	}
}

func BenchmarkIdentity64(b *testing.B) {
	x, y := uint64(0xDEADBEEF), uint64(0xCAFEBABE)
	for b.Loop() {
		_ = identity64(x, y)
	}
}

func BenchmarkOrAbsorb64(b *testing.B) {
	x, y := uint64(0xDEADBEEF), uint64(0xCAFEBABE)
	for b.Loop() {
		_ = orAbsorb64(x, y)
	}
}

func BenchmarkAndAbsorb32(b *testing.B) {
	x, y := uint32(0xDEADBEEF), uint32(0xCAFEBABE)
	for b.Loop() {
		_ = andAbsorb32(x, y)
	}
}

func BenchmarkIdentity32(b *testing.B) {
	x, y := uint32(0xDEADBEEF), uint32(0xCAFEBABE)
	for b.Loop() {
		_ = identity32(x, y)
	}
}

func BenchmarkOrAbsorb32(b *testing.B) {
	x, y := uint32(0xDEADBEEF), uint32(0xCAFEBABE)
	for b.Loop() {
		_ = orAbsorb32(x, y)
	}
}

// Loop benchmark — more realistic
func BenchmarkLoopAbsorb64(b *testing.B) {
	a := make([]uint64, 1024)
	c := make([]uint64, 1024)
	for i := range a {
		a[i] = uint64(i)
		c[i] = uint64(i * 3)
	}
	b.ResetTimer()
	for b.Loop() {
		loopAbsorb64(a, c)
	}
}

func BenchmarkLoopIdentity64(b *testing.B) {
	a := make([]uint64, 1024)
	c := make([]uint64, 1024)
	for i := range a {
		a[i] = uint64(i)
		c[i] = uint64(i * 3)
	}
	b.ResetTimer()
	for b.Loop() {
		loopIdentity64(a, c)
	}
}
