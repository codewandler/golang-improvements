# D03: Generic Shape Types Block Devirtualization

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Devirtualization |
| **Origin** | TODO — `devirtualize.go:84-93` |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard |
| **Impact** | Medium — increasingly relevant as generics adoption grows |
| **Security** | ⚪ none |

## Problem

When a concrete type has shape types (from generic instantiation), devirtualization is completely blocked. Two TODOs describe potential fixes.

## Location

`src/cmd/compile/internal/devirtualize/devirtualize.go:84-93`
