package tinycolor

import (
	"math"
	"testing"
)

func assertColor(t *testing.T, input interface{}, expectedR, expectedG, expectedB float64, expectedA float64, expectedFormat Format, expectedOk bool) {
	t.Helper()
	c := New(input)
	if math.Abs(c.r-expectedR) > 0.001 || math.Abs(c.g-expectedG) > 0.001 || math.Abs(c.b-expectedB) > 0.001 {
		t.Errorf("For input %v, expected RGB(%v, %v, %v) but got RGB(%v, %v, %v)", input, expectedR, expectedG, expectedB, c.r, c.g, c.b)
	}
	if math.Abs(c.a-expectedA) > 0.001 {
		t.Errorf("For input %v, expected alpha %v but got %v", input, expectedA, c.a)
	}
	if c.format != expectedFormat {
		t.Errorf("For input %v, expected format %v but got %v", input, expectedFormat, c.format)
	}
	if c.ok != expectedOk {
		t.Errorf("For input %v, expected ok %v but got %v", input, expectedOk, c.ok)
	}
}

func TestParser_Hex(t *testing.T) {
	assertColor(t, "#f00", 255, 0, 0, 1, FormatHex, true)
	assertColor(t, "f00", 255, 0, 0, 1, FormatHex, true)
	assertColor(t, "#ff0000", 255, 0, 0, 1, FormatHex, true)
	assertColor(t, "ff0000", 255, 0, 0, 1, FormatHex, true)
	assertColor(t, "#f00f", 255, 0, 0, 1, FormatHex8, true)
	assertColor(t, "f00f", 255, 0, 0, 1, FormatHex8, true)
	assertColor(t, "#ff0000ff", 255, 0, 0, 1, FormatHex8, true)
	assertColor(t, "ff0000ff", 255, 0, 0, 1, FormatHex8, true)
	assertColor(t, "#ff000000", 255, 0, 0, 0, FormatHex8, true)
}

func TestParser_Named(t *testing.T) {
	assertColor(t, "red", 255, 0, 0, 1, FormatName, true)
	assertColor(t, "aliceblue", 240, 248, 255, 1, FormatName, true)
	assertColor(t, "transparent", 0, 0, 0, 0, FormatName, true)
	assertColor(t, "RED", 255, 0, 0, 1, FormatName, true)
}

func TestParser_RGB(t *testing.T) {
	assertColor(t, "rgb(255, 0, 0)", 255, 0, 0, 1, FormatRgb, true)
	assertColor(t, "rgb 255 0 0", 255, 0, 0, 1, FormatRgb, true)
	assertColor(t, "rgba(255, 0, 0, 0.5)", 255, 0, 0, 0.5, FormatRgb, true)
	assertColor(t, "rgba 255 0 0 0.5", 255, 0, 0, 0.5, FormatRgb, true)
	assertColor(t, "rgb(100%, 0%, 0%)", 255, 0, 0, 1, FormatPrgb, true)
	assertColor(t, "rgba(100%, 0%, 0%, 0.5)", 255, 0, 0, 0.5, FormatPrgb, true)
}

func TestParser_HSL_HSV(t *testing.T) {
	assertColor(t, "hsl(0, 100%, 50%)", 255, 0, 0, 1, FormatHsl, true)
	assertColor(t, "hsla(0, 100%, 50%, 0.5)", 255, 0, 0, 0.5, FormatHsl, true)
	assertColor(t, "hsv(0, 100%, 100%)", 255, 0, 0, 1, FormatHsv, true)
	assertColor(t, "hsva(0, 100%, 100%, 0.5)", 255, 0, 0, 0.5, FormatHsv, true)
}

func TestParser_Objects(t *testing.T) {
	assertColor(t, map[string]interface{}{"r": 255, "g": 0, "b": 0}, 255, 0, 0, 1, FormatRgb, true)
	assertColor(t, map[string]interface{}{"r": "100%", "g": 0, "b": 0}, 255, 0, 0, 1, FormatPrgb, true)
	assertColor(t, map[string]interface{}{"h": 0, "s": 100, "l": 50}, 255, 0, 0, 1, FormatHsl, true)
	assertColor(t, map[string]interface{}{"h": 0, "s": 100, "v": 100}, 255, 0, 0, 1, FormatHsv, true)
}

func TestParser_Invalid(t *testing.T) {
	assertColor(t, "not a color", 0, 0, 0, 1, "", false)
	assertColor(t, map[string]interface{}{"unknown": 255}, 0, 0, 0, 1, "", false)
}

func TestFromRatio(t *testing.T) {
	c := FromRatio(map[string]interface{}{"r": 1.0, "g": 0.0, "b": 0.0})
	if c.r != 255 || c.g != 0 || c.b != 0 {
		t.Errorf("Expected FromRatio to correctly scale RGB ratios, got R:%v G:%v B:%v", c.r, c.g, c.b)
	}
	c2 := FromRatio(map[string]interface{}{"h": 1.0, "s": 1.0, "l": 0.5})
	if c2.r != 255 || c2.g != 0 || c2.b != 0 {
		t.Errorf("Expected FromRatio to correctly scale HSL ratios, got R:%v G:%v B:%v", c2.r, c2.g, c2.b)
	}
}
