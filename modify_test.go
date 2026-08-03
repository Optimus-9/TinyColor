package tinycolor

import (
	"math"
	"testing"
)

// roundTo3 rounds a float to 3 decimal places for stable comparisons.
func roundTo3(val float64) float64 {
	return math.Round(val*1000) / 1000
}

// approxEqual returns true if a and b are within tolerance.
func approxEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestModify(t *testing.T) {
	// Use map input to correctly initialise a color with H=100, S=50%, L=50%.
	c := New(map[string]interface{}{"h": 100, "s": 50, "l": 50})

	t.Run("Lighten", func(t *testing.T) {
		res := c.Lighten(10)
		hsl := res.ToHsl()
		if !approxEqual(hsl.L, 0.6, 0.002) {
			t.Errorf("Lighten(10): expected L≈0.6, got %v", hsl.L)
		}
	})

	t.Run("Darken", func(t *testing.T) {
		res := c.Darken(10)
		hsl := res.ToHsl()
		if !approxEqual(hsl.L, 0.4, 0.002) {
			t.Errorf("Darken(10): expected L≈0.4, got %v", hsl.L)
		}
	})

	t.Run("Saturate", func(t *testing.T) {
		res := c.Saturate(10)
		hsl := res.ToHsl()
		if !approxEqual(hsl.S, 0.6, 0.002) {
			t.Errorf("Saturate(10): expected S≈0.6, got %v", hsl.S)
		}
	})

	t.Run("Desaturate", func(t *testing.T) {
		res := c.Desaturate(10)
		hsl := res.ToHsl()
		if !approxEqual(hsl.S, 0.4, 0.002) {
			t.Errorf("Desaturate(10): expected S≈0.4, got %v", hsl.S)
		}
	})

	t.Run("Greyscale", func(t *testing.T) {
		res := c.Greyscale()
		hsl := res.ToHsl()
		if !approxEqual(hsl.S, 0.0, 0.002) {
			t.Errorf("Greyscale(): expected S≈0.0, got %v", hsl.S)
		}
	})

	t.Run("Spin", func(t *testing.T) {
		res := c.Spin(10)
		hsl := res.ToHsl()
		if !approxEqual(hsl.H, 110, 1.0) {
			t.Errorf("Spin(10): expected H≈110, got %v", hsl.H)
		}

		res2 := c.Spin(-150)
		hsl2 := res2.ToHsl()
		if !approxEqual(hsl2.H, 310, 1.0) {
			t.Errorf("Spin(-150): expected H≈310, got %v", hsl2.H)
		}
	})
}

func TestBrighten(t *testing.T) {
	// RGB struct is now handled natively by inputToRGB.
	c := New(map[string]interface{}{"r": 100, "g": 100, "b": 100})
	res := c.Brighten(10)
	rgb := res.ToRgb()
	if !approxEqual(rgb.R, 126, 1) || !approxEqual(rgb.G, 126, 1) || !approxEqual(rgb.B, 126, 1) {
		t.Errorf("Brighten(10): expected RGB≈(126,126,126), got R=%v G=%v B=%v", rgb.R, rgb.G, rgb.B)
	}
}

func TestCombinations(t *testing.T) {
	c := New(map[string]interface{}{"h": 100, "s": 50, "l": 50})

	t.Run("Complement", func(t *testing.T) {
		res := c.Complement()
		hsl := res.ToHsl()
		if !approxEqual(hsl.H, 280, 1.0) {
			t.Errorf("Complement: expected H≈280, got %v", hsl.H)
		}
	})

	t.Run("Polyad", func(t *testing.T) {
		res := c.Polyad(4)
		if len(res) != 4 {
			t.Fatalf("Polyad(4): expected 4 colors, got %v", len(res))
		}
		expected := []float64{100, 190, 280, 10}
		for i, ex := range expected {
			h := res[i].ToHsl().H
			if !approxEqual(h, ex, 1.5) {
				t.Errorf("Polyad[%d]: expected H≈%v, got %v", i, ex, h)
			}
		}
	})

	t.Run("Triad", func(t *testing.T) {
		res := c.Triad()
		if len(res) != 3 {
			t.Fatalf("Triad: expected 3 colors, got %v", len(res))
		}
		expected := []float64{100, 220, 340}
		for i, ex := range expected {
			h := res[i].ToHsl().H
			if !approxEqual(h, ex, 1.5) {
				t.Errorf("Triad[%d]: expected H≈%v, got %v", i, ex, h)
			}
		}
	})

	t.Run("Tetrad", func(t *testing.T) {
		res := c.Tetrad()
		if len(res) != 4 {
			t.Fatalf("Tetrad: expected 4 colors, got %v", len(res))
		}
		expected := []float64{100, 190, 280, 10}
		for i, ex := range expected {
			h := res[i].ToHsl().H
			if !approxEqual(h, ex, 1.5) {
				t.Errorf("Tetrad[%d]: expected H≈%v, got %v", i, ex, h)
			}
		}
	})

	t.Run("SplitComplement", func(t *testing.T) {
		res := c.SplitComplement()
		if len(res) != 3 {
			t.Fatalf("SplitComplement: expected 3 colors, got %v", len(res))
		}
		expected := []float64{100, 172, 316}
		for i, ex := range expected {
			h := res[i].ToHsl().H
			if !approxEqual(h, ex, 1.5) {
				t.Errorf("SplitComplement[%d]: expected H≈%v, got %v", i, ex, h)
			}
		}
	})

	t.Run("Analogous", func(t *testing.T) {
		res := c.Analogous(6, 30)
		if len(res) != 6 {
			t.Fatalf("Analogous: expected 6, got %v", len(res))
		}
		expectedH := []float64{100, 76, 88, 100, 112, 124}
		// Note: res[0] is New(c) which keeps H=100; subsequent are the shifted steps.
		for i, ex := range expectedH {
			h := res[i].ToHsl().H
			if !approxEqual(h, ex, 2.0) {
				t.Errorf("Analogous[%d]: expected H≈%v, got %v", i, ex, h)
			}
		}
	})

	t.Run("Monochromatic", func(t *testing.T) {
		// Create via string input to avoid struct-zero alpha ambiguity.
		c2 := New(map[string]interface{}{"h": 100, "s": 50, "v": 20})
		res := c2.Monochromatic(6)
		if len(res) != 6 {
			t.Fatalf("Monochromatic: expected 6, got %v", len(res))
		}
		// V values step by 1/6 ≈ 0.1667; allow ±0.01 to account for RGB round-trip.
		expectedV := []float64{0.2, 0.367, 0.533, 0.700, 0.867, 0.033}
		for i, ex := range expectedV {
			v := res[i].ToHsv().V
			if !approxEqual(v, ex, 0.01) {
				t.Errorf("Monochromatic[%d]: expected V≈%v, got %v", i, ex, roundTo3(v))
			}
		}
	})
}
