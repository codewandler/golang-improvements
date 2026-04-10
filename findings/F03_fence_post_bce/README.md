# F03: `i+1 < len(a)` Doesn't Prove `a[i]` Safe

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Bounds-check elimination (prove pass) |
| **Origin** | NEW — independently discovered |
| **Status** | ✅ CONFIRMED with `check_bce` + benchmark |
| **Tested on** | go1.26.1 linux/amd64 (tip at 4478774aa2) |
| **Lifecycle** | in-progress |
| **Related issues** | None found (source, issue tracker, Gerrit all checked) |
| **Difficulty** | Medium |
| **Impact** | High — ~10% slowdown on pair-access loops, extremely common pattern |
| **Security** | 🔍 watch area — touches prove pass (safety-proving code) |

## Problem

```go
func pairAccess(a []int, i int) int {
    if i >= 0 && i+1 < len(a) {
        return a[i] + a[i+1] // ← BOTH BOUNDS CHECKS REMAIN
    }
    return 0
}
```

The equivalent `i >= 1 && i < len(a)` with `a[i-1] + a[i]` has ZERO bounds checks.

## Reproduction

```bash
go tool compile -d=ssa/check_bce reproduce.go 2>&1 | grep IsInBounds
# Shows: Found IsInBounds at a[i] AND a[i+1]
```

## Benchmark

```
BenchmarkSumPairsWithBCE    389 ns/op   (bounds checks remain)
BenchmarkSumPairsManualBCE  352 ns/op   (bounds checks eliminated)
```

~10% slower on 1000-element slice.

## Root Cause

The prove pass sees `i+1 < len(a)` but doesn't derive `i < len(a)` transitively (since `i < i+1` and `i+1 < len(a)` ⟹ `i < len(a)`). The fence-post derivation at `prove.go:1097` handles the reverse direction but not this one.

## Location

`src/cmd/compile/internal/ssa/prove.go` — fence-post section at line ~1097.

## Investigation Status

See [INVESTIGATION.md](INVESTIGATION.md) for detailed analysis. The naive fix
(adding `v > x+1 ⇒ v > x` to the fence-post switch) is correct but doesn't fire
because `i.max` hasn't been tightened by limit propagation yet when the fence-post
code runs. The fix needs the `update` function's internal ordering to be adjusted.
