package models

import (
	"strings"
	"time"
)

type Phase string

const (
	PhaseDough      Phase = "dough"
	PhaseScald      Phase = "scald"
	PhaseTangzhong  Phase = "tangzhong"
	PhaseYudane     Phase = "yudane"
	PhaseSoak       Phase = "soak"
	PhaseAutolyse   Phase = "autolyse"
)

type YeastType string

const (
	YeastTypeDry       YeastType = "dry"
	YeastTypeSourdough YeastType = "sourdough"
)

type Ingredient struct {
	ID              string  `json:"id" firestore:"id"`
	IngredientName  string  `json:"ingredientName" firestore:"ingredientName"`
	BakerPercentage float64 `json:"bakerPercentage" firestore:"bakerPercentage"`
	Quantity        float64 `json:"quantity,omitempty" firestore:"quantity,omitempty"`
	Unit            string  `json:"unit" firestore:"unit"`
	Grams           float64 `json:"grams" firestore:"grams"`
	Phase           Phase   `json:"phase,omitempty" firestore:"phase,omitempty"`
	DensityGPerMl   float64 `json:"densityGPerMl,omitempty" firestore:"densityGPerMl,omitempty"`
}

type Meta struct {
	YieldGrams     float64   `json:"yieldGrams" firestore:"yieldGrams"`
	CreatedAt      time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" firestore:"updatedAt"`
	Tags           []string  `json:"tags,omitempty" firestore:"tags,omitempty"`
	Servings       string    `json:"servings,omitempty" firestore:"servings,omitempty"`
	PrepTime       string    `json:"prepTime,omitempty" firestore:"prepTime,omitempty"`
	CookTime       string    `json:"cookTime,omitempty" firestore:"cookTime,omitempty"`
	AdditionalTime string    `json:"additionalTime,omitempty" firestore:"additionalTime,omitempty"`
}

type Recipe struct {
	ID               string       `json:"id" firestore:"id"`
	Title            string       `json:"title" firestore:"title"`
	Description      string       `json:"description" firestore:"description"`
	Instructions     []string     `json:"instructions" firestore:"instructions"`
	DoughIngredients []Ingredient `json:"doughIngredients" firestore:"doughIngredients"`
	OtherIngredients []Ingredient `json:"otherIngredients" firestore:"otherIngredients"`
	Meta             Meta         `json:"meta" firestore:"meta"`
	UserID           string       `json:"userId,omitempty" firestore:"userId,omitempty"`
	YeastType        YeastType    `json:"yeastType,omitempty" firestore:"yeastType,omitempty"`
}

// CalculateBakerPercentages computes baker's percentages for DoughIngredients only.
func (r *Recipe) CalculateBakerPercentages() {
	var totalFlour float64

	isFlourPhase := func(p Phase) bool {
		switch strings.ToLower(string(p)) {
		case "dough", "scald", "tangzhong", "yudane", "starter build", "levain", "final dough", "":
			return true
		default:
			return false
		}
	}

	for _, ing := range r.DoughIngredients {
		if isFlourPhase(ing.Phase) && strings.Contains(strings.ToLower(ing.IngredientName), "flour") {
			totalFlour += ing.Grams
		}
	}

	for i := range r.DoughIngredients {
		if totalFlour > 0 && isFlourPhase(r.DoughIngredients[i].Phase) {
			r.DoughIngredients[i].BakerPercentage = (r.DoughIngredients[i].Grams / totalFlour) * 100
		} else {
			r.DoughIngredients[i].BakerPercentage = 0
		}
	}
}
