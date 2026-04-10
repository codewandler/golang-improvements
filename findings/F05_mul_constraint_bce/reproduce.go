package p

// Reproduce: go tool compile -d=ssa/check_bce reproduce.go
// Expected: no IsInBounds
// Actual:   Found IsInBounds at both a[i*2] and a[i*2+1]

func mulBounds(a []int, i int) int {
	if i >= 0 && i*2 < len(a) {
		return a[i*2] + a[i*2+1] // BOUNDS CHECKS REMAIN
	}
	return 0
}
