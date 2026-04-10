# D02: go/defer Calls Skip Devirtualization Entirely

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Devirtualization |
| **Origin** | KNOWN — issue #52072 |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Low-Medium — overly conservative for non-promoted methods |
| **Security** | ⚪ none |
| **Related issues** | #52072 |

## Problem

All `go` and `defer` interface calls are skipped to preserve panic semantics for promoted methods, even when the method is not promoted.

## Location

`src/cmd/compile/internal/devirtualize/devirtualize.go:36-38`
