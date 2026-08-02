package tinycolor

import (
	"math"
)

type Format string

const (
	FormatHex   Format = "hex"
	FormatHex3  Format = "hex3"
	FormatHex4  Format = "hex4"
	FormatHex6  Format = "hex6"
	FormatHex8  Format = "hex8"
	FormatRgb   Format = "rgb"
	FormatPrgb  Format = "prgb"
	FormatHsl   Format = "hsl"
	FormatHsv   Format = "hsv"
	FormatName  Format = "name"
)

type Options struct {
	Format       Format
	GradientType bool
}

type Color struct {
	originalInput interface{}
	r, g, b       float64 // [0, 255]
	a             float64 // [0.0, 1.0]
	roundA        float64 // Cached rounded alpha: Math.Round(100 * a) / 100
	format        Format
	gradientType  bool
	ok            bool // Indicates parsing success
}

func New(input interface{}, opts ...Options) *Color {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	// If input is already a tinycolor, return a copy of it
	if c, ok := input.(*Color); ok {
		return &Color{
			originalInput: c.originalInput,
			r:             c.r,
			g:             c.g,
			b:             c.b,
			a:             c.a,
			roundA:        c.roundA,
			format:        c.format,
			gradientType:  c.gradientType,
			ok:            c.ok,
		}
	}

	// Parse input
	rgb := inputToRGB(input)

	c := &Color{
		originalInput: input,
		r:             rgb.r,
		g:             rgb.g,
		b:             rgb.b,
		a:             rgb.a,
		roundA:        math.Round(rgb.a*100) / 100,
		format:        opt.Format,
		gradientType:  opt.GradientType,
		ok:            rgb.ok,
	}

	if c.format == "" {
		c.format = rgb.format
	}

	// Don't let the range of [0,255] come back in [0,1].
	// Potentially lose a little bit of precision here, but will fix issues where
	// .5 gets interpreted as half of the total, instead of half of 1
	// If it was supposed to be 128, this was already taken care of by `inputToRGB`
	if c.r < 1 {
		c.r = math.Round(c.r)
	}
	if c.g < 1 {
		c.g = math.Round(c.g)
	}
	if c.b < 1 {
		c.b = math.Round(c.b)
	}

	return c
}

// FromRatio takes input from [0, 1] and converts to tinycolor
func FromRatio(input interface{}, opts ...Options) *Color {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if m, ok := input.(map[string]interface{}); ok {
		newInput := make(map[string]interface{})
		for k, v := range m {
			if k == "a" {
				newInput[k] = v
			} else {
				newInput[k] = convertToPercentage(v)
			}
		}
		return New(newInput, opt)
	}

	return New(input, opt)
}

// Struct getters & methods
func (c *Color) IsValid() bool                 { return c.ok }
func (c *Color) GetOriginalInput() interface{} { return c.originalInput }
func (c *Color) GetFormat() Format             { return c.format }
func (c *Color) GetAlpha() float64             { return c.a }
