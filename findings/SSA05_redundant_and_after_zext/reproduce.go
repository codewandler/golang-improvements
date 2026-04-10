package main

// Test 10: Miscellaneous missed optimizations

// Shift of zero-extended value should narrow
//go:noinline
func shiftZext(x uint8) uint64 {
	return uint64(x) << 4 // MOVBQZX then SHL, could be SHL (already narrow enough)
}

// AND with mask after zero extend is redundant if mask covers full width
//go:noinline
func andAfterZext(x uint8) uint64 {
	return uint64(x) & 0xFF // ZeroExt already guarantees <= 0xFF, AND is redundant
}

// OR of value with itself shifted - should just be value (no it shouldn't, this is wrong)
// Actually: x | (x << 0) == x | x == x
//go:noinline
func orShiftZero(x uint64) uint64 {
	return x | (x << 0) // Should simplify to x
}

// Test: y = x + 0 (already handled)
// Test: y = x * 0 (already handled)

// Test: Sub from itself
//go:noinline
func subSelf(x uint64) uint64 {
	return x - x // Should be 0
}

// Test: XOR with itself
//go:noinline
func xorSelf(x uint64) uint64 {
	return x ^ x // Should be 0
}

// Test: Double negation
//go:noinline
func doubleNeg(x int64) int64 {
	return -(-x) // Should be x
}

// Test: Neg(Sub(x,y)) = Sub(y,x) - already handled
//go:noinline
func negSub(x, y int64) int64 {
	return -(x - y) // Should be y - x
}

// Test: And with all ones
//go:noinline
func andAllOnes(x uint64) uint64 {
	return x & ^uint64(0) // Should be x
}

// Test: Or with zero
//go:noinline
func orZero(x uint64) uint64 {
	return x | 0 // Should be x
}

func main() {
	println(shiftZext(15))
	println(andAfterZext(42))
	println(orShiftZero(42))
	println(subSelf(42))
	println(xorSelf(42))
	println(doubleNeg(42))
	println(negSub(10, 3))
	println(andAllOnes(42))
	println(orZero(42))
}
