package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep 'moved to heap'
// Expected: nothing moved to heap (all local)
// Actual:   "moved to heap: i" — even though slice s is local-only

func slicePointersLocal() {
	var s []*int
	i := 0
	s = append(s, &i) // i ESCAPES — even though s never leaves the function
	_ = s
}

func mapValueLocal() int {
	x := 42
	m := map[string]*int{"a": &x} // x ESCAPES — even though m is local
	return *m["a"]
}

func chanBufferedLocal() int {
	x := 42
	ch := make(chan *int, 1)
	ch <- &x // x ESCAPES — even though ch is local and buffered
	return *<-ch
}
