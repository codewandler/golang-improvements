package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep devirt
// After taking &w, w.Work() is NOT devirtualized

type Worker interface {
	Work() string
}

type FastWorker struct{}

func (f FastWorker) Work() string { return "done" }

//go:noinline
func inspectInterface(p *Worker) string { return (*p).Work() }

//go:noinline
func Process() string {
	var w Worker = FastWorker{}
	inspectInterface(&w)     // takes address of w → sets Addrtaken
	return w.Work()          // NOT devirtualized after &w
}
