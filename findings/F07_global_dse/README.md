# F07: Dead Store Elimination is Basic-Block-Only

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Dead code elimination |
| **Origin** | TODO — `deadstore.go:20` |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard — requires cross-block dataflow |
| **Impact** | Medium |
| **Security** | ⚪ none |
| **Related issues** | #67957, #26153 |

## Problem

Stores killed in successor blocks survive because DSE only looks within one basic block.

## Location

`src/cmd/compile/internal/ssa/deadstore.go:20` — *"TODO: use something more global"*
