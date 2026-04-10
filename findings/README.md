# Findings Index

> Last updated: 2025-07-17 · Go version tested: go1.26.1 / tip at `2f3c778b23`

---

## Triage Views

Jump to what matters to you:

- **[🔴 Security-Relevant](#security-relevant)** — Could a bug here cause memory unsafety?
- **[🟠 Performance-Critical](#performance-critical)** — Measurable impact on real workloads
- **[🟡 Performance](#performance)** — Missed optimizations, smaller impact
- **[🔵 Correctness](#correctness)** — Compiler produces wrong results or fragile workarounds
- **[⚪ Code Quality](#code-quality)** — Compiler speed, maintainability

Or browse **[by sub-area](#by-sub-area)** (BCE, devirt, escape, etc.) or see the **[full table](#full-classification-table)**.

---

## Security-Relevant

> Findings where a bug could lead to **memory unsafety**, **buffer overread/overwrite**, or **incorrect safety guarantees**. None of our current findings are confirmed security bugs, but several touch security-sensitive code paths.

| ID | Summary | Why it matters | Severity |
|----|---------|---------------|----------|
| [F13](F13_scheduling_workaround/) | Fragile scheduling workaround in prove pass | The prove pass **removes bounds checks and nil checks**. A miscompilation here means the compiler removes a safety check that was needed → **buffer overread / nil deref**. The current workaround is correct but brittle. | ⚠️ Latent risk |
| [F02](F02_modulo_fixup_bce/) | Modulo + fixup BCE miss | Not a security bug itself (bounds check is *kept*, not *removed*). But bugs in this area of the prove pass *could* cause incorrect check elimination. | 🔍 Watch area |
| [F03](F03_fence_post_bce/) | Fence-post BCE miss | Same — the prove pass is conservative here (safe), but changes to fix this must be carefully reviewed to avoid over-proving. | 🔍 Watch area |

**How to evaluate security relevance of a new finding:**

| Question | If YES → |
|----------|----------|
| Does the compiler **remove** a safety check (bounds, nil, overflow) that it shouldn't? | 🔴 **Security bug** — file immediately |
| Does the compiler **keep** a safety check it could remove? | 🟡 Performance only — safe |
| Is the code path involved in **proving things safe**? (prove.go, nilcheck.go) | ⚠️ Watch area — any fix here needs extreme care |
| Does it affect **escape analysis** correctness? (stack vs heap) | ⚠️ Could cause use-after-free if EA says "stack" but value actually escapes |
| Is it a **missed optimization** in codegen? (rewrite rules, inlining, devirt) | ⚪ No security relevance |

---

## Performance-Critical

> Findings with **measured benchmark impact** or affecting **hot paths in real programs** (loops, slice access, interface dispatch). These are the ones worth upstreaming first.

| ID | Summary | Measured Impact | Difficulty |
|----|---------|----------------|------------|
| [F03](F03_fence_post_bce/) | `i+1 < len(a)` doesn't prove `a[i]` safe | **~10% slower** on pair-access loops (389 vs 352 ns/op) | Medium |
| [F06](F06_escape_local_collections/) | Local collection pointers false-escape to heap | **Heap allocation** on every call for common patterns | Hard |
| [F02](F02_modulo_fixup_bce/) | Modulo + fixup pattern: bounds check remains | Extra branch in every **hash table probe / ring buffer access** | Hard |
| [F08](F08_shortcircuit_multipred/) | Shortcircuit bails on >2 predecessors | Source TODO says **"reasonably high impact"** | Medium |
| [D05](D05_struct_field_devirt/) | Struct-field interfaces not devirtualized | **Indirect call** on every `s.handler.Method()` in hot server loops | Hard |
| [D06](D06_ocallinter_inline_cost/) | Interface calls inflate inlining budget | Functions with interface calls **miss inlining** even when devirtualizable | Medium |

---

## Performance

> Missed optimizations. Correct code, but slower than it could be. No measured benchmark yet — impact is inferred from pattern frequency and instruction count.

| ID | Summary | Origin | Status | Impact | Difficulty |
|----|---------|--------|--------|--------|------------|
| [F01](F01_absorption_rules/) | Missing boolean absorption laws (2 redundant instructions) | NEW | ✅ | Medium | **Easy** |
| [F04](F04_unsigned_len_bce/) | Unsigned arithmetic after length check — BCE miss | NEW | ✅ | Medium | Medium |
| [F05](F05_mul_constraint_bce/) | Multiplication constraints not tracked — BCE miss | NEW | ✅ | Medium | Hard |
| [F07](F07_global_dse/) | Dead store elimination is basic-block-only | TODO | ✅ | Medium | Hard |
| [F09](F09_phiopt_multipred/) | phiopt bails on >2 predecessors | TODO | ✅ | Low-Med | Medium |
| [F10](F10_unsigned_indvar/) | Unsigned induction vars ignored in loop BCE | TODO | ✅ | Medium | Medium |
| [F11](F11_or_facts_bce/) | OR unsigned facts commented out (#57959) | KNOWN | ✅ | Medium | Hard |
| [F12](F12_transitive_equality/) | Transitive equality not tracked | TODO | ❓ | Low | Medium |
| [D01](D01_pparam_devirt/) | Function params (PPARAM) never devirtualized | NEW | ✅ | Medium | Medium |
| [D02](D02_go_defer_devirt/) | go/defer skip devirt entirely (#52072) | KNOWN | ✅ | Low-Med | Medium |
| [D03](D03_generics_shape_devirt/) | Generic shape types block devirt | TODO | ✅ | Medium | Hard |
| [D04](D04_addrtaken_devirt/) | Address-taken interfaces block devirt | NEW | ✅ | Low-Med | Hard |
| [SSA01](SSA01_double_ext/) | Double zero/sign extension collapse | NEW | ✅* | Low | Easy |
| [SSA02](SSA02_consensus_identity/) | Boolean consensus identity not recognized | NEW | ❓ | Low | Easy |
| [SSA03](SSA03_distributive_factoring/) | Distributive AND/OR factoring | NEW | ❓ | Low | Easy |
| [SSA05](SSA05_redundant_and_after_zext/) | Redundant AND after zero-extend | NEW | ✅ | Low | Easy |

\* SSA01 already handled in Go 1.26.1 — may be historical.

---

## Correctness

> The compiler produces wrong results, or relies on a fragile workaround that could break.

| ID | Summary | Origin | Severity |
|----|---------|--------|----------|
| [F13](F13_scheduling_workaround/) | Prove pass manually reorders SSA values to work around scheduling bug (#76060) | KNOWN | Fragile workaround — correct today but brittle |

---

## Code Quality

> Compiler speed or maintainability issues. No impact on generated code.

| ID | Summary | Origin |
|----|---------|--------|
| [F14](F14_getbranch_quadratic/) | `getBranch` has quadratic behavior for jump tables | TODO |

---

## By Sub-Area

For navigating by compiler component:

| Sub-area | Findings | Key files in `go-src/` |
|----------|----------|------------------------|
| **BCE / Prove** | F02, F03, F04, F05, F11, F12 | `ssa/prove.go`, `ssa/loopbce.go` |
| **SSA Rules** | F01, SSA01-03, SSA05 | `ssa/generic.rules`, `ssa/AMD64.rules` |
| **Devirtualization** | D01-D05 | `devirtualize/devirtualize.go`, `devirtualize/pgo.go` |
| **Inlining** | D06 | `inline/inl.go` |
| **Escape Analysis** | F06 | `escape/*.go` |
| **Dead Code/Stores** | F07 | `ssa/deadstore.go` |
| **Control Flow** | F08, F09 | `ssa/shortcircuit.go`, `ssa/phiopt.go` |
| **Loop Optimization** | F10 | `ssa/loopbce.go` |
| **Scheduling** | F13 | `ssa/prove.go` (workaround) |
| **Compiler Speed** | F14 | `ssa/prove.go` |

---

## Full Classification Table

| ID | Summary | Category | Origin | Status | Impact | Difficulty |
|----|---------|----------|--------|--------|--------|------------|
| F01 | Boolean absorption rules missing | perf | NEW | ✅ | Medium | Easy |
| F02 | Modulo + fixup BCE miss | perf | TODO | ✅ | High | Hard |
| F03 | `i+1 < len(a)` fence-post BCE miss | perf | NEW | ✅ | High | Medium |
| F04 | Unsigned arithmetic BCE miss | perf | NEW | ✅ | Medium | Medium |
| F05 | Multiplication constraint BCE miss | perf | NEW | ✅ | Medium | Hard |
| F06 | Local collection false-escape | perf | KNOWN-BAD | ✅ | High | Hard |
| F07 | Basic-block-only DSE | perf | TODO | ✅ | Medium | Hard |
| F08 | Shortcircuit >2 preds | perf | TODO | ✅ | Med-High | Medium |
| F09 | phiopt >2 preds | perf | TODO | ✅ | Low-Med | Medium |
| F10 | Unsigned induction vars | perf | TODO | ✅ | Medium | Medium |
| F11 | OR facts commented out | perf | KNOWN | ✅ | Medium | Hard |
| F12 | Transitive equality | perf | TODO | ❓ | Low | Medium |
| F13 | Scheduling workaround | correctness | KNOWN | ✅ | Low | Medium |
| F14 | getBranch quadratic | compiler-speed | TODO | ❓ | Low | Easy |
| D01 | PPARAM devirt | perf | NEW | ✅ | Medium | Medium |
| D02 | go/defer devirt | perf | KNOWN | ✅ | Low-Med | Medium |
| D03 | Generics shape devirt | perf | TODO | ✅ | Medium | Hard |
| D04 | Addrtaken devirt | perf | NEW | ✅ | Low-Med | Hard |
| D05 | Struct-field devirt | perf | NEW | ✅ | Medium | Hard |
| D06 | OCALLINTER inline cost | perf | NEW | ✅ | Medium | Medium |
| SSA01 | Double extension | perf | NEW | ✅* | Low | Easy |
| SSA02 | Consensus identity | perf | NEW | ❓ | Low | Easy |
| SSA03 | Distributive factoring | perf | NEW | ❓ | Low | Easy |
| SSA05 | Redundant AND after zext | perf | NEW | ✅ | Low | Easy |

---

## Statistics

| Metric | Count |
|--------|-------|
| Total findings | 24 |
| Origin: NEW | 13 |
| Origin: TODO | 7 |
| Origin: KNOWN/KNOWN-BAD | 4 |
| Status: Confirmed (✅) | 20 |
| Category: Performance | 22 |
| Category: Correctness | 1 |
| Category: Compiler speed | 1 |
| Security-relevant watch areas | 3 |

## Recommended Priority

**Quick wins** (easy, upstreamable now):
1. **F01** — 8 SSA rewrite rules, ~30 min of work

**High-impact** (worth the effort):
2. **F03** — 10% benchmark hit, extremely common pattern
3. **F08** — Source TODO says "high impact"
4. **F10** — Extend `findIndVar` to unsigned

**Strategic** (improve whole subsystems):
5. **D01 + D05 + D06** — devirt/inline improvements compound
6. **F02** — unlocks modulo patterns in hash tables
