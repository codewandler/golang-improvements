package p

// Reproduce: go tool compile -S reproduce.go 2>&1 | grep -E 'ORQ|ANDQ|ORL|ANDL|RET'
// Expected:  only RET — no ORQ, ANDQ, ORL, or ANDL instructions
// Actual:    ORQ + ANDQ + RET per function (2 redundant instructions each)
// Tested:    go1.26.1 linux/amd64

// Boolean absorption: x & (x | y) == x
//
//go:noinline
func andAbsorption(x, y uint64) uint64 {
	return x & (x | y) // emits ORQ + ANDQ, should be eliminated
}

// Boolean absorption: x | (x & y) == x
//
//go:noinline
func orAbsorption(x, y uint64) uint64 {
	return x | (x & y) // emits ANDQ + ORQ, should be eliminated
}

// Same patterns at 32-bit width
//
//go:noinline
func andAbsorption32(x, y uint32) uint32 {
	return x & (x | y) // emits ORL + ANDL
}

//go:noinline
func orAbsorption32(x, y uint32) uint32 {
	return x | (x & y) // emits ANDL + ORL
}

// Commuted argument order — should also be optimized
//
//go:noinline
func andAbsorptionCommuted(x, y uint64) uint64 {
	return (x | y) & x // same as x & (x | y)
}

//go:noinline
func orAbsorptionCommuted(x, y uint64) uint64 {
	return (x & y) | x // same as x | (x & y)
}
