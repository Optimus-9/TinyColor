# TinyColor (Go Port)

TinyColor is a small, fast library for color manipulation and conversion in Go. It allows many forms of input, while providing color conversions and other color utility functions. It has no external dependencies.

This repository is a production-quality Go port of the original JavaScript library, submitted for the PORT_MORTEM Code Resurrection 2026 Hackathon.

## Original Project

The original JavaScript repository can be found here: [bgrins/TinyColor](https://github.com/bgrins/TinyColor).
The original JS test suite is preserved in `tests/original/` for verification.

## Build

To build the package:
```sh
go build ./...
```

## Testing

To run the native Go test suite (which includes 46 tests ported from the original JS suite):
```sh
go test -v ./...
```

## Repository Layout

- `types.go`: Public color structs (`RGB`, `HSL`, `HSV`) and options (`ReadabilityOptions`, `Format`).
- `core.go`: The central `Color` struct, `New()`, and basic state methods.
- `parser.go`: Highly permissive string and interface parsing logic.
- `conversions.go`: Core math for color space conversions.
- `modify.go`: Color manipulation methods (e.g., `Lighten()`, `Spin()`).
- `combinations.go`: Methods generating color palettes (e.g., `Analogous()`, `Triad()`).
- `format.go`: String formatting outputs.
- `wcag.go`: Accessibility and contrast ratio math.
- `names.go`: W3C CSS named colors mapping.

## Usage

Import the package in your Go code:

```go
import tc "tinycolor"
```

Call `tc.New(input)` to create a `Color` pointer. See "Accepted Input" below for more information about what is accepted.

## Accepted Input

The string parsing is very permissive. It is meant to make typing a color as input as easy as possible. All commas, percentages, and parenthesis are optional, and most inputs allow either `0-1`, `0%-100%`, or `0-n` (where `n` is either 100, 255, or 360 depending on the value).

HSL and HSV both require either `0%-100%` or `0-1` for the `S`/`L`/`V` properties. The `H` (hue) can have values between `0%-100%` or `0-360`.

RGB input requires either `0-255` or `0%-100%`.

If you call `tc.FromRatio`, RGB and Hue input can also accept `0-1`.

Here are some examples of string input:

### Hex, 8-digit (RGBA) Hex
```go
tc.New("#000")
tc.New("000")
tc.New("#369C")
tc.New("369C")
tc.New("#f0f0f6")
tc.New("f0f0f6")
tc.New("#f0f0f688")
tc.New("f0f0f688")
```

### RGB, RGBA
```go
tc.New("rgb (255, 0, 0)")
tc.New("rgb 255 0 0")
tc.New("rgba (255, 0, 0, .5)")
tc.New(map[string]interface{}{"r": 255, "g": 0, "b": 0})
tc.New(tc.RGB{R: 255, G: 0, B: 0, A: 1})
tc.FromRatio(map[string]interface{}{"r": 1, "g": 0, "b": 0})
tc.FromRatio(map[string]interface{}{"r": 0.5, "g": 0.5, "b": 0.5})
```

### HSL, HSLA
```go
tc.New("hsl(0, 100%, 50%)")
tc.New("hsla(0, 100%, 50%, .5)")
tc.New("hsl 0 1.0 0.5")
tc.New(map[string]interface{}{"h": 0, "s": 1, "l": 0.5})
tc.New(tc.HSL{H: 0, S: 1, L: 0.5, A: 1})
tc.FromRatio(map[string]interface{}{"h": 1, "s": 0, "l": 0})
tc.FromRatio(map[string]interface{}{"h": 0.5, "s": 0.5, "l": 0.5})
```

### HSV, HSVA
```go
tc.New("hsv(0, 100%, 100%)")
tc.New("hsva(0, 100%, 100%, .5)")
tc.New("hsv (0 100% 100%)")
tc.New("hsv 0 1 1")
tc.New(map[string]interface{}{"h": 0, "s": 100, "v": 100})
tc.New(tc.HSV{H: 0, S: 100, V: 100, A: 1})
tc.FromRatio(map[string]interface{}{"h": 1, "s": 0, "v": 0})
tc.FromRatio(map[string]interface{}{"h": 0.5, "s": 0.5, "v": 0.5})
```

### Named Colors
Case insensitive names are accepted, using the list of colors in the CSS spec.
```go
tc.New("RED")
tc.New("blanchedalmond")
tc.New("darkblue")
```

## Methods

### GetFormat
Returns the format used to create the tinycolor instance.
```go
color := tc.New("red")
color.GetFormat() // tc.FormatName

color2 := tc.New(tc.RGB{R: 255, G: 255, B: 255})
color2.GetFormat() // tc.FormatRgb
```

### GetOriginalInput
Returns the input passed into the constructor used to create the tinycolor instance.
```go
color := tc.New("red")
color.GetOriginalInput() // "red"
```

### IsValid
Return a boolean indicating whether the color was successfully parsed. Note: if the color is not valid then it will act like `black` when being used with other methods.
```go
color1 := tc.New("red")
color1.IsValid() // true
color1.ToHexString() // "#ff0000"

color2 := tc.New("not a color")
color2.IsValid() // false
color2.ToString() // "#000000"
```

### GetBrightness
Returns the perceived brightness of a color, from `0-255`, as defined by Web Content Accessibility Guidelines (Version 1.0).
```go
tc.New("#fff").GetBrightness() // 255
tc.New("#000").GetBrightness() // 0
```

### IsLight & IsDark
Return a boolean indicating whether the color's perceived brightness is light or dark.
```go
tc.New("#fff").IsLight() // true
tc.New("#000").IsDark()  // true
```

### GetLuminance
Returns the perceived luminance of a color, from `0-1` as defined by Web Content Accessibility Guidelines (Version 2.0).
```go
tc.New("#fff").GetLuminance() // 1
tc.New("#000").GetLuminance() // 0
```

### GetAlpha / SetAlpha
```go
color := tc.New("red")
color.GetAlpha() // 1.0
color.SetAlpha(0.5)
color.GetAlpha() // 0.5
color.ToRgbString() // "rgba(255, 0, 0, 0.5)"
```

### String Representations
```go
color := tc.New("red")

color.ToHsv() // tc.HSV{H: 0, S: 1, V: 1, A: 1}
color.ToHsvString() // "hsv(0, 100%, 100%)"

color.ToHsl() // tc.HSL{H: 0, S: 1, L: 0.5, A: 1}
color.ToHslString() // "hsl(0, 100%, 50%)"

color.ToHex() // "ff0000"
color.ToHexString() // "#ff0000"
color.ToHex8() // "ff0000ff"
color.ToHex8String() // "#ff0000ff"

color.ToRgb() // tc.RGB{R: 255, G: 0, B: 0, A: 1}
color.ToRgbString() // "rgb(255, 0, 0)"

color.ToPercentageRgb() // tc.RGB{R: 1, G: 0, B: 0, A: 1}
color.ToPercentageRgbString() // "rgb(100%, 0%, 0%)"

color.ToName() // "red", true
color.ToFilter() // "progid:DXImageTransform.Microsoft.gradient(startColorstr=#ffff0000,endColorstr=#ffff0000)"

color.ToString() // "red"
color.ToString(tc.FormatHsv) // "hsv(0, 100%, 100%)"
```

### Color Modification
These methods manipulate the current color, and return it for chaining.
```go
tc.New("red").Lighten().Desaturate().ToHexString() // "#f53d3d"
```

- `Lighten(amount ...float64)`
- `Brighten(amount ...float64)`
- `Darken(amount ...float64)`
- `Desaturate(amount ...float64)`
- `Saturate(amount ...float64)`
- `Greyscale()`
- `Spin(amount float64)`

### Color Combinations
Combination functions return a slice of `*Color` objects (`[]*Color`) unless otherwise noted.
```go
tc.New("#f00").Analogous()
tc.New("#f00").Monochromatic()
tc.New("#f00").SplitComplement()
tc.New("#f00").Triad()
tc.New("#f00").Tetrad()
tc.New("#f00").Complement() // Returns a single *Color
```

### Readability (WCAG)

#### Readability
Returns the contrast ratio between two colors.
```go
tc.Readability(tc.New("#000"), tc.New("#fff")) // 21
```

#### IsReadable
Ensure that foreground and background color combinations meet WCAG guidelines. `ReadabilityOptions` can specify `Level` ("AA" or "AAA") and `Size` ("small" or "large").
```go
tc.IsReadable(tc.New("#000"), tc.New("#111")) // false
tc.IsReadable(tc.New("#ff0088"), tc.New("#5c1a72"), tc.ReadabilityOptions{Level: "AA", Size: "large"}) // true
```

#### MostReadable
Given a base color and a list of possible foreground or background colors for that base, returns the most readable color.
```go
tc.MostReadable(tc.New("#000"), []*tc.Color{tc.New("#f00"), tc.New("#0f0"), tc.New("#00f")}).ToHexString() // "#00ff00"
```

### Clone & Equals
```go
color1 := tc.New("#F00")
color2 := color1.Clone()
color2.SetAlpha(0.5)

tc.Equals(color1, color2) // false
```

## JavaScript → Go API Mapping

The Go API mirrors the original JavaScript API using Idiomatic Go conventions:
- `tinycolor(input)` -> `New(input)`
- `tinycolor.fromRatio(input)` -> `FromRatio(input)`
- `.toRgb()` -> `.ToRgb() RGB`
- `.toHexString()` -> `.ToHexString(...bool) string`
- `.lighten()` -> `.Lighten(...float64) *Color`

See `walkthrough.md` for the full parity mapping.

## Port Notes

The Go port achieves behavioral compatibility with the original JavaScript version by utilizing `interface{}` parameter types to preserve the highly dynamic, permissive input parsing style (strings, objects, percentages) and exact `float64` channel math.

## API Parity

100% API parity has been achieved. All 46 JavaScript exported functions are fully implemented in Go. No intentional behavioral deviations remain beyond unavoidable floating-point representation limits.
