package p

// Reproduce: go build -gcflags='-m -m' . 2>&1 | grep devirt
// Expected: devirtualizing s.Speak to Dog in AnimalSound
// Actual:   only devirtualized AFTER inlining the caller (not in AnimalSound itself)

type Speaker interface {
	Speak() string
}

type Dog struct{}

func (d Dog) Speak() string { return "Woof" }

// s is PPARAM — devirtualize.go:272 bails out for non-PAUTO
//
//go:noinline
func AnimalSound(s Speaker) string {
	return s.Speak() // NOT devirtualized (s is a parameter)
}

// Compare: local variable IS devirtualized
//
//go:noinline
func AnimalSoundDirect() string {
	var s Speaker = Dog{}
	return s.Speak() // devirtualized
}
