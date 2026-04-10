package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep devirt
// Expected: devirtualizing c.Close
// Actual:   no devirtualization (GoDefer = true)

type Closer interface {
	Close() error
}

type MyCloser struct{}

func (m *MyCloser) Close() error { return nil }

//go:noinline
func DeferClose() {
	var c Closer = &MyCloser{}
	defer c.Close() // NOT devirtualized — GoDefer check at devirtualize.go:36-38
}
