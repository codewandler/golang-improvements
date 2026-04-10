package p

// Reproduce: go tool compile -m reproduce.go 2>&1 | grep -E 'inline|cost'
// OCALLINTER charges full extraCallCost=57 regardless of callee

type Sizer interface {
	Size() int
}

type Fixed struct{ N int }

func (f Fixed) Size() int { return f.N }

// Contains interface call → cost inflated by 57 for OCALLINTER
func TotalSize(items []Sizer) int {
	total := 0
	for _, item := range items {
		total += item.Size() // cost=57 as OCALLINTER
	}
	return total
}

// Same logic with concrete type → much cheaper
func TotalSizeFixed(items []Fixed) int {
	total := 0
	for _, item := range items {
		total += item.Size() // cost=inlined via OCALLFUNC
	}
	return total
}
