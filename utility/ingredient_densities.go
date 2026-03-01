package utility

import (
	"regexp"
	"strings"
)

type densityEntry struct {
    keyword string
    density float64
}

// Order matters - more specific entries first
var ingredientDensities = []densityEntry{
    {"whole wheat flour", 0.593},
		{"all purpose flour", 0.53},
    {"bread flour", 0.57},
    {"strong flour", 0.57},
    {"barley flour", 0.59},
    {"white flour", 0.53},
    {"ap flour", 0.53},
    {"rye flour", 0.62},
    {"powdered sugar", 0.56},
    {"brown sugar", 0.93},
    {"active dry yeast", 0.43},
    {"instant yeast", 0.43},
    {"dried yeast", 0.43},
    {"sea salt", 1.18},
    {"olive oil", 0.908},
    {"avocado oil", 0.913},
    {"sourdough", 1.0},
    {"baking soda", 0.88},
    // single word matches last
    {"barley", 0.59},
    {"flour", 0.53},
    {"water", 1.0},
    {"milk", 1.03},
    {"butter", 0.91},
    {"sugar", 0.845},
    {"salt", 1.217},
    {"beer", 1.01},
    {"egg", 1.03},
    {"vanilla", 0.879},
    {"starter", 1.0},
}

func LookupDensity(ingredientName string) float64 {
    normalized := strings.ToLower(ingredientName)
    normalized = regexp.MustCompile(`\(.*?\)`).ReplaceAllString(normalized, "")
    normalized = strings.TrimSpace(normalized)

    for _, entry := range ingredientDensities {
        if strings.Contains(normalized, entry.keyword) {
            return entry.density
        }
    }
    return 0
}
