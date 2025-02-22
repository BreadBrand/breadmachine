package models

import "strings"

// Ingredient struct represents a recipe ingredient
type Ingredient struct {
	Name     string  `json:"ingredientName" firestore:"ingredientName"`
	Quantity float64 `json:"quantity" firestore:"quantity"`
	Unit     string  `json:"unit" firestore:"unit"`
}

// Percentage stuct represents ingredient baker's percentages
type Percentage struct {
	Name    string  `json:"ingredientName" firestore:"ingredientName"`
	Percent float64 `json:"percent" firestore:"percent"`
}

// Recipe stuct represents a full recipe
type Recipe struct {
	Title        string       `json:"title" firestore:"title"`
	Description  string       `json:"description" firestore:"description"`
	Ingredients  []Ingredient `json:"ingredients" firestore:"ingredients"`
	Instructions []string     `json:"instructions" firestore:"instructions"`
	Percentages  []Percentage `json:"percentages" firestore:"percentages"`
}

//CalculateBakersPercentages computes percentages based on flour weight
func (r *Recipe) CalculateBakersPercentages() {
	var totalFlour float64

	//Sum up all flour-based ingredients
	for _, ingredient := range r.Ingredients {
		if strings.Contains(strings.ToLower(ingredient.Name), "flour") {
			totalFlour += ingredient.Quantity
		}
	}

	//Compute percentages
	var percentages []Percentage
	for _, ingredient := range r.Ingredients {
		percent := (ingredient.Quantity / totalFlour) * 100
		percentages = append(percentages, Percentage{Name: ingredient.Name, Percent: percent})
	}

	r.Percentages = percentages
}
