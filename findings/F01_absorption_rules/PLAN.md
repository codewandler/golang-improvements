# F01 Fix Plan: Boolean Absorption Rules

## Overview

Add boolean absorption laws to the Go compiler's SSA generic rewrite rules. This eliminates redundant AND/OR instruction pairs when the compiler encounters patterns like `x & (x | y)` or `x | (x & y)`.

**Branch**: `fix/F01-absorption-rules`  
**Files to modify**: 3  
**Files to create**: 0  
**Estimated time**: 30-45 min  

---

## Proof Strategy

We prove the fix works with **three layers of evidence**:

### 1. Assembly diff (before → after)

Captured in `asm_before.txt`. After the fix, regenerate as `asm_after.txt` and diff:

```bash
# Capture after
go tool compile -S reproduce.go 2>&1 | ... > asm_after.txt
# Diff
diff asm_before.txt asm_after.txt
```

**Expected change per function**:
```
BEFORE:  ORQ + ANDQ + RET   (3 instructions, 7 bytes)
AFTER:   RET                 (1 instruction, 1 byte)
```

Total across 6 test functions: 38 bytes → 6 bytes.

### 2. Benchmark (before → after)

Baseline captured in `bench/bench_before.txt`. After the fix:

```bash
cd bench && go test -bench=. -benchtime=3s -count=6 | tee bench_after.txt
benchstat bench_before.txt bench_after.txt
```

**Baseline numbers** (go1.26.1, Intel i9-10900K):

| Benchmark | Before | Optimal | Gap |
|-----------|--------|---------|-----|
| Single call (64-bit) | 1.05 ns/op | 1.06 ns/op | ~0% (call overhead dominates) |
| Loop 1024 elements | **441.2 ns/op** | **228.6 ns/op** | **-48.2%** (p=0.002) |

The single-call benchmark is dominated by function call overhead (~1ns), so the 2 extra ALU ops are invisible there. The **loop benchmark** is the real proof — it isolates the pattern in a tight loop where the redundant instructions are a significant fraction of the work.

After the fix, `LoopAbsorb64` should drop to match `LoopIdentity64` (~229 ns/op).

### 3. Codegen tests (upstream regression guard)

Added to `test/codegen/bits.go` — these are negative assembly assertions that run in CI:

```go
func absorptionAnd64(x, y uint64) uint64 {
    // amd64:-"ORQ" -"ANDQ"
    return x & (x | y)
}
```

The `-"ORQ"` syntax means "fail if ORQ appears in the output." This prevents future regressions.

---

## Tasks

### Task 1: Create fix branch

```bash
cd go-src
git checkout -b fix/F01-absorption-rules master
```

**Verification**: `git branch --show-current` → `fix/F01-absorption-rules`

---

### Task 2: Add absorption rules to generic.rules

**File**: `src/cmd/compile/internal/ssa/_gen/generic.rules`  
**Location**: After DeMorgan's Laws (lines 208-209), add:

```
// Absorption laws
(And(8|16|32|64) x (Or(8|16|32|64) x y)) => x
(Or(8|16|32|64)  x (And(8|16|32|64) x y)) => x
```

That's 2 lines. The `(8|16|32|64)` syntax covers all 4 widths. Commutativity of AND/OR is handled automatically by the rule engine — it will generate matches for all argument orderings:
- `(And64 x (Or64 x y))` ✓
- `(And64 x (Or64 y x))` ✓ (Or64 is commutative)
- `(And64 (Or64 x y) x)` ✓ (And64 is commutative)
- `(And64 (Or64 y x) x)` ✓ (both commutative)

**Verification**: Rules file parses cleanly (checked in Task 3).

---

### Task 3: Regenerate rewritegeneric.go

```bash
cd go-src/src/cmd/compile/internal/ssa
go generate
```

This runs `_gen/rulegen.go` which reads `generic.rules` and regenerates `rewritegeneric.go`. The generated file should show new `rewriteValuegeneric_OpAnd64`, `rewriteValuegeneric_OpOr64` (etc.) cases containing the absorption patterns.

**Verification**: 
```bash
git diff --stat src/cmd/compile/internal/ssa/rewritegeneric.go
# Should show additions in the OpAnd{8,16,32,64} and OpOr{8,16,32,64} functions
```

---

### Task 4: Add codegen test

**File**: `test/codegen/bits.go` (extend existing bitwise operation tests)

