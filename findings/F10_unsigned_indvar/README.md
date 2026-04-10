# F10: Induction Variable Detection Ignores Unsigned Comparisons

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Loop BCE |
| **Origin** | TODO — `loopbce.go:129` |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Medium — affects `for i := uint(0); i < n; i++` loops |
| **Security** | ⚪ none |
| **Related issues** | #26116 |

## Problem

`findIndVar` only recognizes signed comparison ops (`OpLess64`, `OpLeq64`, etc.). Unsigned loop variables (`uint`) can't benefit from loop-based bounds-check elimination.

## Location

`src/cmd/compile/internal/ssa/loopbce.go:129` — *"TODO: Handle unsigned comparisons?"*
