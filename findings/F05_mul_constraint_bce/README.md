# F05: Multiplication Constraints Not Tracked for BCE

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | BCE (prove pass) |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard |
| **Impact** | Medium — common in SIMD-like manual unrolling |
| **Security** | 🔍 watch area — touches prove pass (safety-proving code) |

## Problem

```go
func mulBounds(a []int, i int) int {
    if i >= 0 && i*2 < len(a) {
        return a[i*2] + a[i*2+1] // ← BOUNDS CHECKS REMAIN
    }
    return 0
}
```

## Location

`src/cmd/compile/internal/ssa/prove.go` — no multiplication-based fact derivation.
