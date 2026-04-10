# D01: Interface Parameters (PPARAM) Never Devirtualized

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Devirtualization |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Medium |
| **Impact** | Medium — affects any function accepting an interface parameter |
| **Security** | ⚪ none |

## Problem

`concreteType1()` in devirtualize.go bails out for non-`PAUTO` variables. Function parameters (`PPARAM`) with interface types are never devirtualized statically.

Note: After inlining the caller, the interface parameter becomes a local and CAN be devirtualized. So the impact is limited to non-inlined call sites.

## Location

`src/cmd/compile/internal/devirtualize/devirtualize.go:272` — `if name.Class != ir.PAUTO { return nil }`
