package parser

import "strings"

// KnownUnits is the canonical list of recognised unit strings.
// Matching is always case-insensitive. Multi-word units (e.g. "fl oz") must be
// checked before single-word units to prevent partial matches.
var KnownUnits = []string{
	// multi-word first
	"fl oz", "fluid ounce", "fluid ounces",
	// weight
	"g", "gram", "grams",
	"kg", "kilogram", "kilograms",
	"oz", "ounce", "ounces",
	"lb", "pound", "pounds",
	// volume
	"ml", "milliliter", "millilitre", "milliliters", "millilitres",
	"l", "liter", "litre", "liters", "litres",
	"tsp", "teaspoon", "teaspoons",
	"tbsp", "tablespoon", "tablespoons",
	"cup", "cups",
	// approximate
	"pinch", "handful", "dash", "smidge", "sprig", "clove", "slice", "piece",
}

// MatchUnit returns the canonical unit string if the token matches a known unit,
// or empty string if no match. Matching is case-insensitive.
func MatchUnit(token string) string {
	lower := strings.ToLower(strings.TrimSpace(token))
	for _, u := range KnownUnits {
		if lower == u {
			return u
		}
	}
	return ""
}
