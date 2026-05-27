package parser

import "strings"

var supportedYeastPhrases = []string{
	"sourdough starter",
	"sourdough discard",
	"active starter",
	"ripe starter",
	"levain",
	"poolish",
	"biga",
	"instant dry yeast",
	"instant yeast",
	"rapid rise yeast",
}

var unsupportedYeastPhrases = []string{
	"active dry yeast",
	"fresh yeast",
	"compressed yeast",
	"cake yeast",
}

// IsUnsupportedYeast returns true if name contains an unsupported leavener
// and does not contain a supported yeast or pre-ferment phrase.
// Supported phrases are checked first to prevent "active starter" from
// matching the "active" in "active dry yeast".
func IsUnsupportedYeast(name string) bool {
	lower := strings.ToLower(name)
	for _, phrase := range supportedYeastPhrases {
		if strings.Contains(lower, phrase) {
			return false
		}
	}
	for _, phrase := range unsupportedYeastPhrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}
