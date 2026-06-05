package utility

import "testing"

func TestLookupDensity_BakingPowder(t *testing.T) {
	d := LookupDensity("baking powder")
	if d == 0 {
		t.Error("baking powder density should be non-zero")
	}
}

func TestLookupDensity_BakingPowder_OneTeaspoonIsFourGrams(t *testing.T) {
	// 1 tsp = 5 mL; 1 tsp baking powder is conventionally 4 g
	// so density should be 4/5 = 0.80 g/mL
	d := LookupDensity("baking powder")
	tspML := 5.0
	grams := tspML * d
	if grams < 3.8 || grams > 4.2 {
		t.Errorf("1 tsp baking powder should be ~4 g, got %.2f g (density=%.3f)", grams, d)
	}
}
