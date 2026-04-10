package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep devirt
// s.handler.ServeHTTP() is NOT devirtualized — ODOT fails ONAME check

type Handler interface {
	ServeHTTP() string
}

type MyHandler struct{}

func (m *MyHandler) ServeHTTP() string { return "served" }

type Server struct {
	handler Handler
}

//go:noinline
func (s *Server) Run() string {
	return s.handler.ServeHTTP() // NOT devirtualized (field access, not ONAME)
}

// Compare: local variable IS devirtualized
//go:noinline
func RunDirect() string {
	var h Handler = &MyHandler{}
	return h.ServeHTTP() // devirtualized
}
