<!-- 
  Issue title: cmd/compile: add boolean absorption laws to SSA rewrite rules
-->

### Go version

```
go version go1.26.1 linux/amd64
```

Also confirmed on `go1.27-devel_4478774aa2` (tip).

### Output of `go env` in your module/workspace

<details>

```
AR='ar'
CC='gcc'
CGO_CFLAGS='-O2 -g'
CGO_CPPFLAGS=''
CGO_CXXFLAGS='-O2 -g'
CGO_ENABLED='1'
CGO_FFLAGS='-O2 -g'
CGO_LDFLAGS='-O2 -g'
CXX='g++'
GCCGO='gccgo'
GO111MODULE='auto'
GOAMD64='v1'
GOARCH='amd64'
GOAUTH='netrc'
GOBIN=''
GOCACHE='/home/timo/.cache/go-build'
GODEBUG=''
GOENV='/home/timo/.config/go/env'
GOEXE=''
GOEXPERIMENT=''
GOFIPS140='off'
GOFLAGS=''
GOGCCFLAGS='-fPIC -m64 -pthread -Wl,--no-gc-sections -fmessage-length=0 -ffile-prefix-map=/tmp/go-build2410965085=/tmp/go-build -gno-record-gcc-switches'
GOHOSTARCH='amd64'
GOHOSTOS='linux'
GOINSECURE=''
GOMOD=''
GOMODCACHE='/home/timo/go/pkg/mod'
GOOS='linux'
GOPATH='/home/timo/go'
GOROOT='/usr/lib/go'
GOSUMDB='sum.golang.org'
GOTELEMETRY='local'
GOTMPDIR=''
GOTOOLCHAIN='auto'
GOTOOLDIR='/usr/lib/go/pkg/tool/linux_amd64'
GOVCS=''
GOVERSION='go1.26.1'
GOWORK=''
PKG_CONFIG='pkg-config'
```

</details>

### What did you do?

The SSA generic rewrite rules implement DeMorgan's laws (`generic.rules:208-209`) but are missing the closely related **boolean absorption laws**:

- `x & (x | y) == x`
- `x | (x & y) == x`

These are fundamental boolean algebra identities ([absorption law](https://en.wikipedia.org/wiki/Absorption_law)) that hold for all bit patterns, all widths, signed and unsigned. Both GCC and LLVM recognize and optimize these patterns at `-O2`.

Minimal reproducer:

```go
package p

//go:noinline
func andAbsorption(x, y uint64) uint64 {
	return x & (x | y) // should simplify to: return x
}

//go:noinline
func orAbsorption(x, y uint64) uint64 {
	return x | (x & y) // should simplify to: return x
}
```

### What did you expect to see?

```asm
andAbsorption:
    MOVQ AX, AX   ; (or nothing — x already in return register)
    RET
```

The compiler recognizes `x & (x | y) == x` and eliminates the dead AND/OR operations.

GCC `-O2` on equivalent C produces just `movq %rdi, %rax; ret`.

### What did you see instead?

```
$ go tool compile -S reproduce.go 2>&1 | grep -E 'ORQ|ANDQ|RET'
```

```asm
andAbsorption:
    ORQ  AX, BX         ; compute x | y
    ANDQ BX, AX         ; compute x & (x | y)  — always equals x
    RET                  ; 3 instructions, 7 bytes

orAbsorption:
    ANDQ AX, BX         ; compute x & y
    ORQ  BX, AX         ; compute x | (x & y)  — always equals x
    RET                  ; 3 instructions, 7 bytes
```

Two redundant ALU instructions per occurrence. Same issue affects 32/16/8-bit variants (`ORL+ANDL`, `ORW+ANDW`, `ORB+ANDB`) and all argument orderings (e.g. `(x | y) & x`).

### Impact

In a tight loop operating on slices (1024 elements), eliminating the redundant instructions yields a significant improvement:

```
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i9-10900K CPU @ 3.70GHz
                  │ bench_before.txt │       bench_after.txt        │
                  │      sec/op      │   sec/op     vs base         │
LoopAbsorb64-20       441.2n ± 1%     226.9n ± 1%  -48.57% (p=0.002 n=6)
LoopIdentity64-20     228.6n ± 3%     235.2n ± 3%       ~ (p=0.065 n=6)
```

After the fix, `LoopAbsorb64` matches the identity baseline (just `a[i] = a[i]`), confirming the redundant instructions are fully eliminated.

No compile-time impact:

```
compilebench -count 6 -pkg math/big
                  │ before      │ after          │
                  │   sec/op    │  sec/op  vs base         │
Pkg                 147.7m ± 4%  148.0m ± 13%  ~ (p=0.699 n=6)
```

### Proposed fix

Add two rules to `src/cmd/compile/internal/ssa/_gen/generic.rules` (next to DeMorgan's laws at line 208):

```
// Absorption laws
(And(8|16|32|64) x (Or(8|16|32|64) x y)) => x
(Or(8|16|32|64)  x (And(8|16|32|64) x y)) => x
```

Commutativity of AND/OR is handled automatically by the rule engine — all argument orderings are matched without additional rules.

**Why this is safe**: By the [absorption law](https://en.wikipedia.org/wiki/Absorption_law), `x & (x | y)` distributes to `(x & x) | (x & y) = x | (x & y) = x`. This holds for all bit patterns, all widths, signed or unsigned. There are no edge cases. The result (`x`) is strictly simpler than the input, so the rules cannot create infinite rewriting loops.

I have a working patch with codegen tests covering 64-bit, 32-bit, and commuted argument orderings.

Fix branch: https://github.com/codewandler/go/tree/fix/F01-absorption-rules
