# Issue Follow-up: Corrected Impact Analysis Using `jlog`

> Response to @Jorropo and @randall77's feedback on
> [#78632](https://github.com/golang/go/issues/78632)

## Correction

My earlier analysis using AND/OR instruction count deltas was flawed — it
conflated register allocation cascading effects with actual rule firings. Thanks
to @Jorropo's `jlog` technique and @randall77's `-genLog` pointer, here are the
real numbers.

## Method

Applied @Jorropo's patch (adding `&& jlog(v, "and absorb")` / `"or absorb"` to
the rules), rebuilt the compiler, then compiled each project with `go build -a
./...` and collected all `absorb` log lines.

## Standard Library

**9 hits** (previously claimed: 88)

| # | File | Function | Rule |
|---|------|----------|------|
| 1 | `runtime/cgocall.go:495` | `cgocallbackg1` | and absorb |
| 2 | `testing/testing.go:1818` | `(*common).runCleanup` | and absorb |
| 3 | `go/parser/parser.go:1763` | `(*parser).parsePrimaryExpr` | and absorb |
| 4 | `go/parser/parser.go:1768` | `(*parser).parsePrimaryExpr` | and absorb |
| 5 | `go/parser/parser.go:1774` | `(*parser).parsePrimaryExpr` | and absorb |
| 6 | `go/parser/parser.go:1782` | `(*parser).parsePrimaryExpr` | and absorb |
| 7 | `go/parser/parser.go:1887` | `(*parser).parseBinaryExpr` | and absorb |
| 8 | `x/net/http3/server.go:197` | `(*server).shutdown` | and absorb |
| 9 | `x/net/http3/server.go:201` | `(*server).shutdown` | and absorb |

The `go/parser` cluster (5 hits) comes from a type-switch with multiple cases
sharing a common `p.exprLev < 0` guard — inlining + control flow flattening
produces the absorption pattern at the SSA level.

## Hugo (github.com/gohugoio/hugo)

**2 project-specific hits** (previously claimed: 260)

| # | File | Function | Rule |
|---|------|----------|------|
| 1 | `publisher/publisher.go:132` | `DestinationPublisher.Publish` | and absorb |
| 2 | `commands/commandeer.go:426` | `(*rootCommand).Run` | and absorb |

Both are cold paths (file I/O, signal wait loop). Plus the 7 stdlib hits above.

## Prometheus (github.com/prometheus/prometheus)

**18 project-specific hits** (previously claimed: 70)

| # | File | Function | Rule |
|---|------|----------|------|
| 1 | `x/net/trace/events.go:102` | `RenderEvents` | and absorb |
| 2 | `x/net/trace/trace.go:286` | `Render` | and absorb |
| 3–5 | `go-openapi/swag/loading.go:151,155,158` | `loadHTTPBytes.func1` | and absorb |
| 6–9 | `microsoft-authentication-library/.../comm.go:247,259,265,272` | `(*Client).do` | and absorb |
| 10–13 | `x/oauth2/deviceauth.go:201,205,210,223` | `(*Config).DeviceAccessToken` | and absorb |
| 14–15 | `client_golang/promhttp/http.go:185,189` | `HandlerForTransactional.func1` | and absorb |
| 16 | `rules/group.go:546` | `(*Group).Eval.func1` | and absorb |
| 17–18 | `web/api/v1/api.go:544,563` | `(*API).query` | and absorb |

Most hits are in dependencies (Azure auth, OAuth2, OpenAPI), not core Prometheus.
Only 3 hits are in Prometheus's own code (`rules/group.go`, `web/api/v1/api.go`).

The `promhttp` hits come from a switch on error-handling strategy
(`ContinueOnError` / `HTTPErrorOnError`) — same pattern as `go/parser`: a
multi-branch control flow where inlining produces redundant AND/OR at SSA level.

## Lazygit (github.com/jesseduffield/lazygit)

**0 project-specific hits** (previously claimed: 27)

Only the 7 stdlib hits. Zero absorption firings in lazygit's own code or
dependencies.

## Summary

| Project | Previously claimed | Actual (total) | Actual (project-specific) |
|---------|-------------------:|---------------:|--------------------------:|
| stdlib | -88 | 9 | 9 |
| hugo | -260 | 9 | 2 |
| prometheus | -70 | 25 | 18 |
| lazygit | -27 | 7 | 0 |

**All hits are `and absorb`. Zero `or absorb` hits anywhere.**

The pattern is real but uncommon. It arises from SSA-level interactions
(inlining + control flow flattening producing redundant boolean ops), not from
source-level `x & (x | y)`. None of the hits are in hot loops or
performance-critical paths.

The rules are still correct and the code impact is non-zero, but the performance
case is much weaker than I originally presented. I apologize for the misleading
instruction-count methodology.

## Remaining argument for inclusion

- The rules are trivially correct (fundamental boolean identity), zero risk.
- They do fire in real code (25 unique locations across stdlib + prometheus).
- They eliminate 2 dead instructions per hit — small but free.
- Zero compile-time cost (confirmed with compilebench).
- Two lines of rules, auto-generated codegen — near-zero maintenance burden.

Whether that clears the bar for inclusion is a judgment call I defer to the
maintainers.
