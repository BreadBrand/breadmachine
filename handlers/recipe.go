package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RedBrand88/breadmachine/models"
)

//GetAllRecipes returns all recipes from Firestore
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := FetchAllRecipesFromFirebase()
	if err != nil {
		http.Error(w, "Failed to fetch recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

//CreateRecipe handles recipe creation
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//Calculate Baker's Percentages before saving
	recipe.CalculateBakersPercentages()

	_, err := SaveRecipe(recipe)
	if err != nil {
		http.Error(w, "Failed to save recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// RecipeHandler processes both GET and DELETE requests for a single recipe
func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/recipes/")

	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		docRef := client.Collection("Recipes").Doc(id)
		doc, err := docRef.Get(r.Context())
		if err != nil || !doc.Exists() {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		var recipe models.Recipe
		doc.DataTo(&recipe)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)

	case "DELETE":
		docRef := client.Collection("Recipes").Doc(id)
		_, err := docRef.Delete(r.Context())
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
