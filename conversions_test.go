package tinycolor

import (
	"math"
	"testing"
)

// Helper to compare floats with epsilon.
func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < 1e-3
}

func TestRgbToHsl(t *testing.T) {
	// White
	hsl := rgbToHsl(255, 255, 255)
	if !floatEquals(hsl.H, 0) || !floatEquals(hsl.S, 0) || !floatEquals(hsl.L, 1.0) {
		t.Errorf("Expected 0, 0, 1.0 for White, got %f, %f, %f", hsl.H, hsl.S, hsl.L)
	}

	// Red
	hsl = rgbToHsl(255, 0, 0)
	if !floatEquals(hsl.H, 0) || !floatEquals(hsl.S, 1.0) || !floatEquals(hsl.L, 0.5) {
		t.Errorf("Expected 0, 1.0, 0.5 for Red, got %f, %f, %f", hsl.H, hsl.S, hsl.L)
	}

	// Green
	hsl = rgbToHsl(0, 255, 0)
	if !floatEquals(hsl.H, 120) || !floatEquals(hsl.S, 1.0) || !floatEquals(hsl.L, 0.5) {
		t.Errorf("Expected 120, 1.0, 0.5 for Green, got %f, %f, %f", hsl.H, hsl.S, hsl.L)
	}
}

func TestRgbToHsv(t *testing.T) {
	// Red
	hsv := rgbToHsv(255, 0, 0)
	if !floatEquals(hsv.H, 0) || !floatEquals(hsv.S, 1.0) || !floatEquals(hsv.V, 1.0) {
		t.Errorf("Expected 0, 1.0, 1.0 for Red, got %f, %f, %f", hsv.H, hsv.S, hsv.V)
	}
}

func TestHslToRgb(t *testing.T) {
	// Red (0, 100%, 50%) -> (255, 0, 0)
	rgb := HslToRgb(0, 100, 50)
	if !floatEquals(rgb.R, 255) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (255, 0, 0) for HSL Red, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Green (120, 100%, 50%) -> (0, 255, 0)
	rgb = HslToRgb(120, 100, 50)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (0, 255, 0) for HSL Green, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Blue (240, 100%, 50%) -> (0, 0, 255)
	rgb = HslToRgb(240, 100, 50)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 255) {
		t.Errorf("Expected (0, 0, 255) for HSL Blue, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// White (0, 0%, 100%) -> (255, 255, 255)
	rgb = HslToRgb(0, 0, 100)
	if !floatEquals(rgb.R, 255) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 255) {
		t.Errorf("Expected (255, 255, 255) for HSL White, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Black (0, 0%, 0%) -> (0, 0, 0)
	rgb = HslToRgb(0, 0, 0)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (0, 0, 0) for HSL Black, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Yellow (60, 100%, 50%) -> (255, 255, 0)
	rgb = HslToRgb(60, 100, 50)
	if !floatEquals(rgb.R, 255) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (255, 255, 0) for HSL Yellow, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}
}

func TestHsvToRgb(t *testing.T) {
	// Red (0, 100%, 100%) -> (255, 0, 0)
	rgb := HsvToRgb(0, 100, 100)
	if !floatEquals(rgb.R, 255) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (255, 0, 0) for HSV Red, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Green (120, 100%, 100%) -> (0, 255, 0)
	rgb = HsvToRgb(120, 100, 100)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (0, 255, 0) for HSV Green, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Blue (240, 100%, 100%) -> (0, 0, 255)
	rgb = HsvToRgb(240, 100, 100)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 255) {
		t.Errorf("Expected (0, 0, 255) for HSV Blue, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// White (0, 0%, 100%) -> (255, 255, 255)
	rgb = HsvToRgb(0, 0, 100)
	if !floatEquals(rgb.R, 255) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 255) {
		t.Errorf("Expected (255, 255, 255) for HSV White, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Black (0, 0%, 0%) -> (0, 0, 0)
	rgb = HsvToRgb(0, 0, 0)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 0) || !floatEquals(rgb.B, 0) {
		t.Errorf("Expected (0, 0, 0) for HSV Black, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}

	// Cyan (180, 100%, 100%) -> (0, 255, 255)
	rgb = HsvToRgb(180, 100, 100)
	if !floatEquals(rgb.R, 0) || !floatEquals(rgb.G, 255) || !floatEquals(rgb.B, 255) {
		t.Errorf("Expected (0, 255, 255) for HSV Cyan, got (%f, %f, %f)", rgb.R, rgb.G, rgb.B)
	}
}

func TestRgbToHex(t *testing.T) {
	hex := rgbToHex(255, 0, 0, false)
	if hex != "ff0000" {
		t.Errorf("Expected ff0000, got %s", hex)
	}

	hex3 := rgbToHex(255, 0, 0, true)
	if hex3 != "f00" {
		t.Errorf("Expected f00, got %s", hex3)
	}
}

func TestToString(t *testing.T) {
	c := &Color{r: 255, g: 0, b: 0, a: 1, roundA: 1, format: FormatHex, ok: true}

	if c.ToHexString() != "#ff0000" {
		t.Errorf("Expected #ff0000, got %s", c.ToHexString())
	}

	if c.ToHexString(true) != "#f00" {
		t.Errorf("Expected #f00, got %s", c.ToHexString(true))
	}

	cInvalid := &Color{ok: false}
	if cInvalid.ToHexString() != "#000000" {
		t.Errorf("Expected #000000 for invalid, got %s", cInvalid.ToHexString())
	}
}

func TestWCAG(t *testing.T) {
	c1 := &Color{r: 255, g: 255, b: 255, a: 1, ok: true}
	c2 := &Color{r: 0, g: 0, b: 0, a: 1, ok: true}

	if !c1.IsLight() {
		t.Errorf("White should be light")
	}

	if !c2.IsDark() {
		t.Errorf("Black should be dark")
	}

	readability := Readability(c1, c2)
	if readability != 21 {
		t.Errorf("Expected contrast between black and white to be 21, got %f", readability)
	}

	if !IsReadable(c1, c2, ReadabilityOptions{Level: "AA", Size: "small"}) {
		t.Errorf("Black and white should be readable")
	}
}
