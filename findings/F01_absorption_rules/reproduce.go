package main

// Test 1: Boolean absorption laws
// x & (x | y) => x
// x | (x & y) => x
// These are standard boolean algebra identities missing from generic.rules

//go:noinline
func andAbsorption(x, y uint64) uint64 {
	return x & (x | y) // Should simplify to just x
}

//go:noinline
func orAbsorption(x, y uint64) uint64 {
	return x | (x & y) // Should simplify to just x
}

//go:noinline
func andAbsorption32(x, y uint32) uint32 {
	return x & (x | y) // Should simplify to just x
}

//go:noinline
func orAbsorption32(x, y uint32) uint32 {
	return x | (x & y) // Should simplify to just x
}

func main() {
	println(andAbsorption(0xFF, 0x0F))
	println(orAbsorption(0xFF, 0x0F))
	println(andAbsorption32(0xFF, 0x0F))
	println(orAbsorption32(0xFF, 0x0F))
}
