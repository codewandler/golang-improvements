package main

// Test 6: Shift and mask patterns, redundant conversions

// (x << n) >> n for unsigned should be AND with mask (already handled)
// But: (uint32(x) << 16) >> 16 should be zero-extend of lower 16 bits

//go:noinline
func andDistributesOverOr(x, y, z uint64) uint64 {
	// x & (y | z) is distributive: (x & y) | (x & z)
	// But the reverse: (x & y) | (x & z) => x & (y | z) is a strength reduction
	return (x & y) | (x & z)
}

//go:noinline
func orDistributesOverAnd(x, y, z uint64) uint64 {
	// (x | y) & (x | z) => x | (y & z)
	return (x | y) & (x | z)
}

// Test: Mul then Div by same constant - should cancel
//go:noinline
func mulDivCancel(x uint64) uint64 {
	return (x * 8) / 8 // Should be just x (if no overflow)
}

// Test: Shift left then shift right by same amount
// Already handled, but double-checking
//go:noinline
func shiftCancel(x uint64) uint64 {
	return (x << 3) >> 3 // Should be AND mask
}

// Test: Redundant zero extension after AND with small constant
//go:noinline
func redundantZext(x uint64) uint64 {
	y := x & 0xFF // Already fits in 8 bits
	return uint64(uint8(y)) // This ZeroExt should be eliminated
}

// Test: sign extension of already-narrow AND result
//go:noinline
func redundantSext(x int64) int64 {
	y := x & 0x7F // Already fits in 7 bits (positive)
	return int64(int8(y)) // Trunc+sext should be zext since value fits
}

func main() {
	println(andDistributesOverOr(0xFF, 0x0F, 0xF0))
	println(orDistributesOverAnd(0xFF, 0x0F, 0xF0))
	println(mulDivCancel(42))
	println(shiftCancel(42))
	println(redundantZext(42))
	println(redundantSext(42))
}
