# API PROPOSALS: BB

## Proposal 1: Expose HslToRgb and HsvToRgb Conversion Functions
**Target Module:** NG (Conversions)
**Rationale:**
The `parser.go` module needs to convert `hsl(...)` and `hsv(...)` string inputs into absolute RGB values during the parsing phase to correctly populate the internal `Color` state (as per `inputToRGB` logic). However, `HslToRgb` and `HsvToRgb` were not formally exposed in the frozen API contract for NG.

**Requested Frozen API Additions for NG:**
```go
// In conversions.go (NG)
func HslToRgb(h, s, l float64) RGB
func HsvToRgb(h, s, v float64) RGB
func RgbToRgb(r, g, b float64) RGB
```

**Temporary Workaround:**
To adhere to the zero-conflict rule and continue parallel development, I have implemented temporary private stubs `stubHslToRgb`, `stubHsvToRgb`, and `stubRgbToRgb` inside `parser.go`. These should be replaced by NG's final exported functions during the final integration phase.
