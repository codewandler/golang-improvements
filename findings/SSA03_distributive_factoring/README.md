# SSA03: Distributive Factoring Not Applied

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | SSA rewrite rules |
| **Origin** | NEW |
| **Status** | HYPOTHETICAL — needs verification |
| **Difficulty** | Easy |
| **Impact** | Low |

## Problem

`(x & y) | (x & z)` should be strength-reduced to `x & (y | z)` (saves one AND).
