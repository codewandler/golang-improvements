package main

// Test 3: Various algebraic simplifications

//go:noinline
func andSelf(x uint64) uint64 {
	return x & x // Should be x (already handled)
}

//go:noinline
func orSelf(x uint64) uint64 {
	return x | x // Should be x (already handled)
}

// Test: (x | y) & (x | ^y) => x
// This is a consensus/resolution identity
//go:noinline
func consensus(x, y uint64) uint64 {
	return (x | y) & (x | ^y)
}

// Test: (x & y) | (x & ^y) => x
//go:noinline
func consensus2(x, y uint64) uint64 {
	return (x & y) | (x & ^y)
}

// Test: x ^ (x & y) => x & ^y  (XOR-AND simplification)
//go:noinline
func xorAndSimplify(x, y uint64) uint64 {
	return x ^ (x & y)
}

// Test: (x | y) ^ y => x & ^y
//go:noinline
func orXorSimplify(x, y uint64) uint64 {
	return (x | y) ^ y
}

// Test: And of Or with same operand: (x & y) | x => x (absorption)
//go:noinline
func orAndAbsorb(x, y uint64) uint64 {
	return (x & y) | x
}

func main() {
	println(consensus(0xF0, 0x0F))
	println(consensus2(0xF0, 0x0F))
	println(xorAndSimplify(0xFF, 0x0F))
	println(orXorSimplify(0xFF, 0x0F))
	println(orAndAbsorb(0xFF, 0x0F))
}
