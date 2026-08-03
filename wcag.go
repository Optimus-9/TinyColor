package tinycolor

import (
	"math"
)

// GetBrightness returns the perceived brightness of a color.
func (c *Color) GetBrightness() float64 {
	rgb := c.ToRgb()
	return (rgb.R*299 + rgb.G*587 + rgb.B*114) / 1000
}

// IsDark returns true if the color is perceived as dark.
func (c *Color) IsDark() bool {
	return c.GetBrightness() < 128
}

// IsLight returns true if the color is perceived as light.
func (c *Color) IsLight() bool {
	return !c.IsDark()
}

// GetLuminance returns the relative luminance of a color.
func (c *Color) GetLuminance() float64 {
	rgb := c.ToRgb()
	rsRGB := rgb.R / 255
	gsRGB := rgb.G / 255
	bsRGB := rgb.B / 255

	var R, G, B float64

	if rsRGB <= 0.03928 {
		R = rsRGB / 12.92
	} else {
		R = math.Pow((rsRGB+0.055)/1.055, 2.4)
	}

	if gsRGB <= 0.03928 {
		G = gsRGB / 12.92
	} else {
		G = math.Pow((gsRGB+0.055)/1.055, 2.4)
	}

	if bsRGB <= 0.03928 {
		B = bsRGB / 12.92
	} else {
		B = math.Pow((bsRGB+0.055)/1.055, 2.4)
	}

	return 0.2126*R + 0.7152*G + 0.0722*B
}

// Readability calculates the contrast ratio between two colors.
func Readability(c1, c2 *Color) float64 {
	l1 := c1.GetLuminance()
	l2 := c2.GetLuminance()
	return (math.Max(l1, l2) + 0.05) / (math.Min(l1, l2) + 0.05)
}

// IsReadable determines if a color combination meets WCAG guidelines.
func IsReadable(c1, c2 *Color, opts ...ReadabilityOptions) bool {
	readability := Readability(c1, c2)

	level := "AA"
	size := "small"

	if len(opts) > 0 {
		if opts[0].Level != "" {
			level = opts[0].Level
		}
		if opts[0].Size != "" {
			size = opts[0].Size
		}
	}

	out := false
	switch level + size {
	case "AAsmall", "AAAlarge":
		out = readability >= 4.5
	case "AAlarge":
		out = readability >= 3.0
	case "AAAsmall":
		out = readability >= 7.0
	}
	return out
}

// MostReadable returns the most readable color from a list of options.
func MostReadable(baseColor *Color, colorList []*Color, opts ...ReadabilityOptions) *Color {
	var bestColor *Color
	bestScore := 0.0

	for _, c := range colorList {
		score := Readability(baseColor, c)
		if score > bestScore {
			bestScore = score
			bestColor = c
		}
	}

	var options ReadabilityOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	if IsReadable(baseColor, bestColor, options) || !options.IncludeFallbackColors {
		return bestColor
	}

	// IncludeFallbackColors is true, and the best color isn't readable
	white := &Color{r: 255, g: 255, b: 255, a: 1, ok: true}
	black := &Color{r: 0, g: 0, b: 0, a: 1, ok: true}

	fallbackOptions := options
	fallbackOptions.IncludeFallbackColors = false

	return MostReadable(baseColor, []*Color{white, black}, fallbackOptions)
}
