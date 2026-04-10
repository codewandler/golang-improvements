# D05: Struct-Field and Package-Level Interface Calls Not Devirtualized

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | Devirtualization |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Hard |
| **Impact** | Medium — common pattern in server/handler architectures |

## Problem

Only `PAUTO` (local auto) variables are tracked for devirtualization. Struct field accesses (`s.handler.ServeHTTP()`) and package-level variables fail the `ONAME` check.

## Location

`src/cmd/compile/internal/devirtualize/devirtualize.go:267-268`
