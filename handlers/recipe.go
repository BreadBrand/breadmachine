package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BreadBrand/breadmachine/internal/parser"
	"github.com/BreadBrand/breadmachine/models"
	"github.com/BreadBrand/breadmachine/utility"
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
	uid, ok := authenticate(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		recipe.YeastType != models.YeastTypeSourdough &&
		recipe.YeastType != models.YeastTypeNone {
		http.Error(w, "Invalid yeastType: must be 'dry', 'sourdough', or 'none'", http.StatusBadRequest)
		return
	}

	// Assign Firestore-generated ID for this recipe
	docRef := client.Collection("Recipes").NewDoc()
	recipe.ID = docRef.ID
	recipe.UserID = uid

	// Normalize ingredients and generate IDs for each entry
	now := time.Now()
	doughTotal, normalizedDough := normalizeIngredients(recipe.DoughIngredients)
	recipe.DoughIngredients = convertToGrams(normalizedDough)
	_, normalizedOther := normalizeIngredients(recipe.OtherIngredients)
	recipe.OtherIngredients = convertToGrams(normalizedOther)

	// Populate computed Meta fields; preserve any client-supplied fields (e.g. PrepTime).
	// YieldGrams is dough weight only — otherIngredients (toppings, fillings) are excluded.
	recipe.Meta.YieldGrams = doughTotal
	recipe.Meta.CreatedAt = now
	recipe.Meta.UpdatedAt = now

	// Calculate baker's percentages into each ingredient
	models.CalculateBakerPercentages(recipe.DoughIngredients)

	// Persist the recipe document to Firestore
	if _, err := docRef.Set(r.Context(), recipe); err != nil {
		http.Error(w, "Failed to save recipe", http.StatusInternalServerError)
		return
	}

	// Return the created recipe (with ID and computed fields)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// GetRecipe returns a single recipe by ID.
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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
}

// DeleteRecipe deletes a single recipe by ID. The caller must be authenticated
// and must own the recipe (recipe.UserID must match the verified token's UID).
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	uid, ok := authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Please sign in to delete a recipe.")
		return
	}

	id := r.PathValue("id")
	docRef := client.Collection("Recipes").Doc(id)
	docSnap, err := docRef.Get(r.Context())
	if err != nil || !docSnap.Exists() {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Recipe not found.")
		return
	}

	var recipe models.Recipe
	if err := docSnap.DataTo(&recipe); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong.")
		return
	}
	if recipe.UserID != uid {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "You don't have permission to delete this recipe.")
		return
	}

	if _, err := docRef.Delete(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete recipe.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func normalizeIngredients(ingredients []models.Ingredient) (float64, []models.Ingredient) {
	totalDough := 0.0

	for i := range ingredients {
		ing := &ingredients[i]

		if ing.ID == "" {
			ing.ID = uuid.NewString()
		}

		unit := parser.CanonicalUnit(strings.ToLower(ing.Unit))
		ing.Unit = unit

		if ing.Quantity > 0 {
			if unit == "count" {
				// Always recompute count ingredients — stored grams may be stale if
				// the per-unit weight table changed since the recipe was first saved.
				ing.Grams = ing.Quantity * utility.LookupCountWeight(ing.IngredientName)
			} else if ing.Grams == 0 {
				mult := unitToGrams[unit]
				if mult == 0 {
					mult = 1
				}
				ing.Grams = ing.Quantity * mult
			}
		}

		if ing.DensityGPerMl == 0 {
			ing.DensityGPerMl = utility.LookupDensity(ing.IngredientName)
		}

		totalDough += ing.Grams
	}
	return totalDough, ingredients
}
