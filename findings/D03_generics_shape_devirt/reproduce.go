package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep devirt
// Shows that generic instantiations block devirtualization

type Stringer interface {
	String() string
}

type MyString struct{ Val string }

func (m MyString) String() string { return m.Val }

// Generic: interface calls through type params can't be devirtualized
func PrintIt[T Stringer](val T) string {
	return val.String() // NOT devirtualized — shape type blocks it
}

// Non-generic: CAN be devirtualized after inlining
func PrintItConcrete(val Stringer) string {
	return val.String()
}
