# F01: Missing Boolean Absorption Laws

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | SSA rewrite rules |
| **Origin** | NEW — independently discovered |
| **Status** | ✅ CONFIRMED with compiler output |
| **Difficulty** | Easy (~30 min) |
| **Impact** | Medium — 2 redundant instructions per pattern occurrence |
| **Security** | ⚪ none |

## Problem

The Go compiler's SSA rewrite rules don't implement boolean absorption:
- `x & (x | y) == x`
- `x | (x & y) == x`

These are fundamental boolean algebra identities.

## Reproduction

```bash
go tool compile -S findings/F01_absorption_rules/reproduce.go 2>&1 | grep -A5 'andAbsorption(SB), NOSPLIT'
```

**Actual output**: `ORQ AX, BX; ANDQ BX, AX; RET` (3 instructions)  
**Expected output**: `RET` (1 instruction — just return x)

## Fix

Add to `src/cmd/compile/internal/ssa/generic.rules`:

```
(And64 x (Or64 x y)) => x
(And32 x (Or32 x y)) => x
(And16 x (Or16 x y)) => x
(And8  x (Or8  x y)) => x
(Or64  x (And64 x y)) => x
(Or32  x (And32 x y)) => x
(Or16  x (And16 x y)) => x
(Or8   x (And8  x y)) => x
```

## Location

`src/cmd/compile/internal/ssa/generic.rules` — rules missing entirely.
