package tinycolor

import "math"

func clamp01(val float64) float64 {
	return math.Min(1, math.Max(0, val))
}

func (c *Color) Lighten(amount ...float64) *Color {
	amt := 10.0
	if len(amount) > 0 {
		amt = amount[0]
	}
	hsl := c.ToHsl()
	hsl.L += amt / 100
	hsl.L = clamp01(hsl.L)
	return New(hsl)
}

func (c *Color) Brighten(amount ...float64) *Color {
	amt := 10.0
	if len(amount) > 0 {
		amt = amount[0]
	}
	rgb := c.ToRgb()
	rgb.R = math.Max(0, math.Min(255, rgb.R-math.Round(255*-(amt/100))))
	rgb.G = math.Max(0, math.Min(255, rgb.G-math.Round(255*-(amt/100))))
	rgb.B = math.Max(0, math.Min(255, rgb.B-math.Round(255*-(amt/100))))
	return New(rgb)
}

func (c *Color) Darken(amount ...float64) *Color {
	amt := 10.0
	if len(amount) > 0 {
		amt = amount[0]
	}
	hsl := c.ToHsl()
	hsl.L -= amt / 100
	hsl.L = clamp01(hsl.L)
	return New(hsl)
}

func (c *Color) Desaturate(amount ...float64) *Color {
	amt := 10.0
	if len(amount) > 0 {
		amt = amount[0]
	}
	hsl := c.ToHsl()
	hsl.S -= amt / 100
	hsl.S = clamp01(hsl.S)
	return New(hsl)
}

func (c *Color) Saturate(amount ...float64) *Color {
	amt := 10.0
	if len(amount) > 0 {
		amt = amount[0]
	}
	hsl := c.ToHsl()
	hsl.S += amt / 100
	hsl.S = clamp01(hsl.S)
	return New(hsl)
}

func (c *Color) Greyscale() *Color {
	return c.Desaturate(100)
}

func (c *Color) Spin(amount float64) *Color {
	hsl := c.ToHsl()
	hue := math.Mod(hsl.H+amount, 360)
	if hue < 0 {
		hue += 360
	}
	hsl.H = hue
	return New(hsl)
}
