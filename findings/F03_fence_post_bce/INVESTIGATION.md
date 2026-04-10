# F03 Investigation Log

## Summary

`i+1 < len(a)` does not prove `a[i]` safe — both bounds checks remain.
The equivalent `i >= 1 && i < len(a)` with `a[i-1] + a[i]` has zero bounds checks.

## Confirmed (go1.26.1 linux/amd64)

```
$ go tool compile -d=ssa/check_bce reproduce.go
reproduce.go:9:11: Found IsInBounds
reproduce.go:9:18: Found IsInBounds
```

Both `a[i]` and `a[i+1]` retain bounds checks despite `i >= 0 && i+1 < len(a)`.

## Root Cause Analysis

The fence-post section in `prove.go:1097` handles these derivations:

| Case | Existing rule | Comment |
|------|--------------|---------|
| `gt` | `x+1 > w ⇒ x >= w` | ✅ handles delta on LHS |
| `gt` | `v > x-1 ⇒ v >= x` | ✅ handles delta on RHS (negative) |
| `gt\|eq` | `x-1 >= w && x > min ⇒ x > w` | ✅ |
| `gt\|eq` | `v >= x+1 && x < max ⇒ v > x` | ✅ |
| **`gt`** | **`v > x+1 && x < max ⇒ v > x`** | **❌ MISSING** |

The missing case is: when the RHS has a +1 delta under strict `gt`, derive the
strict inequality without the delta (with an overflow guard).

## Fix Attempted

Added to `case gt:` in the fence-post switch, mirroring the existing `gt|eq` pattern:

```go
if x, delta := isConstDelta(w); x != nil && delta == 1 {
    // v > x+1 && x < max  ⇒  v > x
    lim := ft.limits[x.ID]
    if (d == signed && lim.max < opMax[w.Op]) || (d == unsigned && lim.umax < opUMax[w.Op]) {
        ft.update(parent, v, x, d, gt)
    }
}
```

## Why the Fix Doesn't Fire

The overflow guard `lim.max < opMax[w.Op]` fails because **`i.max` is still
`MaxInt64` when the fence-post code runs**.

Debug trace (with instrumentation):

```
fence-post gt: v=v12 > x+1=v15 (x=v8 delta=1) d=signed
    lim.max=9223372036854775807  opMax=9223372036854775807   → FAILS (not <)
```

### Why `i.max` isn't tightened yet

The SSA for `pairAccess` processes facts in dominator order:

1. **b1→b2**: `Leq64(0, i)` → sets `i.min = 0` (signed). `i.max` stays `MaxInt64`.
2. **b2→b5**: `Less64(i+1, len(a))` → fence-post runs HERE.

At step 2, nothing has tightened `i.max` below `MaxInt64`. The comparison
`i+1 < len(a)` is `Less64` (signed domain), and the overflow guard checks
`i.max < MaxInt64` — which is `MaxInt64 < MaxInt64 = false`.

The limit *does* get tightened later (to `2305843009213693949` = `MaxSliceLen-2`)
by limit propagation *after* the fence-post section runs. On the second encounter
(at b5→b6, processing `IsInBounds`), the limits are tight enough:

```
fence-post gt: v=v12 > x+1=v15 (x=v8 delta=1) d=signed
    lim.max=2305843009213693949  opMax=9223372036854775807   → PASSES
fence-post gt: ... d=unsigned
    lim.umax=2305843009213693949 opUMax=18446744073709551615 → PASSES
```

But by that point, `v > x` is derived in the signed/unsigned domain — however the
`IsInBounds` at b5 has already been evaluated and the prove pass doesn't revisit it.

## Possible Next Steps

1. **Reorder within `update`**: Move fence-post processing AFTER limit propagation
   (lines 930–990) so limits are tightened before the overflow guard is checked.
   Risk: recursive `update` calls from fence-post could interact with partially
   updated limits.

2. **Relax the overflow guard for slice indices**: When `x` is known non-negative
   (`lim.min >= 0`) and `w.Op` is `OpAdd64`, then `x+1` cannot overflow signed
   (since `x <= MaxInt64` and `x >= 0` means `x+1 <= MaxInt64+1` which overflows,
   but practically slice lengths are bounded by `MaxSliceLen`). Could check
   `lim.min >= 0 && lim.max <= MaxSliceLen` instead.

3. **Two-pass approach**: Run fence-post implications a second time after all
   limit propagation is complete.

4. **Direct transitivity**: Instead of fence-post, teach the prove pass that
   `i < i+1` (when `i >= 0`) and combine with `i+1 < len(a)` via the poset's
   transitive closure.

Option 1 seems most promising — investigate whether the limit propagation at
lines 930–990 can safely precede the fence-post section.
