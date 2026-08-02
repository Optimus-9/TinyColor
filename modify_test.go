package tinycolor

import (
	"math"
	"testing"
)

// -- MOCKS START --
// Implementing local testing mocks directly inside modify_test.go as per Zero-Conflict File Ownership Rules.

type Options struct {
	Format       Format
	GradientType bool
}

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

type Color struct {
	originalInput interface{}
	r, g, b       float64
	a             float64
	roundA        float64
	format        Format
	gradientType  bool
	ok            bool
}

type HSL struct {
	H, S, L float64
}

type HSV struct {
	H, S, V float64
}

type RGB struct {
	R, G, B float64
	A       float64
}

func New(input interface{}, opts ...Options) *Color {
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
	if hsl, ok := input.(HSL); ok {
		return &Color{ok: true, r: hsl.H, g: hsl.S, b: hsl.L}
	}
	if rgb, ok := input.(RGB); ok {
		return &Color{ok: true, r: rgb.R, g: rgb.G, b: rgb.B}
	}
	if hsv, ok := input.(HSV); ok {
		return &Color{ok: true, r: hsv.H, g: hsv.S, b: hsv.V}
	}
	return &Color{ok: true}
}

func (c *Color) ToHsl() HSL {
	return HSL{H: c.r, S: c.g, L: c.b}
}

func (c *Color) ToRgb() RGB {
	return RGB{R: c.r, G: c.g, B: c.b, A: c.a}
}

func (c *Color) ToHsv() HSV {
	return HSV{H: c.r, S: c.g, V: c.b}
}
// -- MOCKS END --

func round(val float64) float64 {
	return math.Round(val*1000) / 1000
}

func TestModify(t *testing.T) {
	c := New(HSL{H: 100, S: 0.5, L: 0.5})

	t.Run("Lighten", func(t *testing.T) {
		res := c.Lighten(10)
		if round(res.b) != 0.6 {
			t.Errorf("Lighten(10): expected L=0.6, got %v", res.b)
		}
	})

	t.Run("Darken", func(t *testing.T) {
		res := c.Darken(10)
		if round(res.b) != 0.4 {
			t.Errorf("Darken(10): expected L=0.4, got %v", res.b)
		}
	})

	t.Run("Saturate", func(t *testing.T) {
		res := c.Saturate(10)
		if round(res.g) != 0.6 {
			t.Errorf("Saturate(10): expected S=0.6, got %v", res.g)
		}
	})

	t.Run("Desaturate", func(t *testing.T) {
		res := c.Desaturate(10)
		if round(res.g) != 0.4 {
			t.Errorf("Desaturate(10): expected S=0.4, got %v", res.g)
		}
	})

	t.Run("Greyscale", func(t *testing.T) {
		res := c.Greyscale()
		if round(res.g) != 0.0 {
			t.Errorf("Greyscale(): expected S=0.0, got %v", res.g)
		}
	})

	t.Run("Spin", func(t *testing.T) {
		res := c.Spin(10)
		if round(res.r) != 110 {
			t.Errorf("Spin(10): expected H=110, got %v", res.r)
		}

		res2 := c.Spin(-150)
		if round(res2.r) != 310 {
			t.Errorf("Spin(-150): expected H=310, got %v", res2.r)
		}
	})
}

func TestBrighten(t *testing.T) {
	c := New(RGB{R: 100, G: 100, B: 100})
	res := c.Brighten(10)
	if res.r != 126 || res.g != 126 || res.b != 126 {
		t.Errorf("Brighten(10): expected 126, got R=%v G=%v B=%v", res.r, res.g, res.b)
	}
}

func TestCombinations(t *testing.T) {
	c := New(HSL{H: 100, S: 0.5, L: 0.5})

	t.Run("Complement", func(t *testing.T) {
		res := c.Complement()
		if res.r != 280 {
			t.Errorf("expected 280, got %v", res.r)
		}
	})

	t.Run("Polyad", func(t *testing.T) {
		res := c.Polyad(4)
		if len(res) != 4 {
			t.Fatalf("expected 4, got %v", len(res))
		}
		if res[0].r != 100 || res[1].r != 190 || res[2].r != 280 || res[3].r != 10 {
			t.Errorf("unexpected polyad H values: %v %v %v %v", res[0].r, res[1].r, res[2].r, res[3].r)
		}
	})

	t.Run("Triad", func(t *testing.T) {
		res := c.Triad()
		if len(res) != 3 {
			t.Fatalf("expected 3, got %v", len(res))
		}
		if res[0].r != 100 || res[1].r != 220 || res[2].r != 340 {
			t.Errorf("unexpected triad H values: %v %v %v", res[0].r, res[1].r, res[2].r)
		}
	})

	t.Run("Tetrad", func(t *testing.T) {
		res := c.Tetrad()
		if len(res) != 4 {
			t.Fatalf("expected 4, got %v", len(res))
		}
		if res[0].r != 100 || res[1].r != 190 || res[2].r != 280 || res[3].r != 10 {
			t.Errorf("unexpected tetrad H values: %v %v %v %v", res[0].r, res[1].r, res[2].r, res[3].r)
		}
	})

	t.Run("SplitComplement", func(t *testing.T) {
		res := c.SplitComplement()
		if len(res) != 3 {
			t.Fatalf("expected 3, got %v", len(res))
		}
		if res[0].r != 100 || res[1].r != 172 || res[2].r != 316 {
			t.Errorf("unexpected splitcomplement H values")
		}
	})

	t.Run("Analogous", func(t *testing.T) {
		res := c.Analogous(6, 30)
		if len(res) != 6 {
			t.Fatalf("expected 6, got %v", len(res))
		}
		if res[1].r != 76 || res[2].r != 88 || res[3].r != 100 || res[4].r != 112 || res[5].r != 124 {
			t.Errorf("unexpected analogous values")
		}
	})

	t.Run("Monochromatic", func(t *testing.T) {
		c2 := New(HSV{H: 100, S: 0.5, V: 0.2})
		res := c2.Monochromatic(6)
		if len(res) != 6 {
			t.Fatalf("expected 6, got %v", len(res))
		}
		if round(res[0].b) != 0.2 || round(res[1].b) != 0.367 || round(res[2].b) != 0.533 || round(res[3].b) != 0.7 || round(res[4].b) != 0.867 || round(res[5].b) != 0.033 {
			t.Errorf("unexpected monochromatic V values: %v %v %v %v %v %v", round(res[0].b), round(res[1].b), round(res[2].b), round(res[3].b), round(res[4].b), round(res[5].b))
		}
	})
}
