# F08: Shortcircuit Only Handles 2-Predecessor Blocks

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Shortcircuit optimization |
| **Origin** | TODO — `shortcircuit.go:482` |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Medium-High — TODO says "reasonably high impact" |
| **Related issues** | #45175, #33903, #44465 |

## Problem

`a || b || c` creates a 3-predecessor join block. The shortcircuit pass bails out for >2 predecessors.

## Location

`src/cmd/compile/internal/ssa/shortcircuit.go:482`
