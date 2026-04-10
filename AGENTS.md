# Go Compiler Bug Hunting & Performance Improvement

## Purpose

Research workspace for finding **bugs** and **performance improvements** in the Go compiler. The Go source lives in `go-src/` (cloned from `https://go.googlesource.com/go` at tip).

Goal: produce **concrete, upstreamable patches** backed by failing tests or benchmarks.

---

## Repository Layout

```
├── AGENTS.md                       ← you are here
├── go-src/                         ← Go source tree (upstream tip, read-only)
│   └── src/cmd/compile/internal/   ← the compiler
│       ├── ssa/                    ← SSA backend (prove, rewrite rules, regalloc, ...)
│       ├── inline/                 ← inlining heuristics
│       ├── escape/                 ← escape analysis
│       ├── devirtualize/           ← interface call devirtualization
│       └── ...
└── findings/                       ← all findings, one directory each
    ├── README.md                   ← master index with triage views and classification
    ├── F01_absorption_rules/       ← example finding
    │   ├── README.md               ← description, classification, location, fix
    │   └── reproduce.go            ← minimal reproducer
    └── ...
```

### Rules

- **Nothing goes in the repo root** except `AGENTS.md`, `.gitignore`, and `go-src/`.
- **All findings live under `findings/`**. One directory per finding.
- **No binaries, .o files, or build artifacts** committed.
- **No scattered test files** — every `.go` file belongs inside a finding directory.

---

## Finding Structure

Each finding gets a directory under `findings/` with a standardized name:

```
findings/<ID>_<short_name>/
├── README.md       ← required: description, metadata table, reproduction, location
├── reproduce.go    ← required if applicable: minimal Go program demonstrating the issue
├── bench_test.go   ← optional: benchmark proving performance impact
└── fix.patch       ← optional: proposed patch against go-src/
```

### Naming Convention

| Prefix | Area |
|--------|------|
| `F##`  | General compiler finding (BCE, DSE, control flow, etc.) |
| `D##`  | Devirtualization / inlining |
| `SSA##`| SSA rewrite rules |
| `E##`  | Escape analysis |
| `R##`  | Runtime |
| `L##`  | Linker |

---

## Classification System

### Required README.md Metadata Table

Every finding README.md **must** start with a metadata table:

```markdown
# <ID>: <One-line summary>

| Field | Value |
|-------|-------|
| **Category** | performance / correctness / security / compiler-speed / code-quality |
| **Sub-area** | bce / escape / devirt / inline / ssa-rules / nilcheck / deadcode / ... |
| **Origin** | NEW / TODO / KNOWN / KNOWN-BAD |
| **Status** | CONFIRMED / HYPOTHETICAL |
| **Difficulty** | Easy / Medium / Hard |
| **Impact** | Low / Medium / High |
| **Security** | vuln / latent-risk / watch-area / none |
| **Tested on** | go1.X.Y / tip at <commit> |
| **Lifecycle** | open / in-progress / submitted / upstream-rejected / fixed-in-<version> |
| **Related issues** | #12345 (if any — both source-referenced AND issue-tracker matches) |
```

### Category — "What kind of problem is this?"

| Category | What it means | Examples |
|----------|---------------|----------|
| **performance** | Compiler generates correct but slower code than it could | Missed BCE, false escape, redundant instructions |
| **correctness** | Compiler generates wrong code or has fragile workaround | Miscompilation, incorrect optimization |
| **security** | Bug causes memory unsafety in compiled programs | Wrong bounds-check removal, incorrect escape to stack |
| **compiler-speed** | The compiler itself is slower than it needs to be | Quadratic algorithms, redundant passes |
| **code-quality** | Source quality of the compiler | Dead code, misleading comments |

### Security — "Could this hurt users?"

Compiler bugs can have security implications. **Every finding must be evaluated**:

```
                              ┌─ Does the compiler REMOVE a safety check it shouldn't?
                              │  (bounds check, nil check, overflow check)
                              │  → SECURITY BUG (vuln) — file immediately
                              │
  Is this in a safety-        ├─ Does escape analysis say "stack" when value escapes?
  critical code path?  ──YES──┤  → SECURITY BUG (vuln) — use-after-free risk
  (prove.go, nilcheck.go,    │
   escape/)                   ├─ Is there a fragile workaround that could regress?
                              │  → latent-risk — correct today, breakable tomorrow
                              │
                              └─ Is the code conservative? (keeps checks it could remove)
                                 → watch-area — safe, but changes here need extra care

  Is this a missed      ──YES── none — no security relevance
  optimization only?
  (rewrite rules, devirt,
   inlining, dead stores)
```

| Security label | Icon | Meaning | Action |
|----------------|------|---------|--------|
| **vuln** | 🔴 | Active security bug | File upstream immediately, request CVE |
| **latent-risk** | ⚠️ | Fragile correctness in safety path | Document, test heavily, flag in reviews |
| **watch-area** | 🔍 | Conservative safety code we want to change | Any PR touching this needs security review |
| **none** | ⚪ | No security relevance | Normal review process |

