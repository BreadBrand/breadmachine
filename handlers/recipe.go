package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/RedBrand88/breadmachine/models"
	"github.com/google/uuid"
)

// GetAllRecipes returns all recipes from Firestore
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := FetchAllRecipesFromFirebase()
	if err != nil {
		http.Error(w, "Failed to fetch recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// CreateRecipe handles recipe creation
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Expect "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
		return
	}

	// Verify the Firebase ID token
	token, err := authClient.VerifyIDToken(r.Context(), tokenString)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	var recipe models.Recipe

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Assign Firestore-generated ID for this recipe
	docRef := client.Collection("Recipes").NewDoc()
	recipe.ID = docRef.ID

	// Set the UserID from the verified token
	recipe.UserID = token.UID

	// Sum total dough and generate IDs for each ingredient
	totalDough := 0.0
	for i := range recipe.Ingredients {
		totalDough += recipe.Ingredients[i].Quantity
		if recipe.Ingredients[i].ID == "" {
			recipe.Ingredients[i].ID = uuid.NewString()
		}
	}

	// Populate Meta data
	now := time.Now()
	recipe.Meta = models.Meta{
		YieldGrams: totalDough,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Calculate baker's percentages into each ingredient
	recipe.CalculateBakerPercentages()

	// Persist the recipe document to Firestore
	if _, err := docRef.Set(r.Context(), recipe); err != nil {
		http.Error(w, "Failed to save recipe", http.StatusInternalServerError)
		return
	}

	// Return the created recipe (with ID and computed fields)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// RecipeHandler processes GET and DELETE for a single recipe
func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/recipes/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		docRef := client.Collection("Recipes").Doc(id)
		docSnap, err := docRef.Get(r.Context())
		if err != nil || !docSnap.Exists() {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		var recipe models.Recipe
		if err := docSnap.DataTo(&recipe); err != nil {
			http.Error(w, "Error parsing recipe", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)

	case "DELETE":
		docRef := client.Collection("Recipes").Doc(id)
		if _, err := docRef.Delete(r.Context()); err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func RecipesHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    GetAllRecipes(w, r)
  case http.MethodPost:
    CreateRecipe(w, r)
  default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}
