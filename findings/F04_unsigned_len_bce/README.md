# F04: Unsigned Arithmetic After Length Check — BCE Miss

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | BCE (prove pass) |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Medium |

## Problem

```go
func unsignedAfterLen(a []int, i uint) int {
    if i < uint(len(a))-1 {
        return a[i] + a[i+1] // ← BOTH BOUNDS CHECKS REMAIN
    }
    return 0
}
```

## Location

`src/cmd/compile/internal/ssa/prove.go` — unsigned domain fact propagation.
