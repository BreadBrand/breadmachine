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
	YeastTypeNone      YeastType = "none"
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
	PrepTime       int       `json:"prepTime,omitempty" firestore:"prepTime,omitempty"`
	CookTime       int       `json:"cookTime,omitempty" firestore:"cookTime,omitempty"`
	AdditionalTime int       `json:"additionalTime,omitempty" firestore:"additionalTime,omitempty"`
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

// baseIngredientKeywords lists ingredient name substrings that serve as the
// baker's-math base (100%) in specialty breads. "flour" covers the common case;
// the rest handle grain-free and legume-based recipes.
var baseIngredientKeywords = []string{
	"flour",
	"lentil",
	"oat",
	"cauliflower",
	"chickpea",
	"tapioca",
}

func isBaseIngredient(name string) bool {
	lower := strings.ToLower(name)
	for _, kw := range baseIngredientKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// CalculateBakerPercentages computes baker's percentages for DoughIngredients only.
// Base weight is the sum of any ingredient whose name matches a known base keyword
// (flour, lentil, oat, cauliflower, chickpea, tapioca). Recipes with no matching
// base ingredient get zero percentages.
func (r *Recipe) CalculateBakerPercentages() {
	isFlourPhase := func(p Phase) bool {
		switch strings.ToLower(string(p)) {
		case "dough", "scald", "tangzhong", "yudane", "starter build", "levain", "final dough", "":
			return true
		default:
			return false
		}
	}

	var totalBase float64
	for _, ing := range r.DoughIngredients {
		if isFlourPhase(ing.Phase) && isBaseIngredient(ing.IngredientName) {
			totalBase += ing.Grams
		}
	}

	for i := range r.DoughIngredients {
		if totalBase > 0 && isFlourPhase(r.DoughIngredients[i].Phase) {
			r.DoughIngredients[i].BakerPercentage = (r.DoughIngredients[i].Grams / totalBase) * 100
		} else {
			r.DoughIngredients[i].BakerPercentage = 0
		}
	}
}
