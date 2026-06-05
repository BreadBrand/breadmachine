package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/RedBrand88/breadmachine/models"
	"github.com/RedBrand88/breadmachine/utility"
	"github.com/google/uuid"
)

// unitToGrams maps canonical unit strings to their gram equivalent.
// Volume units use water density (1 g/ml) as the base; density lookup in
// normalizeIngredients adjusts the effective weight per ingredient.
var unitToGrams = map[string]float64{
	"g":     1,
	"kg":    1000,
	"oz":    28.35,
	"lb":    453.592,
	"fl oz": 29.57,
	"ml":    1,
	"l":     1000,
	"tsp":   5,
	"tbsp":  15,
	"cup":   240,
}

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

	// Validate yeastType if provided
	if recipe.YeastType != "" &&
		recipe.YeastType != models.YeastTypeDry &&
		recipe.YeastType != models.YeastTypeSourdough {
		http.Error(w, "Invalid yeastType: must be 'dry' or 'sourdough'", http.StatusBadRequest)
		return
	}

	// Assign Firestore-generated ID for this recipe
	docRef := client.Collection("Recipes").NewDoc()
	recipe.ID = docRef.ID
	recipe.UserID = token.UID

	// Normalize ingredients and generate IDs for each entry
	now := time.Now()
	doughTotal, normalizedDough := normalizeIngredients(recipe.DoughIngredients)
	recipe.DoughIngredients = normalizedDough
	otherTotal, normalizedOther := normalizeIngredients(recipe.OtherIngredients)
	recipe.OtherIngredients = normalizedOther

	// Populate computed Meta fields; preserve any client-supplied fields (e.g. PrepTime).
	recipe.Meta.YieldGrams = doughTotal + otherTotal
	recipe.Meta.CreatedAt = now
	recipe.Meta.UpdatedAt = now

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

func normalizeIngredients(ingredients []models.Ingredient) (float64, []models.Ingredient) {
	totalDough := 0.0

	for i := range ingredients {
		ing := &ingredients[i]

		if ing.ID == "" {
			ing.ID = uuid.NewString()
		}

		unit := strings.ToLower(ing.Unit)

		if ing.Grams == 0 && ing.Quantity > 0 {
			mult := unitToGrams[unit]
			if mult == 0 {
				mult = 1
			}
			ing.Grams = ing.Quantity * mult
		}

		if ing.DensityGPerMl == 0 {
			ing.DensityGPerMl = utility.LookupDensity(ing.IngredientName)
		}

		totalDough += ing.Grams
	}
	return totalDough, ingredients
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
