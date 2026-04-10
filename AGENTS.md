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
│       ├── ssa/                    ← SSA backend (prove, rewrite rules, regalloc, …)
│       ├── inline/                 ← inlining heuristics
│       ├── escape/                 ← escape analysis
│       ├── devirtualize/           ← interface call devirtualization
│       └── …
└── findings/                       ← all findings, one directory each
    ├── README.md                   ← master index with classification table
    ├── F01_absorption_rules/       ← example finding
    │   ├── README.md               ← description, classification, location, fix
    │   └── reproduce.go            ← minimal reproducer
    ├── F02_modulo_fixup_bce/
    │   ├── README.md
    │   └── reproduce.go
    └── …
```

### Rules

- **Nothing goes in the repo root** except `AGENTS.md` and `go-src/`.
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
└── fix.patch       ← optional: proposed patch
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

### Required README.md Metadata Table

Every finding README.md **must** start with a metadata table:

```markdown
# <ID>: <One-line summary>

| Field | Value |
|-------|-------|
| **Category** | performance / correctness / compiler-speed / code-quality |
| **Sub-area** | bce / escape / devirt / inline / ssa-rules / nilcheck / deadcode / … |
| **Origin** | NEW / TODO / KNOWN / KNOWN-BAD |
| **Status** | ✅ CONFIRMED / ❓ HYPOTHETICAL |
| **Difficulty** | Easy / Medium / Hard |
| **Impact** | Low / Medium / High |
| **Related issues** | #12345 (if any) |
```

### Origin Classification

| Origin | Meaning | How to identify |
|--------|---------|-----------------|
| **NEW** | We discovered this independently | No TODO/FIXME or issue reference in source |
| **TODO** | Derived from an existing TODO/FIXME in source | The source code explicitly mentions it |
| **KNOWN** | Has an existing upstream Go issue | Referenced by `go.dev/issue/N` or `#N` in source |
| **KNOWN-BAD** | Known limitation in Go's test suite | Marked `// BAD` or `// known limitation` in tests |

### reproduce.go Requirements

- Use `package p` (not `package main`) so `go tool compile` works without linking.
- Include the exact command to reproduce as a comment at the top.
- Include both the expected and actual output.
- Keep it **minimal** — one function demonstrating one issue.
- No `fmt` or external imports if possible.
- Use `//go:noinline` where needed to prevent the issue from being optimized away.

---

## Workflow

### 1. Pick a Target Area

Start with one of these high-value targets:

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

### 3. Investigate

1. Write a minimal `reproduce.go` that should trigger the issue.
2. Compile with the appropriate diagnostic flag.
3. Record the **actual** vs **expected** output.
4. If performance-related, write a `bench_test.go` with before/after.

### 4. Classify

Before recording, determine:
- Is this NEW or does a TODO/FIXME already exist?
- Does an upstream issue already exist? (`grep -r 'issue' go-src/src/cmd/compile/internal/ssa/<file>`)
- Is it confirmed (have compiler output proving it) or hypothetical?

### 5. Record

```bash
mkdir -p findings/<ID>_<name>
# Write README.md with metadata table
# Write reproduce.go with minimal reproducer
# Optionally: bench_test.go, fix.patch
```

Then update `findings/README.md` — add one row to the appropriate table.

### 6. Verify

```bash
# BCE findings
go tool compile -d=ssa/check_bce findings/<ID>/reproduce.go

# SSA rule findings  
go tool compile -S findings/<ID>/reproduce.go 2>&1 | grep -A5 '<function_name>'

# Escape findings
go tool compile -m findings/<ID>/reproduce.go

# Devirt findings
go tool compile -m findings/<ID>/reproduce.go 2>&1 | grep devirt
```

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
```

---

## Rules

1. **Read before you write.** The Go compiler is large and subtle.
2. **One finding per directory.** Keep changes isolated.
3. **Every finding needs a README.md** with the metadata table.
4. **Reproducers must be minimal** — one function, one issue.
5. **Benchmark performance claims.** "Faster" without numbers is not a finding.
6. **No files in the repo root** — everything under `findings/`.
7. **No binaries or build artifacts** — add `*.o` and named binaries to `.gitignore`.
8. **Never modify `go-src/` directly** — use worktrees or branches.
9. **Never commit without explicit user instruction.**
10. **Update `findings/README.md`** when adding or modifying findings.
