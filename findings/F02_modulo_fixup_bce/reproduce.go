package p

// Reproduce: go tool compile -d=ssa/check_bce reproduce.go
// Expected: no IsInBounds on the return line
// Actual:   Found IsInBounds

func modBoundsCheck(a []int, i int) int {
	if len(a) > 0 {
		idx := i % len(a)
		if idx < 0 {
			idx += len(a)
		}
		return a[idx] // BOUNDS CHECK NOT ELIMINATED
	}
	return 0
}
