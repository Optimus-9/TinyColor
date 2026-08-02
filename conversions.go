package tinycolor

import (
	"fmt"
	"math"
)

// pad2 pads a hex string with a leading zero if it has only one character
func pad2(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func convertDecimalToHex(d float64) string {
	return fmt.Sprintf("%x", int(math.Round(d*255)))
}

func mathMax3(a, b, c float64) float64 {
	return math.Max(a, math.Max(b, c))
}

func mathMin3(a, b, c float64) float64 {
	return math.Min(a, math.Min(b, c))
}

// bound01 takes input in [0, max] and returns it as [0, 1].
func bound01(n, max float64) float64 {
	n = math.Min(max, math.Max(0, n))
	if math.Abs(n-max) < 0.000001 {
		return 1.0
	}
	return math.Mod(n, max) / max
}

// hue2rgb is a helper for HSL to RGB conversion
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

// HslToRgb converts an HSL color value to RGB.
// Assumes h is in [0, 360] and s and l are in [0, 100].
// Returns RGB with R, G, B in the range [0, 255] and A = 1.0.
func HslToRgb(h, s, l float64) RGB {
	h = bound01(h, 360)
	s = bound01(s, 100)
	l = bound01(l, 100)

	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l // achromatic
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	return RGB{
		R: r * 255,
		G: g * 255,
		B: b * 255,
		A: 1.0,
	}
}

// HsvToRgb converts an HSV color value to RGB.
// Assumes h is in [0, 360] and s and v are in [0, 100].
// Returns RGB with R, G, B in the range [0, 255] and A = 1.0.
func HsvToRgb(h, s, v float64) RGB {
	h = bound01(h, 360) * 6
	s = bound01(s, 100)
	v = bound01(v, 100)

	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	mod := int(i) % 6

	rVals := [6]float64{v, q, p, p, t, v}
	gVals := [6]float64{t, v, v, q, p, p}
	bVals := [6]float64{p, p, t, v, v, q}

	return RGB{
		R: rVals[mod] * 255,
		G: gVals[mod] * 255,
		B: bVals[mod] * 255,
		A: 1.0,
	}
}

// rgbToHsl converts RGB to HSL.
// r, g, b are [0, 255]. Returns h in [0, 360], s and l in [0, 1].
func rgbToHsl(r, g, b float64) HSL {
	r = bound01(r, 255)
	g = bound01(g, 255)
	b = bound01(b, 255)

	max := mathMax3(r, g, b)
	min := mathMin3(r, g, b)
	l := (max + min) / 2

	var h, s float64

	if max == min {
		h, s = 0, 0 // achromatic
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return HSL{H: h * 360, S: s, L: l}
}

// rgbToHsv converts RGB to HSV.
// r, g, b are [0, 255]. Returns h in [0, 360], s and v in [0, 1].
func rgbToHsv(r, g, b float64) HSV {
	r = bound01(r, 255)
	g = bound01(g, 255)
	b = bound01(b, 255)

	max := mathMax3(r, g, b)
	min := mathMin3(r, g, b)
	v := max

	d := max - min
	var s float64
	if max == 0 {
		s = 0
	} else {
		s = d / max
	}

	var h float64
	if max == min {
		h = 0 // achromatic
	} else {
		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return HSV{H: h * 360, S: s, V: v}
}

// rgbToHex converts RGB to Hex.
func rgbToHex(r, g, b float64, allow3Char bool) string {
	hex := []string{
		pad2(fmt.Sprintf("%x", int(math.Round(r)))),
		pad2(fmt.Sprintf("%x", int(math.Round(g)))),
		pad2(fmt.Sprintf("%x", int(math.Round(b)))),
	}

	if allow3Char &&
		hex[0][0] == hex[0][1] &&
		hex[1][0] == hex[1][1] &&
		hex[2][0] == hex[2][1] {
		return string([]byte{hex[0][0], hex[1][0], hex[2][0]})
	}

	return hex[0] + hex[1] + hex[2]
}

// rgbaToHex converts RGBA to Hex.
func rgbaToHex(r, g, b, a float64, allow4Char bool) string {
	hex := []string{
		pad2(fmt.Sprintf("%x", int(math.Round(r)))),
		pad2(fmt.Sprintf("%x", int(math.Round(g)))),
		pad2(fmt.Sprintf("%x", int(math.Round(b)))),
		pad2(convertDecimalToHex(a)),
	}

	if allow4Char &&
		hex[0][0] == hex[0][1] &&
		hex[1][0] == hex[1][1] &&
		hex[2][0] == hex[2][1] &&
		hex[3][0] == hex[3][1] {
		return string([]byte{hex[0][0], hex[1][0], hex[2][0], hex[3][0]})
	}

	return hex[0] + hex[1] + hex[2] + hex[3]
}

// rgbaToArgbHex converts RGBA to ARGB Hex.
func rgbaToArgbHex(r, g, b, a float64) string {
	hex := []string{
		pad2(convertDecimalToHex(a)),
		pad2(fmt.Sprintf("%x", int(math.Round(r)))),
		pad2(fmt.Sprintf("%x", int(math.Round(g)))),
		pad2(fmt.Sprintf("%x", int(math.Round(b)))),
	}
	return hex[0] + hex[1] + hex[2] + hex[3]
}

// Exported Conversion Methods on *Color

func (c *Color) ToRgb() RGB {
	return RGB{
		R: math.Round(c.r),
		G: math.Round(c.g),
		B: math.Round(c.b),
		A: c.a,
	}
}

func (c *Color) ToHsl() HSL {
	hsl := rgbToHsl(c.r, c.g, c.b)
	hsl.A = c.a
	return hsl
}

func (c *Color) ToHsv() HSV {
	hsv := rgbToHsv(c.r, c.g, c.b)
	hsv.A = c.a
	return hsv
}

func (c *Color) ToHex(allow3Char ...bool) string {
	a3c := false
	if len(allow3Char) > 0 {
		a3c = allow3Char[0]
	}
	return rgbToHex(c.r, c.g, c.b, a3c)
}

func (c *Color) ToHex8(allow4Char ...bool) string {
	a4c := false
	if len(allow4Char) > 0 {
		a4c = allow4Char[0]
	}
	return rgbaToHex(c.r, c.g, c.b, c.a, a4c)
}

func (c *Color) ToName() (string, bool) {
	if c.a == 0 {
		return "transparent", true
	}
	if c.a < 1 {
		return "", false
	}
	hex := c.ToHex()
	for k, v := range hexNames {
		if k == hex {
			return v, true
		}
	}
	return "", false
}
