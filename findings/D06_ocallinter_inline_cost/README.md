# D06: OCALLINTER Always Charged Full extraCallCost in Inlining Budget

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Inlining |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Medium — interface-containing functions penalized even when devirtualizable |
| **Security** | ⚪ none |

## Problem

`OCALLFUNC` checks if the callee is inlinable and charges only `callee.Inl.Cost`. `OCALLINTER` always charges the full `extraCallCost` (57), preventing inlining of functions that contain interface calls even when those calls will be devirtualized.

## Location

`src/cmd/compile/internal/inline/inl.go:657-662`
