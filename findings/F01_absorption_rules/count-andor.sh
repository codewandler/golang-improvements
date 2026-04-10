#!/usr/bin/env bash
# count-andor.sh — Compare AND/OR instruction counts between two Go compilers.
#
# Usage:
#   ./count-andor.sh <pkg1> <pkg2> ...
#   ./count-andor.sh runtime sync math/bits math/big fmt strings bytes
#
# For external repos (not stdlib), run from within the repo:
#   cd /path/to/repo && ../count-andor.sh ./...
#
# Environment:
#   PATCHED_GOROOT  — GOROOT for the patched compiler (required)
#
# How it works:
#   For each package, compiles with `go build -gcflags="-S"` using both
#   the stock and patched compilers.  Counts all AND/OR instructions
#   (ORQ, ORL, ORW, ORB, ANDQ, ANDL, ANDW, ANDB) in the generated
#   assembly.  Reports per-package and total deltas.

set -euo pipefail

PATCHED_GOROOT="${PATCHED_GOROOT:?Set PATCHED_GOROOT to the patched compiler's GOROOT}"
PATCHED_GO="$PATCHED_GOROOT/bin/go"

PATTERN='\b(ORQ|ORL|ORW|ORB|ANDQ|ANDL|ANDW|ANDB)\b'

count_andor() {
    local pkg="$1" goroot="${2:-}" go_bin="${3:-go}"
    local tmpfile
    tmpfile=$(mktemp)

    if [ -n "$goroot" ]; then
        GOROOT="$goroot" "$go_bin" build -gcflags="-S" "$pkg" >"$tmpfile" 2>&1 || true
    else
        go build -gcflags="-S" "$pkg" >"$tmpfile" 2>&1 || true
    fi

    grep -cE "$PATTERN" "$tmpfile" 2>/dev/null || echo 0
    rm -f "$tmpfile"
}

printf "%-30s %8s %8s %8s\n" "Package" "Stock" "Patched" "Delta"
printf "%-30s %8s %8s %8s\n" "-------" "-----" "-------" "-----"

total_stock=0
total_patched=0

for pkg in "$@"; do
    n1=$(count_andor "$pkg")
    n2=$(count_andor "$pkg" "$PATCHED_GOROOT" "$PATCHED_GO")

    delta=$((n2 - n1))
    total_stock=$((total_stock + n1))
    total_patched=$((total_patched + n2))

    marker=""
    [ "$delta" -lt 0 ] && marker="  <---"
    printf "%-30s %8d %8d %8d%s\n" "$pkg" "$n1" "$n2" "$delta" "$marker"
done

echo ""
printf "%-30s %8d %8d %8d\n" "TOTAL" "$total_stock" "$total_patched" "$((total_patched - total_stock))"