> **Key insight**: A *kept* bounds check (not removed when it could be) = **performance** issue.
> A *wrong* bounds check removal (removed when it shouldn't be) = **security** issue.
> The prove pass does both — so any change to it must be evaluated for both sides.

### Origin — "Is this new or already known?"

| Origin | Meaning | How to identify |
|--------|---------|-----------------|
| **NEW** | We discovered this independently | Not in source comments AND not in issue tracker |
| **TODO** | Derived from an existing TODO/FIXME in source | The source code explicitly mentions it |
| **KNOWN** | Has an existing upstream Go issue | Found via issue tracker search or source reference |
| **KNOWN-BAD** | Known limitation in Go's test suite | Marked `// BAD` or `// known limitation` in tests |

**To verify origin, you MUST do both:**
1. `grep -rn 'TODO\|FIXME\|issue\|go.dev' go-src/src/cmd/compile/internal/ssa/<file>`
2. Search the upstream tracker: `https://github.com/golang/go/issues?q=is:issue+<keywords>`

A finding is only NEW if both come up empty.

### Status — "Is this real?"

| Status | Meaning | What you need |
|--------|---------|---------------|
| **CONFIRMED** | Proven with actual compiler output | Paste the `go tool compile` output showing the issue |
| **HYPOTHETICAL** | Code review suggests it, but no reproducer yet | Flag it — don't claim it's real until verified |

"I read the code and it looks buggy" is **HYPOTHETICAL**, not CONFIRMED.

### Impact — "How much does this matter?"

| Impact | Criteria |
|--------|----------|
| **High** | Measured benchmark regression (>5%), or affects hot paths in common programs (hash tables, HTTP handlers, slice loops) |
| **Medium** | Theoretical perf impact from redundant instructions/allocations, affects moderately common patterns |
| **Low** | Rare patterns, minor instruction count difference, or compiler-internal only |

### Lifecycle — "What's happening with this?"

| Lifecycle | Meaning |
|-----------|---------|
| **open** | Finding recorded, no work started on a fix |
| **in-progress** | Someone is actively working on a fix |
| **submitted** | CL/patch submitted upstream (link it in Related issues) |
| **upstream-rejected** | Upstream declined the fix (document why in README) |
| **fixed-in-\<version\>** | Fixed in a Go release (e.g. `fixed-in-1.24`) |

---

### reproduce.go Requirements

- Use `package p` (not `package main`) so `go tool compile` works without linking.
- Include the exact command to reproduce as a comment at the top.
- Include both the **expected** and **actual** output.
- Keep it **minimal** — one function demonstrating one issue.
- No `fmt` or external imports if possible.
- Use `//go:noinline` where needed to prevent the issue from being optimized away.
- Note the Go version: `// Tested: go1.26.1 linux/amd64`

---

## Workflow

### 1. Pick a Target Area

| Area | Key files | Diagnostic flags |
|------|-----------|------------------|
| BCE / Prove | `ssa/prove.go`, `ssa/loopbce.go` | `-d=ssa/check_bce`, `-d=ssa/prove/debug=2` |
| SSA Rules | `ssa/generic.rules`, `ssa/AMD64.rules` | `-S` (assembly output) |
| Escape | `escape/*.go` | `-m`, `-m -m` |
| Devirt | `devirtualize/*.go` | `-m`, `-m -m` |
| Inlining | `inline/*.go` | `-m` (inline decisions) |
| Nilcheck | `ssa/nilcheck.go` | `-S` (look for TESTB/JEQ patterns) |
| Dead stores | `ssa/deadstore.go` | `GOSSAFUNC=X` for SSA HTML dump |

### 2. Read Before Writing

1. **Read the source file** top-to-bottom. Use `file_read` with ranges for large files.
2. **Read the tests** — `go-src/test/prove.go`, `go-src/test/escape*.go`, etc.
3. **Grep for TODOs**: `grep -rn 'TODO\|FIXME\|HACK\|BUG' <path>`
4. **Check git history**: `cd go-src && git log --oneline -20 -- <path>`

### 3. Investigate & Reproduce

1. Write a minimal `reproduce.go` in a **new finding directory** (not the repo root!).
2. **Test the reproducer with the current compiler FIRST** — before reading source to explain it:
   ```bash
   go tool compile -d=ssa/check_bce findings/<ID>/reproduce.go   # BCE
   go tool compile -S findings/<ID>/reproduce.go 2>&1 | grep ...  # SSA rules
   go tool compile -m findings/<ID>/reproduce.go                   # escape/devirt
   ```
3. If the issue **doesn't reproduce** on the current compiler, it may be already fixed. Check `git log` for the relevant file. If fixed, record as `fixed-in-<version>` and move on.
4. Record the **actual** vs **expected** output in the README.
5. If performance-related, write a `bench_test.go` with before/after.

### 4. Classify

Before recording, answer **all** of these:

1. **Category**: Performance, correctness, or security?
2. **Security**: Walk the decision tree above. If it touches prove.go/nilcheck.go/escape, explain why it's safe (or not).
3. **Origin** (two checks required):
   - Source code: `grep -r 'TODO\|FIXME\|issue' go-src/src/cmd/compile/internal/ssa/<file>`
   - Issue tracker: search `https://github.com/golang/go/issues?q=<keywords>`
   - Only mark NEW if **both** are empty.
4. **Status**: Do you have `go tool compile` output proving it? → CONFIRMED. Otherwise → HYPOTHETICAL.
5. **Impact**: Do you have a benchmark? If claiming High, you **must** have numbers.
6. **Tested on**: Record exact `go version` output.

### 5. Record

```bash
mkdir -p findings/<ID>_<name>
# Write README.md with FULL metadata table (all fields including Security, Tested on, Lifecycle)
# Write reproduce.go with minimal reproducer
# Optionally: bench_test.go, fix.patch
```

Then update `findings/README.md` — add one row to the appropriate triage section AND the full table.

### 6. Verify (Checklist)

Before considering a finding complete:

- [ ] `reproduce.go` compiles: `go tool compile findings/<ID>/reproduce.go`
- [ ] Issue confirmed with diagnostic output (pasted in README)
- [ ] Origin checked against **both** source code and issue tracker
- [ ] Security evaluated using the decision tree
- [ ] `findings/README.md` updated with new row in correct section
- [ ] No stray files outside the finding directory (no `.o`, no binaries, no loose `.go`)

### 7. Fix (optional — when ready to submit upstream)

1. Create a worktree: `cd go-src && git worktree add ../fix-<ID> master`
2. Write the fix in the worktree.
3. Write/update tests in the worktree.
4. Run the relevant test suite: `go test cmd/compile/internal/ssa`
5. Run the broader suite: `cd go-src/src && ./run.bash`
6. Save the patch: `cd fix-<ID> && git diff > ../findings/<ID>_<name>/fix.patch`
7. Update the finding's lifecycle to `in-progress` or `submitted`.

---

## Handling Side Findings

While investigating area X, you will often spot issues in area Y. **Do not mix them.**

1. Create a **new** finding directory immediately: `mkdir -p findings/<ID>_<name>`
2. Write a **quick** README.md with at least the metadata table (mark Status as HYPOTHETICAL if not yet verified).
3. Optionally drop a draft `reproduce.go` in there.
4. **Go back to your original investigation.** Don't context-switch.
5. Come back to the side finding later.

---

## Parallel Hunting

When running multiple sub-agents to cover different areas simultaneously:

1. **Assign each agent a single area** (BCE, devirt, escape, etc.) — no overlap.
2. **Each agent must write directly into `findings/<ID>_<name>/`** — never into the repo root or temporary directories.
3. **The coordinating agent** reviews results, deduplicates, and runs the verify checklist.
4. **Deduplication**: If two agents found the same underlying issue from different angles, merge into one finding and note both perspectives.

---

## Useful Commands

```bash
# Build compiler from source
cd go-src/src && ./make.bash

# Run specific package tests
cd go-src && go test cmd/compile/internal/ssa
cd go-src && go test cmd/compile/...

# Dump SSA for a function
GOSSAFUNC=FuncName go tool compile reproduce.go   # writes ssa.html

# Search rewrite rules
grep -n 'pattern' go-src/src/cmd/compile/internal/ssa/generic.rules

# Recent changes to a file
cd go-src && git log --oneline -20 -- src/cmd/compile/internal/ssa/prove.go

# Search upstream issue tracker
# https://github.com/golang/go/issues?q=is:issue+bounds+check+elimination

# Check current Go version
go version
```

---

## Rules

1. **Read before you write.** The Go compiler is large and subtle.
2. **One finding per directory.** Keep changes isolated.
3. **Every finding needs a README.md** with the full metadata table (all fields).
4. **Evaluate security for every finding.** Use the decision tree above.
5. **Verify origin against both source and issue tracker.** Only NEW if both empty.
6. **Test reproducers against the current compiler** before recording. Don't record stale findings.
7. **Reproducers must be minimal** — one function, one issue, `package p`.
8. **Benchmark performance claims.** "Faster" without numbers is not a finding. High impact requires measured numbers.
9. **No files outside `findings/`** — no root-level test files, no temp directories.
10. **No binaries or build artifacts** — `.gitignore` handles this.
11. **Never modify `go-src/` directly** — use worktrees for fixes.
12. **Never commit without explicit user instruction.**
13. **Update `findings/README.md`** when adding or modifying any finding.
14. **Record Go version** in every finding and reproducer.
