# API Proposals - NG

## 1. Extension of `ReadabilityOptions`

In the original JavaScript implementation of `tinycolor`, the `mostReadable` function accepts an options object that may contain an `includeFallbackColors` property (boolean). 

Since the `ReadabilityOptions` struct defined in the frozen API contract currently only specifies:

```go
type ReadabilityOptions struct {
    Level string // "AA" or "AAA"
    Size  string // "small" or "large"
}
```

**Proposal:** Add the `IncludeFallbackColors` boolean to the `ReadabilityOptions` struct so that `MostReadable` can properly support the fallback behaviors described in the JS implementation.

```go
type ReadabilityOptions struct {
    Level                 string // "AA" or "AAA"
    Size                  string // "small" or "large"
    IncludeFallbackColors bool
}
```

## 2. Shared `hexNames` Map
I need the `hexNames` map inside `names.go` to be exported or globally available within the package so that my implementation of `c.ToName()` can look up the named color string for a given hex value. I've left a basic loop using `hexNames` in `conversions.go` assuming this will be available in BB's files.
