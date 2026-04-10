package main

// Test 2: Redundant double zero/sign extensions  
// ZeroExt16to64(ZeroExt8to16(x)) => ZeroExt8to64(x)
// SignExt16to64(SignExt8to16(x)) => SignExt8to64(x)

//go:noinline
func doubleZeroExt(x uint8) uint64 {
	return uint64(uint16(x)) // Should be one MOVBQZX, not two extends
}

//go:noinline
func doubleSignExt(x int8) int64 {
	return int64(int16(x)) // Should be one MOVBQSX, not two extends
}

//go:noinline
func doubleZeroExt32(x uint8) uint32 {
	return uint32(uint16(x)) // Should be one MOVBQZX
}

//go:noinline
func tripleZeroExt(x uint8) uint64 {
	return uint64(uint32(uint16(x))) // Should be one MOVBQZX
}

func main() {
	println(doubleZeroExt(42))
	println(doubleSignExt(-42))
	println(doubleZeroExt32(42))
	println(tripleZeroExt(42))
}
