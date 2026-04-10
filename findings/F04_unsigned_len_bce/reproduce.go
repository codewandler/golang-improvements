package p

// Reproduce: go tool compile -d=ssa/check_bce reproduce.go
// Expected: no IsInBounds
// Actual:   Found IsInBounds at both a[i] and a[i+1]

func unsignedAfterLen(a []int, i uint) int {
	if i < uint(len(a))-1 {
		return a[i] + a[i+1] // BOUNDS CHECKS REMAIN
	}
	return 0
}
