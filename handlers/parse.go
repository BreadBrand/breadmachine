package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RedBrand88/breadmachine/internal/parser"
)

type apiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiError{Error: code, Message: message})
}

// ParseHandler handles POST /api/recipes/parse.
// Accepts JSON body {"text": "..."} and returns a RecipeDTO.
func ParseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST only")
		return
	}

	// Auth — same pattern as CreateRecipe
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Please sign in to import recipes.")
		return
	}
	if _, err := authClient.VerifyIDToken(r.Context(), tokenString); err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Please sign in to import recipes.")
		return
	}

	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Something went wrong with the request. Please try again.")
		return
	}

	if strings.TrimSpace(body.Text) == "" {
		writeError(w, http.StatusBadRequest, "INPUT_EMPTY", "Please paste some recipe text before submitting.")
		return
	}

	dto, err := parser.Parse(body.Text)
	if err != nil {
		switch err {
		case parser.ErrInputTooLarge:
			writeError(w, http.StatusBadRequest, "INPUT_TOO_LARGE",
				"This looks like more than a recipe. Try copying from the print view instead.")
		case parser.ErrInputEmpty:
			writeError(w, http.StatusBadRequest, "INPUT_EMPTY",
				"Please paste some recipe text before submitting.")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong.")
		}
		return
	}

	// Minimum viability: must have at least ingredients OR instructions
	if len(dto.DoughIngredients) == 0 && len(dto.Instructions) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "PARSE_FAILED",
			"We couldn't find a recipe in that text. Try the print view, or paste just the recipe section.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
