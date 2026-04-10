# D04: Address-Taken Interface Variables Block Devirtualization

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Devirtualization |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard |
| **Impact** | Low-Medium |

## Problem

When an interface variable has its address taken (`&w`), devirtualization bails out even if the address is only used for reading.

## Location

`src/cmd/compile/internal/devirtualize/devirtualize.go:285-286`
