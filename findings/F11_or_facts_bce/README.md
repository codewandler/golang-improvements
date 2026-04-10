# F11: OR-Based Unsigned Facts Commented Out

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | BCE (prove pass) |
| **Origin** | KNOWN — issue #57959 |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard — adding facts caused compile-time slowdown |
| **Impact** | Medium |
| **Security** | 🔍 watch area — touches prove pass (safety-proving code) |
| **Related issues** | #57959 |

## Problem

The fact `x | y >= x` (unsigned) is commented out because it caused slowdowns:

```go
case OpOr64, OpOr32, OpOr16, OpOr8:
    // TODO: investigate how to always add facts without much slowdown
    //ft.update(b, v, v.Args[0], unsigned, gt|eq)
    //ft.update(b, v, v.Args[1], unsigned, gt|eq)
```

## Location

`src/cmd/compile/internal/ssa/prove.go:2466-2469`
