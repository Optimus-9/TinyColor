// Package tinycolor_test provides black-box integration tests for the tinycolor package.
// Tests are derived from the original JavaScript test suite (test.js).
package tinycolor_test

import (
	"math"
	"testing"

	tc "tinycolor"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

func approx(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

// ─── Initialization ───────────────────────────────────────────────────────────

func TestInitialization(t *testing.T) {
	c := tc.New("red")
	if !c.IsValid() {
		t.Errorf("tinycolor(\"red\") should be valid")
	}
}

// ─── Hex Parsing ─────────────────────────────────────────────────────────────

func TestHexParsing(t *testing.T) {
	cases := []struct {
		input   string
		wantR   float64
		wantG   float64
		wantB   float64
		wantFmt tc.Format
	}{
		{"#f00", 255, 0, 0, tc.FormatHex},
		{"f00", 255, 0, 0, tc.FormatHex},
		{"#ff0000", 255, 0, 0, tc.FormatHex},
		{"ff0000", 255, 0, 0, tc.FormatHex},
		{"#f00f", 255, 0, 0, tc.FormatHex8},
		{"#ff0000ff", 255, 0, 0, tc.FormatHex8},
		{"#000", 0, 0, 0, tc.FormatHex},
		{"#ffffff", 255, 255, 255, tc.FormatHex},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			color := tc.New(c.input)
			if !color.IsValid() {
				t.Fatalf("%s: expected valid color", c.input)
			}
			rgb := color.ToRgb()
			if !approx(rgb.R, c.wantR, 0.5) || !approx(rgb.G, c.wantG, 0.5) || !approx(rgb.B, c.wantB, 0.5) {
				t.Errorf("%s: expected RGB(%v,%v,%v), got RGB(%v,%v,%v)",
					c.input, c.wantR, c.wantG, c.wantB, rgb.R, rgb.G, rgb.B)
			}
			if color.GetFormat() != c.wantFmt {
				t.Errorf("%s: expected format %q, got %q", c.input, c.wantFmt, color.GetFormat())
			}
		})
	}
}

// ─── Named Colors ─────────────────────────────────────────────────────────────

func TestNamedColors(t *testing.T) {
	cases := []struct{ name, hex string }{
		{"red", "#ff0000"},
		{"green", "#008000"},
		{"blue", "#0000ff"},
		{"white", "#ffffff"},
		{"black", "#000000"},
		{"aliceblue", "#f0f8ff"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			color := tc.New(c.name)
			if !color.IsValid() {
				t.Fatalf("Named color %q should be valid", c.name)
			}
			got := color.ToHexString()
			if got != c.hex {
				t.Errorf("%s: expected %s, got %s", c.name, c.hex, got)
			}
		})
	}
}

// ─── Transparent ──────────────────────────────────────────────────────────────

func TestTransparent(t *testing.T) {
	c := tc.New("transparent")
	if !c.IsValid() {
		t.Fatal("transparent should be valid")
	}
	rgb := c.ToRgb()
	if rgb.A != 0 {
		t.Errorf("transparent: expected alpha=0, got %v", rgb.A)
	}
	if c.GetFormat() != tc.FormatName {
		t.Errorf("transparent: expected format name, got %q", c.GetFormat())
	}
}

// ─── RGB String Parsing ───────────────────────────────────────────────────────

func TestRGBStringParsing(t *testing.T) {
	cases := []struct {
		input      string
		r, g, b, a float64
	}{
		{"rgb(255, 0, 0)", 255, 0, 0, 1},
		{"rgb 255 0 0", 255, 0, 0, 1},
		{"rgba(255, 0, 0, 0.5)", 255, 0, 0, 0.5},
		{"rgba 255 0 0 0.5", 255, 0, 0, 0.5},
		{"rgb(100%, 0%, 0%)", 255, 0, 0, 1},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			color := tc.New(c.input)
			if !color.IsValid() {
				t.Fatalf("%q: expected valid, got invalid", c.input)
			}
			rgb := color.ToRgb()
			if !approx(rgb.R, c.r, 0.5) || !approx(rgb.G, c.g, 0.5) || !approx(rgb.B, c.b, 0.5) {
				t.Errorf("%q: expected RGB(%v,%v,%v), got (%v,%v,%v)",
					c.input, c.r, c.g, c.b, rgb.R, rgb.G, rgb.B)
			}
			if !approx(rgb.A, c.a, 0.001) {
				t.Errorf("%q: expected A=%v, got %v", c.input, c.a, rgb.A)
			}
		})
	}
}

// ─── HSL String Parsing ───────────────────────────────────────────────────────

