package handlers

import "github.com/BreadBrand/breadmachine/models"

func isMassUnit(unit string) bool {
	switch unit {
	case "g", "kg", "oz", "lb":
		return true
	}
	return false
}

func isVolumeUnit(unit string) bool {
	switch unit {
	case "tsp", "tbsp", "cup", "ml", "l", "fl oz":
		return true
	}
	return false
}

// isGramDominant returns true when the recipe contains at least one mass-unit
// ingredient (g, kg, oz, lb). Any mass unit is enough to trigger conversion
// of volume measurements; all-volume recipes are left unchanged.
func isGramDominant(ingredients []models.Ingredient) bool {
	for _, ing := range ingredients {
		if isMassUnit(ing.Unit) {
			return true
		}
	}
	return false
}

// convertToGrams rewrites volume-unit ingredients to grams in gram-dominant
// recipes. Ingredients with an unknown density (Grams == 0) are left unchanged.
func convertToGrams(ingredients []models.Ingredient) []models.Ingredient {
	if !isGramDominant(ingredients) {
		return ingredients
	}
	for i := range ingredients {
		ing := &ingredients[i]
		if !isVolumeUnit(ing.Unit) || ing.Grams <= 0 {
			continue
		}
		ing.Quantity = ing.Grams
		ing.Unit = "g"
	}
	return ingredients
}
