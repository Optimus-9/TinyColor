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

func inputToRGB(color interface{}) parsedRGB {
	rgb := parsedRGB{r: 0, g: 0, b: 0, a: 1, ok: false, format: ""}
	var format Format

	if str, ok := color.(string); ok {
		obj, success := stringInputToObject(str)
		if success {
			color = obj
		}
	}

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
			rgb.r, rgb.g, rgb.b = stubHsvToRgb(m["h"], s, v)
			rgb.ok = true
			format = "hsv"
		} else if hasH && hasS && hasL && isValidCSSUnit(m["h"]) && isValidCSSUnit(m["s"]) && isValidCSSUnit(m["l"]) {
			s := convertToPercentage(m["s"])
			l := convertToPercentage(m["l"])
			rgb.r, rgb.g, rgb.b = stubHslToRgb(m["h"], s, l)
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

// stubHsvToRgb temporary implementation for BB
func stubHsvToRgb(h, s, v interface{}) (float64, float64, float64) {
	hf := bound01(h, 360) * 6
	sf := bound01(s, 100)
	vf := bound01(v, 100)

	i := math.Floor(hf)
	f := hf - i
	p := vf * (1 - sf)
	q := vf * (1 - f*sf)
	t := vf * (1 - (1-f)*sf)
	mod := int(i) % 6
	r, g, b := 0.0, 0.0, 0.0
	switch mod {
	case 0:
		r, g, b = vf, t, p
	case 1:
		r, g, b = q, vf, p
	case 2:
		r, g, b = p, vf, t
	case 3:
		r, g, b = p, q, vf
	case 4:
		r, g, b = t, p, vf
	case 5:
		r, g, b = vf, p, q
	}
	return r * 255, g * 255, b * 255
}

// stubHslToRgb temporary implementation for BB
func stubHslToRgb(h, s, l interface{}) (float64, float64, float64) {
	hf := bound01(h, 360)
	sf := bound01(s, 100)
	lf := bound01(l, 100)
	r, g, b := 0.0, 0.0, 0.0

	if sf == 0 {
		r, g, b = lf, lf, lf
	} else {
		var q float64
		if lf < 0.5 {
			q = lf * (1 + sf)
		} else {
			q = lf + sf - lf*sf
		}
		p := 2*lf - q
		r = hue2rgb(p, q, hf+1.0/3.0)
		g = hue2rgb(p, q, hf)
		b = hue2rgb(p, q, hf-1.0/3.0)
	}
	return r * 255, g * 255, b * 255
}

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