```go
// Absorption: x & (x | y) => x
func absorptionAnd64(x, y uint64) uint64 {
	// amd64:-"ORQ" -"ANDQ"
	return x & (x | y)
}

func absorptionOr64(x, y uint64) uint64 {
	// amd64:-"ORQ" -"ANDQ"
	return x | (x & y)
}

func absorptionAnd32(x, y uint32) uint32 {
	// amd64:-"ORL" -"ANDL"
	return x & (x | y)
}

func absorptionOr32(x, y uint32) uint32 {
	// amd64:-"ORL" -"ANDL"
	return x | (x & y)
}
```

**Verification**:
```bash
cd go-src
go test cmd/compile/internal/ssa        # SSA unit tests pass
go run cmd/internal/testdir -run=codegen # codegen tests pass (includes our new tests)
```

---

### Task 5: Capture proof — assembly + benchmark AFTER fix

```bash
# Assembly after
go tool compile -S findings/F01_absorption_rules/reproduce.go 2>&1 | ... > asm_after.txt

# Benchmark after
cd findings/F01_absorption_rules/bench
go test -bench=. -benchtime=3s -count=6 | tee bench_after.txt
benchstat bench_before.txt bench_after.txt
```

**Expected**: LoopAbsorb64 drops from ~441ns to ~229ns (matching LoopIdentity64).

---

### Task 6: Run full test suite

```bash
cd go-src
go test cmd/compile/...              # all compiler tests
go test std                          # standard library (catches regressions)
```

---

## Deliverables

After the fix, the finding directory contains:

```
findings/F01_absorption_rules/
├── README.md             ← finding description + metadata
├── PLAN.md               ← this file
├── ISSUE.md              ← GitHub issue template (pre-filled)
├── reproduce.go          ← minimal reproducer (package p)
├── asm_before.txt        ← assembly BEFORE fix
├── asm_after.txt         ← assembly AFTER fix
├── fix.patch             ← human-authored changes ONLY (rules + tests)
└── bench/
    ├── bench_test.go     ← benchmark code
    ├── bench_before.txt  ← benchmark BEFORE fix (6 runs)
    └── bench_after.txt   ← benchmark AFTER fix (6 runs)
```

### About the patch file

`fix.patch` contains **only the human-authored changes**:
- `src/cmd/compile/internal/ssa/_gen/generic.rules` — 2 new rules (4 lines)
- `test/codegen/bits.go` — 6 codegen test functions (32 lines)

It deliberately **excludes** `rewritegeneric.go` because that file is auto-generated.

### Applying the patch

```bash
cd go-src
git apply ../findings/F01_absorption_rules/fix.patch   # apply human-authored changes
cd src/cmd/compile/internal/ssa && go generate          # regenerate rewritegeneric.go
cd ../../../../../src && ./make.bash                     # rebuild compiler
```

### What goes in the issue report

- Side-by-side assembly diff (asm_before.txt vs asm_after.txt)
- benchstat output showing the loop improvement (~48% expected)
- Link to fix branch: `https://github.com/codewandler/go/tree/fix/F01-absorption-rules`
- Link to the patch file

### Benchmark baseline interpretation

**Single-call benchmarks** (~1.05 ns/op for both absorb and identity):
No measurable difference — function call overhead (~1ns) dominates the 2 extra ALU ops (~0.3ns each).

**Loop benchmark** — the real proof:
| Benchmark | ns/op | Description |
|-----------|-------|-------------|
| LoopAbsorb64 | **441.2 ± 1%** | `a[i] = a[i] & (a[i] | b[i])` × 1024, with redundant ORQ+ANDQ |
| LoopIdentity64 | **228.6 ± 3%** | `a[i] = a[i]` × 1024, optimal baseline |
| **Gap** | **-48.2%** | p=0.002, n=6 — statistically significant |

After the fix, `LoopAbsorb64` should compile to the same code as `LoopIdentity64` and produce matching numbers.

---

## Risk Assessment

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Rule matches too aggressively | Very low — absorption is a mathematical identity | The pattern is unambiguous: same `x` on both sides |
| Performance regression (compile time) | Very low — 2 simple rules, no conditions | Monitor `go test -bench` on the compiler itself |
| Interacts badly with other rules | Very low — the output `x` is simpler than the input, so no infinite loops | The result has fewer ops than the pattern |

**Why this is safe**: The absorption law holds for ALL bit patterns. There's no edge case with signed values, overflow, or special constants. `x & (x | y)` is *always* `x`. The proof is one line of boolean algebra: `x & (x | y) = (x & x) | (x & y) = x | (x & y) = x`.
