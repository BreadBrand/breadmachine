package parser

import "testing"

func TestIsUnsupportedYeast(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		// unsupported
		{"active dry yeast", true},
		{"fresh yeast", true},
		{"compressed yeast", true},
		{"cake yeast", true},
		// supported sourdough variants
		{"sourdough starter", false},
		{"sourdough discard", false},
		{"ripe starter", false},
		{"levain", false},
		{"poolish", false},
		{"biga", false},
		// supported dry yeast
		{"instant dry yeast", false},
		{"instant yeast", false},
		{"rapid rise yeast", false},
		// critical: "active starter" must NOT match "active dry yeast"
		{"active starter", false},
		{"bubbly active sourdough starter", false},
		// non-yeast ingredients
		{"bread flour", false},
		{"water", false},
		{"salt", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsUnsupportedYeast(c.name); got != c.want {
				t.Errorf("IsUnsupportedYeast(%q) = %v, want %v", c.name, got, c.want)
			}
		})
	}
}
