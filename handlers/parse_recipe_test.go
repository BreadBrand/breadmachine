package handlers

import (
	"testing"

	"github.com/RedBrand88/breadmachine/models"
	"github.com/RedBrand88/breadmachine/utility"
)

func TestParseRecipeText(t *testing.T) {
	input := `
		# Title
		Sour Dough Discard Cinnamon Rolls

		# Description
		Soft, rich cinnamon rolls made with sourdough discard and a tangzhong roux.

		# Ingredients
		[tangzhong]
		25g bread flour (about 3 Tbsp)
		½ cup (125g) water

		[dough]
		475g bread flour
		1 Tbsp sugar
		2 Tbsp starter
		1 ½ cups milk
		1 1/3 cups barley flour


		# Instructions
		1. Make the tangzhong: whisk 125g water and 25g flour together.
		2. Mix wet and dry ingredients.
		3. Bake at 325°F for 35-40 minutes.
	`

	recipe := ParseRecipeText(input)
	recipe.CalculateBakerPercentages()

	// Title
	if recipe.Title != "Sour Dough Discard Cinnamon Rolls" {
		t.Errorf("expected title %q, got %q", "Sour Dough Discard Cinnamon Rolls", recipe.Title)
	}

	// Description
	if recipe.Description == "" {
		t.Errorf("expected description to be parsed, got empty string")
	}

	expectedIngredients := []models.Ingredient{
		{IngredientName: "bread flour (about 3 Tbsp)", Quantity: 25, Unit: "g", Grams: 25, Phase: "tangzhong", BakerPercentage: 3.05, DensityGPerMl: 0.57},
		{IngredientName: "(125g) water", Quantity: 0.5, Unit: "cup", Grams: 120, Phase: "tangzhong", BakerPercentage: 14.63, DensityGPerMl: 1.0},
		{IngredientName: "bread flour", Quantity: 475, Unit: "g", Grams: 475, Phase: "dough", BakerPercentage: 57.93, DensityGPerMl: 0.57},
		{IngredientName: "sugar", Quantity: 1, Unit: "Tbsp", Grams: 15, Phase: "dough", BakerPercentage: 1.83, DensityGPerMl: 0.845},
		{IngredientName: "starter", Quantity: 2, Unit: "Tbsp", Grams: 30, Phase: "dough", BakerPercentage: 3.66, DensityGPerMl: 1.0},
		{IngredientName: "milk", Quantity: 1.5, Unit: "cups", Grams: 360, Phase: "dough", BakerPercentage: 43.90, DensityGPerMl: 1.03},
		{IngredientName: "barley flour", Quantity: 1.33, Unit: "cups", Grams: 320, Phase: "dough", BakerPercentage: 39.02, DensityGPerMl: 0.59},
	}
	// Ingredient count
	if len(recipe.Ingredients) != len(expectedIngredients) {
		t.Errorf("expected %d ingredients, got %d", len(expectedIngredients), len(recipe.Ingredients))
	}

	for i, exp := range expectedIngredients {
		got := recipe.Ingredients[i]

		if got.Phase != exp.Phase {
			t.Errorf("ingredient %s phase mismatch: expected %q, got %q", exp.IngredientName, exp.Phase, got.Phase)
		}

		if got.Unit != exp.Unit {
			t.Errorf("ingredient %s unit mismatch: expected %q, got %q", exp.IngredientName, exp.Unit, got.Unit)
		}

		if got.IngredientName != exp.IngredientName {
			t.Errorf("ingredient %s name mismatch: expected %q, got %q", exp.IngredientName, exp.IngredientName, got.IngredientName)
		}

		diff := got.Quantity - exp.Quantity
		if diff < -0.01 || diff > 0.01 {
			t.Errorf("ingredient %s quantity mismatch: expected %.3f, got %.3f", exp.IngredientName, exp.Quantity, got.Quantity)
		}

		grams_diff := got.Grams - exp.Grams
		if grams_diff < -0.01 || grams_diff > 0.01 {
			t.Errorf("ingredient %s grams mismatch: expected %.2f, got %.2f", exp.IngredientName, exp.Grams, got.Grams)
		}

		percent_diff := got.BakerPercentage - exp.BakerPercentage
		if percent_diff < -0.01 || percent_diff > 0.01 {
			t.Errorf("ingredient %s bakers percent mismatch: expected %.2f, got %.2f", exp.IngredientName, exp.BakerPercentage, got.BakerPercentage)
		}

		density_diff := got.DensityGPerMl - exp.DensityGPerMl
		if density_diff < -0.001 || density_diff > 0.001 {
			t.Errorf("ingredient %s density mismatch: expected %.3f, got %.3f", exp.IngredientName, exp.DensityGPerMl, got.DensityGPerMl)
		}
	}

	expectedInstructions := []string{
		"Make the tangzhong: whisk 125g water and 25g flour together.",
		"Mix wet and dry ingredients.",
		"Bake at 325°F for 35-40 minutes.",
	}

	// Instructions
	if len(recipe.Instructions) != len(expectedInstructions) {
		t.Errorf("expected %d instructions, got %d", len(expectedInstructions), len(recipe.Instructions))
	}

	for i, exp := range expectedInstructions {
		if recipe.Instructions[i] != exp {
			t.Errorf("instruction %d mismatch:\nexpected: %q\ngot: %q", i, exp, recipe.Instructions[i])
		}
	}
}

func TestLookupDensity(t *testing.T) {
	cases := []struct {
		name     string
		expected float64
	}{
		{"bread flour", 0.57},
		{"AP flour (11.7% protein)", 0.53}, // parenthetical stripped
		{"strong white bread flour", 0.57},
		{"Fine Sea Salt", 1.18},
		{"Instant Yeast", 0.43},
		{"olive oil", 0.908},
		{"beer", 1.01},
		{"sourdough starter", 1.0},
		{"barley flour", 0.59},
		{"baking soda", 0.88},
	}

	for _, tc := range cases {
		got := utility.LookupDensity(tc.name)
		diff := got - tc.expected
		if diff < -0.001 || diff > 0.001 {
			t.Errorf("LookupDensity(%q): expected %.3f, got %.3f", tc.name, tc.expected, got)
		}
	}
}

func TestParseIngredientLineNoMatch(t *testing.T) {
    cases := []struct {
        line  string
        phase models.Phase
    }{
        {"some text with no number", "dough"},
        {"", "dough"},
        {"just words here", "dough"},
    }

    for _, tc := range cases {
        // should not panic
        got := parseIngredientLine(tc.line, tc.phase)
        if got.Phase != tc.phase {
            t.Errorf("expected phase %q got %q", tc.phase, got.Phase)
        }
    }
}
