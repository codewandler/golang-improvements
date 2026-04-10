# SSA05: Redundant AND After Zero-Extend

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | SSA rewrite rules |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED |
| **Difficulty** | Easy |
| **Impact** | Low |
| **Security** | ⚪ none |

## Problem

`uint64(x) & 0xFF` where `x` is `uint8` — the AND is redundant since zero-extend already guarantees the value fits in 8 bits.
