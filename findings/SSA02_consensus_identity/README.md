# SSA02: Boolean Consensus Identity Not Recognized

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | SSA rewrite rules |
| **Origin** | NEW |
| **Status** | HYPOTHETICAL — needs re-verification |
| **Difficulty** | Easy |
| **Impact** | Low |

## Problem

`(x | y) & (x | ^y) => x` — the boolean consensus/resolution identity is not in generic.rules.
