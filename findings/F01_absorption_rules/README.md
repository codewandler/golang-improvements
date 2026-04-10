# F01: Missing Boolean Absorption Laws

| Field | Value |
|-------|-------|
| **Category** | performance |
| **Sub-area** | ssa-rules |
| **Origin** | NEW |
| **Status** | CONFIRMED |
| **Difficulty** | Easy |
| **Impact** | Medium — 2 redundant instructions per pattern occurrence |
| **Security** | none |
| **Tested on** | go1.26.1 linux/amd64, tip at `4478774aa2` |
| **Lifecycle** | submitted |
| **Related issues** | https://github.com/golang/go/issues/78632 |

## Problem

The Go compiler's SSA rewrite rules implement DeMorgan's laws but are missing the closely related [boolean absorption laws](https://en.wikipedia.org/wiki/Absorption_law):
- `x & (x | y) == x`
- `x | (x & y) == x`

These are fundamental boolean algebra identities. Both GCC and LLVM recognize and optimize these patterns at `-O2`.

## Reproduction

```bash
go tool compile -S findings/F01_absorption_rules/reproduce.go 2>&1 | grep -E 'ORQ|ANDQ|RET'
```

**Actual output**: `ORQ AX, BX; ANDQ BX, AX; RET` (3 instructions, 7 bytes)  
**Expected output**: `RET` (1 instruction, 1 byte — just return x)

Same issue affects 32/16/8-bit variants and all commuted argument orderings.

## Fix

Add to `src/cmd/compile/internal/ssa/_gen/generic.rules` (next to DeMorgan's at line 208):

```
// Absorption laws
(And(8|16|32|64) x (Or(8|16|32|64) x y)) => x
(Or(8|16|32|64)  x (And(8|16|32|64) x y)) => x
```

Commutativity is handled automatically by the rule engine.

## Proof

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Assembly (64-bit) | ORQ+ANDQ+RET (7 bytes) | RET (1 byte) | -86% code size |
| Assembly (32-bit) | ORL+ANDL+RET (5 bytes) | RET (1 byte) | -80% code size |
| Benchmark (loop 1024) | 441.2 ns/op | 226.9 ns/op | -48.57% (p=0.002) |
| Compilebench (math/big) | 147.7 ms | 148.0 ms | ~ (p=0.699) |

## Location

`src/cmd/compile/internal/ssa/_gen/generic.rules` — rules missing entirely.
