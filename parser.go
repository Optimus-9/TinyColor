package tinycolor

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	trimLeft  = regexp.MustCompile(`^\s+`)
	trimRight = regexp.MustCompile(`\s+$`)

	cssInteger = `[-\+]?\d+%?`
	cssNumber  = `[-\+]?\d*\.\d+%?`
	cssUnit    = `(?:` + cssNumber + `)|(?:` + cssInteger + `)`

	permissiveMatch3 = `[\s|\(]+(` + cssUnit + `)[,|\s]+(` + cssUnit + `)[,|\s]+(` + cssUnit + `)\s*\)?`
	permissiveMatch4 = `[\s|\(]+(` + cssUnit + `)[,|\s]+(` + cssUnit + `)[,|\s]+(` + cssUnit + `)[,|\s]+(` + cssUnit + `)\s*\)?`

	matchers = struct {
		cssUnit *regexp.Regexp
		rgb     *regexp.Regexp
		rgba    *regexp.Regexp
		hsl     *regexp.Regexp
		hsla    *regexp.Regexp
		hsv     *regexp.Regexp
		hsva    *regexp.Regexp
		hex3    *regexp.Regexp
		hex6    *regexp.Regexp
		hex4    *regexp.Regexp
		hex8    *regexp.Regexp
	}{
		cssUnit: regexp.MustCompile(`^` + cssUnit + `$`),
		rgb:     regexp.MustCompile(`(?i)^rgb` + permissiveMatch3 + `$`),
		rgba:    regexp.MustCompile(`(?i)^rgba` + permissiveMatch4 + `$`),
		hsl:     regexp.MustCompile(`(?i)^hsl` + permissiveMatch3 + `$`),
		hsla:    regexp.MustCompile(`(?i)^hsla` + permissiveMatch4 + `$`),
		hsv:     regexp.MustCompile(`(?i)^hsv` + permissiveMatch3 + `$`),
		hsva:    regexp.MustCompile(`(?i)^hsva` + permissiveMatch4 + `$`),
		hex3:    regexp.MustCompile(`(?i)^#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})$`),
		hex6:    regexp.MustCompile(`(?i)^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`),
		hex4:    regexp.MustCompile(`(?i)^#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})$`),
		hex8:    regexp.MustCompile(`(?i)^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`),
	}
)

func isValidCSSUnit(color interface{}) bool {
	str := fmt.Sprintf("%v", color)
	return matchers.cssUnit.MatchString(str)
}

func parseIntFromHex(val string) float64 {
	v, _ := strconv.ParseInt(val, 16, 64)
	return float64(v)
}

func convertHexToDecimal(h string) float64 {
	return parseIntFromHex(h) / 255.0
}

