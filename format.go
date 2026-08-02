package tinycolor

import (
	"fmt"
	"math"
)

// ToRgbString returns the RGB string representation.
func (c *Color) ToRgbString() string {
	r := math.Round(c.r)
	g := math.Round(c.g)
	b := math.Round(c.b)

	if c.a == 1 {
		return fmt.Sprintf("rgb(%d, %d, %d)", int(r), int(g), int(b))
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %v)", int(r), int(g), int(b), c.roundA)
}

// ToHslString returns the HSL string representation.
func (c *Color) ToHslString() string {
	hsl := c.ToHsl()
	h := math.Round(hsl.H)
	s := math.Round(hsl.S * 100)
	l := math.Round(hsl.L * 100)

	if c.a == 1 {
		return fmt.Sprintf("hsl(%d, %d%%, %d%%)", int(h), int(s), int(l))
	}
	return fmt.Sprintf("hsla(%d, %d%%, %d%%, %v)", int(h), int(s), int(l), c.roundA)
}

// ToHsvString returns the HSV string representation.
func (c *Color) ToHsvString() string {
	hsv := c.ToHsv()
	h := math.Round(hsv.H)
	s := math.Round(hsv.S * 100)
	v := math.Round(hsv.V * 100)

	if c.a == 1 {
		return fmt.Sprintf("hsv(%d, %d%%, %d%%)", int(h), int(s), int(v))
	}
	return fmt.Sprintf("hsva(%d, %d%%, %d%%, %v)", int(h), int(s), int(v), c.roundA)
}

// ToHexString returns the hex string representation.
func (c *Color) ToHexString(allow3Char ...bool) string {
	if !c.ok {
		return "#000000"
	}
	return "#" + c.ToHex(allow3Char...)
}

// ToHex8String returns the hex8 string representation.
func (c *Color) ToHex8String(allow4Char ...bool) string {
	if !c.ok {
		return "#00000000"
	}
	return "#" + c.ToHex8(allow4Char...)
}

// ToString formats the color based on the provided format.
func (c *Color) ToString(formatOverride ...Format) string {
	format := c.format
	if len(formatOverride) > 0 {
		format = formatOverride[0]
	}

	// Handle invalid color
	if !c.ok {
		return c.ToHexString()
	}

	// Format overrides
	switch format {
	case FormatHex, FormatHex6:
		return c.ToHexString()
	case FormatHex3:
		return c.ToHexString(true)
	case FormatHex4:
		return c.ToHex8String(true)
	case FormatHex8:
		return c.ToHex8String()
	case FormatRgb:
		return c.ToRgbString()
	case FormatPrgb:
		// Percentage RGB is not strictly listed in the API I must implement, but if requested:
		// It would be "rgb(100%, 0%, 0%)". I will implement a basic version or fallback to RGB.
		// Let's implement it for completeness based on JS mod.js
		r := math.Round((c.r / 255) * 100)
		g := math.Round((c.g / 255) * 100)
		b := math.Round((c.b / 255) * 100)
		if c.a == 1 {
			return fmt.Sprintf("rgb(%d%%, %d%%, %d%%)", int(r), int(g), int(b))
		}
		return fmt.Sprintf("rgba(%d%%, %d%%, %d%%, %v)", int(r), int(g), int(b), c.roundA)
	case FormatHsl:
		return c.ToHslString()
	case FormatHsv:
		return c.ToHsvString()
	case FormatName:
		name, found := c.ToName()
		if found {
			return name
		}
		return c.ToHexString()
	}

	// If no format matches, fallback to Hex string (like in JS)
	return c.ToHexString()
}
