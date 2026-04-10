package p

// Reproduce: go tool compile -d=ssa/check_bce reproduce.go
// Expected: no IsInBounds
// Actual:   Found IsInBounds at BOTH a[i] and a[i+1]

func pairAccess(a []int, i int) int {
	if i >= 0 && i+1 < len(a) {
		return a[i] + a[i+1] // BOTH BOUNDS CHECKS REMAIN
	}
	return 0
}

// Compare: this equivalent formulation has ZERO bounds checks
func pairAccessFixed(a []int, i int) int {
	if i >= 1 && i < len(a) {
		return a[i-1] + a[i] // NO bounds checks
	}
	return 0
}
