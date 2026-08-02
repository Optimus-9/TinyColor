package tinycolor

import "math"

func (c *Color) Complement() *Color {
	hsl := c.ToHsl()
	hsl.H = math.Mod(hsl.H+180, 360)
	return New(hsl)
}

func (c *Color) Polyad(number int) []*Color {
	if number <= 0 {
		panic("Argument to polyad must be a positive number")
	}
	hsl := c.ToHsl()
	result := []*Color{New(c)}
	step := 360.0 / float64(number)
	for i := 1; i < number; i++ {
		result = append(result, New(HSL{
			H: math.Mod(hsl.H+float64(i)*step, 360),
			S: hsl.S,
			L: hsl.L,
		}))
	}
	return result
}

func (c *Color) Triad() []*Color {
	return c.Polyad(3)
}

func (c *Color) Tetrad() []*Color {
	return c.Polyad(4)
}

func (c *Color) SplitComplement() []*Color {
	return []*Color{
		New(c),
		c.Spin(72),
		c.Spin(216),
	}
}

func (c *Color) Analogous(results, slices int) []*Color {
	if results == 0 {
		results = 6
	}
	if slices == 0 {
		slices = 30
	}

	hsl := c.ToHsl()
	part := 360.0 / float64(slices)
	ret := []*Color{New(c)}

	// (part * results) >> 1
	shift := int(part*float64(results)) >> 1
	hsl.H = math.Mod(hsl.H-float64(shift)+720, 360)

	for results > 1 {
		hsl.H = math.Mod(hsl.H+part, 360)
		ret = append(ret, New(hsl))
		results--
	}
	return ret
}

func (c *Color) Monochromatic(results int) []*Color {
	if results == 0 {
		results = 6
	}
	hsv := c.ToHsv()
	h, s, v := hsv.H, hsv.S, hsv.V
	ret := []*Color{}
	modification := 1.0 / float64(results)

	for results > 0 {
		ret = append(ret, New(HSV{H: h, S: s, V: v}))
		v = math.Mod(v+modification, 1)
		results--
	}
	return ret
}