func stringInputToObject(color string) (map[string]interface{}, bool) {
	color = trimLeft.ReplaceAllString(color, "")
	color = trimRight.ReplaceAllString(color, "")
	color = strings.ToLower(color)

	named := false
	if hex, ok := Names[color]; ok {
		color = hex
		named = true
	} else if color == "transparent" {
		return map[string]interface{}{"r": 0, "g": 0, "b": 0, "a": 0, "format": "name"}, true
	}

	if match := matchers.rgb.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"r": match[1], "g": match[2], "b": match[3]}, true
	}
	if match := matchers.rgba.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"r": match[1], "g": match[2], "b": match[3], "a": match[4]}, true
	}
	if match := matchers.hsl.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"h": match[1], "s": match[2], "l": match[3]}, true
	}
	if match := matchers.hsla.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"h": match[1], "s": match[2], "l": match[3], "a": match[4]}, true
	}
	if match := matchers.hsv.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"h": match[1], "s": match[2], "v": match[3]}, true
	}
	if match := matchers.hsva.FindStringSubmatch(color); match != nil {
		return map[string]interface{}{"h": match[1], "s": match[2], "v": match[3], "a": match[4]}, true
	}
	if match := matchers.hex8.FindStringSubmatch(color); match != nil {
		format := "hex8"
		if named {
			format = "name"
		}
		return map[string]interface{}{
			"r":      parseIntFromHex(match[1]),
			"g":      parseIntFromHex(match[2]),
			"b":      parseIntFromHex(match[3]),
			"a":      convertHexToDecimal(match[4]),
			"format": format,
		}, true
	}
	if match := matchers.hex6.FindStringSubmatch(color); match != nil {
		format := "hex"
		if named {
			format = "name"
		}
		return map[string]interface{}{
			"r":      parseIntFromHex(match[1]),
			"g":      parseIntFromHex(match[2]),
			"b":      parseIntFromHex(match[3]),
			"format": format,
		}, true
	}
	if match := matchers.hex4.FindStringSubmatch(color); match != nil {
		format := "hex8"
		if named {
			format = "name"
		}
		return map[string]interface{}{
			"r":      parseIntFromHex(match[1] + match[1]),
			"g":      parseIntFromHex(match[2] + match[2]),
			"b":      parseIntFromHex(match[3] + match[3]),
			"a":      convertHexToDecimal(match[4] + match[4]),
			"format": format,
		}, true
	}
	if match := matchers.hex3.FindStringSubmatch(color); match != nil {
		format := "hex"
		if named {
			format = "name"
		}
		return map[string]interface{}{
			"r":      parseIntFromHex(match[1] + match[1]),
			"g":      parseIntFromHex(match[2] + match[2]),
			"b":      parseIntFromHex(match[3] + match[3]),
			"format": format,
		}, true
	}

	return nil, false
}

func convertToPercentage(n interface{}) interface{} {
	switch v := n.(type) {
	case float64:
		if v <= 1 {
			return fmt.Sprintf("%v%%", v*100)
		}
		return v
	case float32:
		if v <= 1 {
			return fmt.Sprintf("%v%%", float64(v)*100)
		}
		return v
	case int:
		if v <= 1 {
			return fmt.Sprintf("%v%%", v*100)
		}
		return v
	}
	return n
}

func isOnePointZero(n interface{}) bool {
	str := fmt.Sprintf("%v", n)
	return strings.Contains(str, ".") && parseFloat(str) == 1
}

func isPercentage(n interface{}) bool {
	str := fmt.Sprintf("%v", n)
	return strings.Contains(str, "%")
}

func parseFloat(n interface{}) float64 {
	switch v := n.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		str := strings.ReplaceAll(v, "%", "")
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return math.NaN()
		}
		return f
	}
	str := fmt.Sprintf("%v", n)
	str = strings.ReplaceAll(str, "%", "")
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

// bound01 converts a CSS value to a normalised float in [0, 1].
// It handles integer values, percentage strings, and the "1.0 == 100%" ratio rule.
func bound01(n interface{}, max float64) float64 {
	if isOnePointZero(n) {
		n = "100%"
	}

	processPercent := isPercentage(n)
	f := parseFloat(n)
	f = math.Min(max, math.Max(0, f))

	if processPercent {
		f = float64(int64(f*max)) / 100.0
	}

	if math.Abs(f-max) < 0.000001 {
		return 1
	}

	return math.Mod(f, max) / max
}

func boundAlpha(a interface{}) float64 {
	f := parseFloat(a)
	if math.IsNaN(f) || f < 0 || f > 1 {
		f = 1
	}
	return f
}

// hue2rgb is a helper for HSL-to-RGB conversion, matching the JS implementation.
func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

type parsedRGB struct {
	r, g, b float64
	a       float64
	ok      bool
	format  Format
}

