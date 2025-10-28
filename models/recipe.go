package models

import (
	"strings"
	"time"
)

// Phase represents the stage in which an ingredient is used.
type Phase string

const (
	PhaseDough    Phase = "dough"
	PhaseScald    Phase = "scald"
	PhaseSoak     Phase = "soak"
	PhaseAutolyse Phase = "autolyse"
)

// Ingredient represents a single recipe component with baker's percentage and metadata.
type Ingredient struct {
	ID              string  `json:"id" firestore:"id"`
	IngredientName  string  `json:"ingredientName" firestore:"ingredientName"`
	BakerPercentage float64 `json:"bakerPercentage" firestore:"bakerPercentage"`
	Quantity        float64 `json:"quantity,omitempty" firestore:"quantity,omitempty"`
	Unit            string  `json:"unit" firestore:"unit"`
	Grams           float64 `json:"grams" firestore:"grams"`
	Phase           Phase   `json:"phase" firestore:"phase"`
	DensityGPerMl   float64 `json:"densityGPerMl,omitempty" firestore:"densityGPerMl,omitempty"`
}

// Meta holds auxiliary information about a recipe.
type Meta struct {
	YieldGrams float64   `json:"yieldGrams" firestore:"yieldGrams"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" firestore:"updatedAt"`
	Tags       []string  `json:"tags,omitempty" firestore:"tags,omitempty"`
}

// Recipe is the top‑level model for storing and scaling bread recipes.
type Recipe struct {
	ID           string       `json:"id" firestore:"id"`
	Title        string       `json:"title" firestore:"title"`
	Description  string       `json:"description" firestore:"description"`
	Instructions []string     `json:"instructions" firestore:"instructions"`
	Ingredients  []Ingredient `json:"ingredients" firestore:"ingredients"`
	Meta         Meta         `json:"meta" firestore:"meta"`
	UserID       string       `json:"userId,omitempty" firestore:"userId,omitempty"`
}

// CalculateBakerPercentages inspects r.Ingredients, sums all flour quantities,
// and populates each Ingredient.BakerPercentage accordingly.
func (r *Recipe) CalculateBakerPercentages() {
	var totalFlour float64

	for _, ing := range r.Ingredients {
		if (ing.Phase == PhaseDough || ing.Phase == PhaseScald) &&
			strings.Contains(strings.ToLower(ing.IngredientName), "flour") {
			totalFlour += ing.Grams
		}
	}

	for i := range r.Ingredients {
		if totalFlour > 0 && (r.Ingredients[i].Phase == PhaseDough || r.Ingredients[i].Phase == PhaseScald) {
			r.Ingredients[i].BakerPercentage = (r.Ingredients[i].Grams / totalFlour) * 100
		} else {
			r.Ingredients[i].BakerPercentage = 0
		}
	}
}
