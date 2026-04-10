# SSA01: Double Zero/Sign Extension Not Collapsed

| Field | Value |
|-------|-------|
| **Category** | Performance |
| **Sub-area** | SSA rewrite rules |
| **Origin** | NEW |
| **Status** | ✅ CONFIRMED — but Go 1.26 may already handle this |
| **Difficulty** | Easy |
| **Impact** | Low |
| **Security** | ⚪ none |

## Problem

`ZeroExt16to64(ZeroExt8to16(x))` should collapse to `ZeroExt8to64(x)`.

Note: Testing showed this IS already optimized in Go 1.26.1 (single `MOVBLZX`). May have been a valid finding on older versions only.
