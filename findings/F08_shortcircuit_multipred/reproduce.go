// Test: shortcircuit missed case - three preds joining
package main

//go:noinline
func cond1() bool { return true }
//go:noinline
func cond2() bool { return true }  
//go:noinline  
func cond3() bool { return true }

// This creates a pattern where b has 3 preds and the shortcircuit
// optimization bails out (line 133: nOtherPhi > 0 && len(b.Preds) != 2)
func threeWayJoin(a, b, c int) int {
	x := 0
	switch {
	case a > 0:
		x = 1
	case b > 0:
		x = 2
	case c > 0:
		x = 3
	}
	if x > 0 {
		return x * 2
	}
	return 0
}

// Pattern with phi used as control but additional phis in block
// with more than 2 predecessors - this is missed
func multiPredPhi(a, b, c bool, x, y, z int) int {
	var val int
	var flag bool
	if a {
		val = x
		flag = true
	} else if b {
		val = y
		flag = true
	} else {
		val = z
		flag = c
	}
	// At this join point, both val and flag are phis.
	// If flag is a phi used as control with a const true arg,
	// shortcircuit would like to redirect, but it needs > 2 preds
	// handling for the val phi.
	if flag {
		return val
	}
	return -1
}

func main() {
	println(threeWayJoin(1, 2, 3))
	println(multiPredPhi(false, false, true, 10, 20, 30))
}