func TestHSLStringParsing(t *testing.T) {
	c := tc.New("hsl(0, 100%, 50%)")
	if !c.IsValid() {
		t.Fatal("hsl(0,100%,50%) should be valid")
	}
	rgb := c.ToRgb()
	if !approx(rgb.R, 255, 0.5) || !approx(rgb.G, 0, 0.5) || !approx(rgb.B, 0, 0.5) {
		t.Errorf("hsl(0,100%%,50%%) expected red, got (%v,%v,%v)", rgb.R, rgb.G, rgb.B)
	}

	ca := tc.New("hsla(0, 100%, 50%, 0.5)")
	if !approx(ca.GetAlpha(), 0.5, 0.001) {
		t.Errorf("hsla alpha: expected 0.5, got %v", ca.GetAlpha())
	}
}

// ─── HSV String Parsing ───────────────────────────────────────────────────────

func TestHSVStringParsing(t *testing.T) {
	c := tc.New("hsv(0, 100%, 100%)")
	if !c.IsValid() {
		t.Fatal("hsv(0,100%,100%) should be valid")
	}
	rgb := c.ToRgb()
	if !approx(rgb.R, 255, 0.5) || !approx(rgb.G, 0, 0.5) || !approx(rgb.B, 0, 0.5) {
		t.Errorf("hsv(0,100%%,100%%) expected red, got (%v,%v,%v)", rgb.R, rgb.G, rgb.B)
	}
}

// ─── Invalid Inputs ───────────────────────────────────────────────────────────

func TestInvalidInputs(t *testing.T) {
	cases := []interface{}{"not a color", "nope", "", "# 123456"}
	for _, input := range cases {
		c := tc.New(input)
		if c.IsValid() {
			t.Errorf("input %v should be invalid", input)
		}
	}
}

// ─── Alpha Handling ───────────────────────────────────────────────────────────

func TestAlphaHandling(t *testing.T) {
	// Out-of-range alpha defaults to 1.
	c := tc.New(map[string]interface{}{"r": 255, "g": 0, "b": 0, "a": -1})
	if c.GetAlpha() != 1 {
		t.Errorf("negative alpha should default to 1, got %v", c.GetAlpha())
	}

	c2 := tc.New(map[string]interface{}{"r": 255, "g": 0, "b": 0, "a": 2})
	if c2.GetAlpha() != 1 {
		t.Errorf("alpha > 1 should default to 1, got %v", c2.GetAlpha())
	}

	// Valid alpha 0.5 passes through.
	c3 := tc.New(map[string]interface{}{"r": 255, "g": 0, "b": 0, "a": 0.5})
	if !approx(c3.GetAlpha(), 0.5, 0.001) {
		t.Errorf("alpha=0.5 should be preserved, got %v", c3.GetAlpha())
	}
}

// ─── Clone & Equals ───────────────────────────────────────────────────────────

func TestCloneAndEquals(t *testing.T) {
	c := tc.New("red")
	clone := c.Clone()

	if !c.Equals(clone) {
		t.Error("clone should equal original")
	}

	// Mutation of clone should not affect original.
	lightened := clone.Lighten(50)
	if c.Equals(lightened) {
		t.Error("lightened clone should not equal original")
	}
}

// ─── SetAlpha ─────────────────────────────────────────────────────────────────

func TestSetAlpha(t *testing.T) {
	c := tc.New("red").SetAlpha(0.5)
	if !approx(c.GetAlpha(), 0.5, 0.001) {
		t.Errorf("SetAlpha(0.5): expected 0.5, got %v", c.GetAlpha())
	}
}

// ─── String Formatting ────────────────────────────────────────────────────────

func TestStringFormatting(t *testing.T) {
	red := tc.New("red")

	if red.ToHexString() != "#ff0000" {
		t.Errorf("ToHexString: expected #ff0000, got %s", red.ToHexString())
	}
	if red.ToRgbString() != "rgb(255, 0, 0)" {
		t.Errorf("ToRgbString: expected 'rgb(255, 0, 0)', got %s", red.ToRgbString())
	}
	if red.ToHslString() != "hsl(0, 100%, 50%)" {
		t.Errorf("ToHslString: expected 'hsl(0, 100%%, 50%%)', got %s", red.ToHslString())
	}

	// With alpha
	redAlpha := tc.New("rgba(255, 0, 0, 0.5)")
	if redAlpha.ToRgbString() != "rgba(255, 0, 0, 0.5)" {
		t.Errorf("ToRgbString alpha: expected 'rgba(255, 0, 0, 0.5)', got %s", redAlpha.ToRgbString())
	}
}

// ─── Color Modifications ─────────────────────────────────────────────────────

func TestColorModifications(t *testing.T) {
	c := tc.New("red")

	t.Run("Lighten", func(t *testing.T) {
		l := c.Lighten(10)
		if l.ToHslString() == c.ToHslString() {
			t.Error("Lighten should change the color")
		}
	})

	t.Run("Darken", func(t *testing.T) {
		d := c.Darken(10)
		if d.ToHslString() == c.ToHslString() {
			t.Error("Darken should change the color")
		}
	})

	t.Run("Greyscale", func(t *testing.T) {
		g := c.Greyscale()
		hsl := g.ToHsl()
		if hsl.S > 0.001 {
			t.Errorf("Greyscale: expected S=0, got %v", hsl.S)
		}
	})

	t.Run("Complement", func(t *testing.T) {
		comp := c.Complement()
		hsl := comp.ToHsl()
		// Red is H=0, complement should be H≈180.
		if !approx(hsl.H, 180, 1.5) {
			t.Errorf("Complement of red: expected H≈180, got %v", hsl.H)
		}
	})
}

