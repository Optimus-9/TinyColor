package tinycolor

// RGB represents a color in the RGB color space.
// R, G, B values are in [0, 255]; A is in [0.0, 1.0].
type RGB struct {
	R, G, B, A float64
}

// HSL represents a color in the HSL color space.
// H is in [0, 360]; S and L are in [0, 1]; A is in [0.0, 1.0].
type HSL struct {
	H, S, L, A float64
}

// HSV represents a color in the HSV color space.
// H is in [0, 360]; S and V are in [0, 1]; A is in [0.0, 1.0].
type HSV struct {
	H, S, V, A float64
}

// ReadabilityOptions configures WCAG readability checking.
type ReadabilityOptions struct {
	Level                 string // "AA" or "AAA"
	Size                  string // "small" or "large"
	IncludeFallbackColors bool
}
