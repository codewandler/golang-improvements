# F01: Missing Boolean Absorption Laws

| Field | Value |
|-------|-------|
| **Category** | performance |
| **Sub-area** | ssa-rules |
| **Origin** | NEW |
| **Status** | CONFIRMED |
| **Difficulty** | Easy |
| **Impact** | Low — 9 hits in stdlib, all in cold paths (defer bit-manipulation) |
| **Security** | none |
| **Tested on** | go1.26.1 linux/amd64, tip at `4478774aa2` |
| **Lifecycle** | submitted — maintainer approved ✅ |
| **Related issues** | https://github.com/golang/go/issues/78632, [CL 736541](https://go-review.googlesource.com/c/go/+/736541), [CL 739720](https://go-review.googlesource.com/c/go/+/739720) |

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

## Real-World Impact (Corrected)

Our original instruction-count analysis was flawed — it conflated register
allocation cascading effects with actual rule firings. Using @Jorropo's `jlog`
technique, the actual numbers are much smaller:

| Project | Originally claimed | Actual (project-specific) |
|---------|-------------------:|--------------------------:|
| stdlib | -88 | 9 |
| hugo | -260 | 2 |
| prometheus | -70 | 18 |
| lazygit | -27 | 0 |

**All hits are `and absorb`. Zero `or absorb` hits anywhere.** See
[ISSUE_followup_01.md](ISSUE_followup_01.md) for the full location-by-location
breakdown.

### How the pattern actually arises

@Jorropo analyzed every stdlib hit and found they all come from **defer bit
manipulation**: the compiler optimizes `(x | 2) & 2 == 0` → `false` via
absorption for the last always-taken defer in a function. The pattern is:

- `cgocallbackg1` — 2nd defer always executes
- `(*common).runCleanup` — 3rd defer always executes
- `(*parser).parsePrimaryExpr` — defers in type-switch branches
- `(*parser).parseBinaryExpr` — same

The rules "work by accident" on this defer pattern rather than on source-level
`x & (x | y)`. Jorropo noted a broader optimization (tracking always-set /
always-cleared bits) could handle all always-taken defers, not just the last one.

## Upstream Status

- **Issue**: [golang/go#78632](https://github.com/golang/go/issues/78632)
- **randall77** (Keith Randall): "Triggering at all is good enough for me. Thanks." ✅
- **Next step**: Submit CL via Gerrit

## Lessons Learned

1. **Instruction count deltas ≠ rule firings.** Counting AND/OR instructions
   before/after conflates regalloc cascading with actual optimizations. Use
   `jlog` or `-genLog` to count real rule applications.
2. **Check Gerrit CLs before filing.** [CL 736541](https://go-review.googlesource.com/c/go/+/736541)
   proposed similar rules and got the exact pushback we received. We could have
   prepared real-code evidence upfront.
3. **Honest correction builds trust.** Posting corrected numbers with the flawed
   methodology acknowledged got immediate approval from the maintainer.

## Location

`src/cmd/compile/internal/ssa/_gen/generic.rules` — rules missing entirely.
