# F02: Modulo + Fixup Pattern Not Proven for BCE

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Bounds-check elimination (prove pass) |
| **Origin** | TODO — derived from FIXME at `prove.go:2456` |
| **Status** | ✅ CONFIRMED with `check_bce` |
| **Difficulty** | Hard |
| **Impact** | High — common hash table / ring buffer pattern |

## Problem

```go
func modBoundsCheck(a []int, i int) int {
    if len(a) > 0 {
        idx := i % len(a)
        if idx < 0 { idx += len(a) }
        return a[idx] // ← BOUNDS CHECK NOT ELIMINATED
    }
    return 0
}
```

After `idx = i % len(a)` followed by the negative fixup, the compiler should know `0 <= idx < len(a)` but can't prove it.

## Reproduction

```bash
go tool compile -d=ssa/check_bce reproduce.go 2>&1 | grep IsInBounds
```

## Root Cause

The prove pass doesn't derive signed facts for subtraction/addition results. The FIXME at `prove.go:2456` says: *"we could also do signed facts but the overflow checks are much trickier and I don't need it yet."*

## Location

`src/cmd/compile/internal/ssa/prove.go:2456` (FIXME)
