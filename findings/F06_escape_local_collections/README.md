# F06: False Heap Escape for Pointers in Local Collections

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Escape analysis |
| **Origin** | KNOWN-BAD — marked in upstream test suite |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard — fundamental EA limitation |
| **Impact** | High — causes heap allocations in common patterns |
| **Security** | 🔍 watch area — escape analysis correctness affects stack/heap safety |

## Problem

```go
func slicePointersLocal() {
    var s []*int
    i := 0
    s = append(s, &i) // i MOVED TO HEAP — even though s is local-only
    _ = s
}
```

Also affects: local maps (`map[string]*int{...}`), local buffered channels.

## Reproduction

```bash
go tool compile -m reproduce.go 2>&1 | grep 'moved to heap'
```

## Location

`src/cmd/compile/internal/escape/` — known limitation in Go's escape analysis.
