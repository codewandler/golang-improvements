# Go Compiler Bug Hunting & Performance Improvements

Research workspace for finding **bugs** and **performance improvements** in the [Go compiler](https://github.com/golang/go). The goal is to produce concrete, upstreamable patches backed by failing tests or benchmarks.

## Findings

24 findings across six areas of the Go compiler:

| Area | Findings | Description |
|------|----------|-------------|
| **BCE / Prove** | F02–F05, F10–F12 | Bounds check elimination misses |
| **SSA Rewrite Rules** | F01, SSA01–SSA03, SSA05 | Missing algebraic simplifications |
| **Devirtualization** | D01–D06 | Interface call devirtualization gaps |
| **Escape Analysis** | F06 | False escapes to heap |
| **Control Flow** | F08, F09 | Short-circuit and phi optimization limits |
| **Other** | F07, F13, F14 | Dead stores, scheduling, compiler speed |

See [`findings/README.md`](findings/README.md) for the full classified index with triage views.

### Upstream Submissions

| Finding | Issue | Status |
|---------|-------|--------|
| [F01: Boolean Absorption Laws](findings/F01_absorption_rules/) | [golang/go#78632](https://github.com/golang/go/issues/78632) | Submitted |

## Repository Layout

```
├── findings/               ← all findings, one directory each
│   ├── README.md           ← master index with triage views
│   ├── F01_absorption_rules/
│   │   ├── README.md       ← description, classification, proof
│   │   ├── reproduce.go    ← minimal reproducer
│   │   ├── fix.patch       ← human-authored patch
│   │   ├── ISSUE.md        ← upstream issue template
│   │   └── bench/          ← benchmark data
│   └── ...
├── go-src/                 ← Go source tree (submodule)
└── AGENTS.md               ← methodology and conventions
```

## How It Works

Each finding follows a structured process:

1. **Read** the compiler source — understand before proposing
2. **Reproduce** — minimal `package p` program + `go tool compile` output
3. **Classify** — category, origin, security impact, difficulty
4. **Prove** — assembly diff, benchmarks, compilebench
5. **Submit** — file upstream issue, push fix branch, open CL

Findings are classified by **security relevance** (could a bug here cause memory unsafety?), **performance impact** (measured with benchstat), and **origin** (independently discovered vs existing TODO/known issue).

See [`AGENTS.md`](AGENTS.md) for the full methodology.

## Go Source

The `go-src/` submodule tracks upstream tip from `https://go.googlesource.com/go`. Fix branches are pushed to the [codewandler/go](https://github.com/codewandler/go) fork.

## License

The findings, reproducers, and analysis in this repository are original work. The Go compiler source in `go-src/` is licensed under the [Go BSD license](https://go.dev/LICENSE).
