package handlers

import (
	"testing"

	"github.com/RedBrand88/breadmachine/models"
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
		1 1/3 cups barley


		# Instructions
		1. Make the tangzhong: whisk 125g water and 25g flour together.
		2. Mix wet and dry ingredients.
		3. Bake at 325°F for 35-40 minutes.
	`

	recipe := ParseRecipeText(input)

	// Title
	if recipe.Title != "Sour Dough Discard Cinnamon Rolls" {
		t.Errorf("expected title %q, got %q", "Sour Dough Discard Cinnamon Rolls", recipe.Title)
	}

	// Description
	if recipe.Description == "" {
		t.Errorf("expected description to be parsed, got empty string")
	}

	expectedIngredients := []models.Ingredient{
		{IngredientName: "bread flour (about 3 Tbsp)", Quantity: 25, Unit: "g", Grams: 25, Phase: "tangzhong"},
		{IngredientName: "(125g) water", Quantity: 0.5, Unit: "cup", Grams: 120, Phase: "tangzhong"},
		{IngredientName: "bread flour", Quantity: 475, Unit: "g", Grams: 475, Phase: "dough"},
		{IngredientName: "sugar", Quantity: 1, Unit: "Tbsp", Grams: 15, Phase: "dough"},
		{IngredientName: "starter", Quantity: 2, Unit: "Tbsp", Grams: 30, Phase: "dough"},
		{IngredientName: "milk", Quantity: 1.5, Unit: "cups", Grams: 360, Phase: "dough"},
		{IngredientName: "barley", Quantity: 1.33, Unit: "cups", Grams: 320, Phase: "dough"},
	}
	// Ingredient count
	if len(recipe.Ingredients) != len(expectedIngredients) {
		t.Errorf("expected %d ingredients, got %d", len(expectedIngredients), len(recipe.Ingredients))
	}

	for i, exp := range expectedIngredients {
		got := recipe.Ingredients[i]

		if got.Phase != exp.Phase {
			t.Errorf("ingredient %d phase mismatch: expected %q, got %q", i, exp.Phase, got.Phase)
		}

		if got.Unit != exp.Unit {
			t.Errorf("ingredient %d unit mismatch: expected %q, got %q", i, exp.Unit, got.Unit)
		}

		if got.IngredientName != exp.IngredientName {
			t.Errorf("ingredient %d name mismatch: expected %q, got %q", i, exp.IngredientName, got.IngredientName)
		}

		diff := got.Quantity - exp.Quantity
		if diff < -0.01 || diff > 0.01 {
			t.Errorf("ingredient %d quantity mismatch: expected %.3f, got %.3f", i, exp.Quantity, got.Quantity)
		}

		grams_diff := got.Grams - exp.Grams
		if grams_diff < -0.01 || grams_diff > 0.01 {
			t.Errorf("ingredient %d grams mismatch: expected %.2f, got %.2f", i, exp.Grams, got.Grams)
		}
	}

	expectedInstructions := []string{
		"1. Make the tangzhong: whisk 125g water and 25g flour together.",
		"2. Mix wet and dry ingredients.",
		"3. Bake at 325°F for 35-40 minutes.",
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