// ─── WCAG Readability ─────────────────────────────────────────────────────────

func TestWCAGReadability(t *testing.T) {
	white := tc.New("white")
	black := tc.New("black")

	if !white.IsLight() {
		t.Error("white should be IsLight")
	}
	if !black.IsDark() {
		t.Error("black should be IsDark")
	}

	ratio := tc.Readability(white, black)
	if !approx(ratio, 21, 0.1) {
		t.Errorf("White/black contrast ratio: expected 21, got %v", ratio)
	}

	if !tc.IsReadable(white, black, tc.ReadabilityOptions{Level: "AA", Size: "small"}) {
		t.Error("white/black should meet AA small readability")
	}
	if !tc.IsReadable(white, black, tc.ReadabilityOptions{Level: "AAA", Size: "small"}) {
		t.Error("white/black should meet AAA small readability")
	}
}

// ─── MostReadable ─────────────────────────────────────────────────────────────

func TestMostReadable(t *testing.T) {
	white := tc.New("white")
	black := tc.New("black")
	yellow := tc.New("yellow")

	// Against white background, black should be more readable than yellow.
	best := tc.MostReadable(white, []*tc.Color{black, yellow})
	if best.ToHexString() != "#000000" {
		t.Errorf("MostReadable(white, [black, yellow]): expected black, got %s", best.ToHexString())
	}
}

// ─── ToName ───────────────────────────────────────────────────────────────────

func TestToName(t *testing.T) {
	// From test.js line 694-695: tinycolor("#f00").toName() === "red"
	name, ok := tc.New("#f00").ToName()
	if !ok || name != "red" {
		t.Errorf("ToName(#f00): expected (\"red\", true), got (%q, %v)", name, ok)
	}

	// Non-named color returns false (test.js line 695)
	_, ok = tc.New("#fa0a0a").ToName()
	if ok {
		t.Errorf("ToName(#fa0a0a): expected false, got true")
	}

	// Transparent (a=0) returns "transparent" (test.js line 883-887)
	name, ok = tc.New(map[string]interface{}{"r": 255, "g": 20, "b": 10, "a": 0}).ToName()
	if !ok || name != "transparent" {
		t.Errorf("ToName(a=0): expected (\"transparent\", true), got (%q, %v)", name, ok)
	}

	// Semi-transparent returns false (test.js line 812-815)
	_, ok = tc.New("rgba(255, 0, 0, 0.5)").ToName()
	if ok {
		t.Errorf("ToName(semi-transparent): expected false, got true")
	}

	// toString on named color should return the name, not hex
	got := tc.New("red").ToString()
	if got != "red" {
		t.Errorf("ToString() on named color: expected \"red\", got %q", got)
	}
}

// ─── ToFilter (IE Legacy Gradient) ───────────────────────────────────────────

func TestToFilter(t *testing.T) {
	// Derived from test.js "Filters" Deno.test block.

	got := tc.New("red").ToFilter()
	want := "progid:DXImageTransform.Microsoft.gradient(startColorstr=#ffff0000,endColorstr=#ffff0000)"
	if got != want {
		t.Errorf("ToFilter(red): got %q, want %q", got, want)
	}

	got = tc.New("red").ToFilter(tc.New("blue"))
	want = "progid:DXImageTransform.Microsoft.gradient(startColorstr=#ffff0000,endColorstr=#ff0000ff)"
	if got != want {
		t.Errorf("ToFilter(red,blue): got %q, want %q", got, want)
	}

	got = tc.New("transparent").ToFilter()
	want = "progid:DXImageTransform.Microsoft.gradient(startColorstr=#00000000,endColorstr=#00000000)"
	if got != want {
		t.Errorf("ToFilter(transparent): got %q, want %q", got, want)
	}

	got = tc.New("transparent").ToFilter(tc.New("red"))
	want = "progid:DXImageTransform.Microsoft.gradient(startColorstr=#00000000,endColorstr=#ffff0000)"
	if got != want {
		t.Errorf("ToFilter(transparent,red): got %q, want %q", got, want)
	}

	got = tc.New("#f0f0f0dd").ToFilter()
	want = "progid:DXImageTransform.Microsoft.gradient(startColorstr=#ddf0f0f0,endColorstr=#ddf0f0f0)"
	if got != want {
		t.Errorf("ToFilter(#f0f0f0dd): got %q, want %q", got, want)
	}

	got = tc.New("rgba(0, 0, 255, .5)").ToFilter()
	want = "progid:DXImageTransform.Microsoft.gradient(startColorstr=#800000ff,endColorstr=#800000ff)"
	if got != want {
		t.Errorf("ToFilter(rgba(0,0,255,.5)): got %q, want %q", got, want)
	}
}