func hasKey(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

// inputToRGB is the core parser that converts any supported input type into
// internal RGB components. It mirrors the JavaScript inputToRGB function exactly.
func inputToRGB(color interface{}) parsedRGB {
	rgb := parsedRGB{r: 0, g: 0, b: 0, a: 1, ok: false, format: ""}
	var format Format

	// ── Native Go struct inputs (from modify/combination chain) ───────────────
	// These are generated by ToHsl(), ToHsv(), ToRgb() and passed back via New().
	// S, L, V from rgbToHsl/rgbToHsv are already in [0,1], so multiply by 100
	// to feed the [0,100] expected by HslToRgb/HsvToRgb.
	if rgbIn, ok := color.(RGB); ok {
		a := rgbIn.A
		if a <= 0 {
			a = 1.0 // zero-value default: treat as opaque
		}
		return parsedRGB{r: rgbIn.R, g: rgbIn.G, b: rgbIn.B, a: a, ok: true, format: "rgb"}
	}
	if hsl, ok := color.(HSL); ok {
		result := HslToRgb(hsl.H, hsl.S*100, hsl.L*100)
		a := hsl.A
		if a <= 0 {
			a = 1.0
		}
		return parsedRGB{r: result.R, g: result.G, b: result.B, a: a, ok: true, format: "hsl"}
	}
	if hsv, ok := color.(HSV); ok {
		result := HsvToRgb(hsv.H, hsv.S*100, hsv.V*100)
		a := hsv.A
		if a <= 0 {
			a = 1.0
		}
		return parsedRGB{r: result.R, g: result.G, b: result.B, a: a, ok: true, format: "hsv"}
	}

	// ── String input ──────────────────────────────────────────────────────────
	if str, ok := color.(string); ok {
		obj, success := stringInputToObject(str)
		if success {
			color = obj
		}
	}

	// ── Map / object input ────────────────────────────────────────────────────
	if m, ok := color.(map[string]interface{}); ok {
		hasR := hasKey(m, "r")
		hasG := hasKey(m, "g")
		hasB := hasKey(m, "b")
		hasH := hasKey(m, "h")
		hasS := hasKey(m, "s")
		hasV := hasKey(m, "v")
		hasL := hasKey(m, "l")
		hasA := hasKey(m, "a")

		if hasR && hasG && hasB && isValidCSSUnit(m["r"]) && isValidCSSUnit(m["g"]) && isValidCSSUnit(m["b"]) {
			rgb.r = bound01(m["r"], 255) * 255
			rgb.g = bound01(m["g"], 255) * 255
			rgb.b = bound01(m["b"], 255) * 255
			rgb.ok = true
			strR := fmt.Sprintf("%v", m["r"])
			if len(strR) > 0 && strR[len(strR)-1:] == "%" {
				format = "prgb"
			} else {
				format = "rgb"
			}
		} else if hasH && hasS && hasV && isValidCSSUnit(m["h"]) && isValidCSSUnit(m["s"]) && isValidCSSUnit(m["v"]) {
			s := convertToPercentage(m["s"])
			v := convertToPercentage(m["v"])
			// Use NG's exported HsvToRgb; values come through bound01 via the
			// string percentage path so we parse them using parseFloat after converting.
			sf := bound01(s, 100) * 100
			vf := bound01(v, 100) * 100
			hf := bound01(m["h"], 360) * 360
			result := HsvToRgb(hf, sf, vf)
			rgb.r, rgb.g, rgb.b = result.R, result.G, result.B
			rgb.ok = true
			format = "hsv"
		} else if hasH && hasS && hasL && isValidCSSUnit(m["h"]) && isValidCSSUnit(m["s"]) && isValidCSSUnit(m["l"]) {
			s := convertToPercentage(m["s"])
			l := convertToPercentage(m["l"])
			sf := bound01(s, 100) * 100
			lf := bound01(l, 100) * 100
			hf := bound01(m["h"], 360) * 360
			result := HslToRgb(hf, sf, lf)
			rgb.r, rgb.g, rgb.b = result.R, result.G, result.B
			rgb.ok = true
			format = "hsl"
		}

		if hasA {
			rgb.a = boundAlpha(m["a"])
		} else {
			rgb.a = boundAlpha(1)
		}

		if f, ok := m["format"]; ok {
			format = Format(fmt.Sprintf("%v", f))
		}
	} else {
		rgb.a = boundAlpha(1)
	}

	rgb.r = math.Min(255, math.Max(rgb.r, 0))
	rgb.g = math.Min(255, math.Max(rgb.g, 0))
	rgb.b = math.Min(255, math.Max(rgb.b, 0))
	rgb.format = format

	return rgb
}
