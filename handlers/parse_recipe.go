package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/RedBrand88/breadmachine/models"
	"github.com/RedBrand88/breadmachine/utility"
)

// Improved ingredient line detection

// Base unit conversions → grams
var unitToGrams = map[string]float64{
	"g":     1,
	"gram":  1,
	"grams": 1,
	"kg":    1000,
	"oz":    28.35,
	"ml":    1,
	"tsp":   5,
	"tbsp":  15,
	"tbs":   15,
	"cup":   240,
	"cups":  240,
}

// ParseRecipeHandler receives raw text and returns structured JSON.
func ParseRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	recipe := ParseRecipeText(body.Text)

	recipe.CalculateBakerPercentages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// ParseRecipeText converts a raw recipe blob into a structured Recipe.
func ParseRecipeText(raw string) models.Recipe {
	lines := strings.Split(raw, "\n")

	headerRegex := regexp.MustCompile(`^#\s*(\w+)\s*$`)
	subHeaderRegex := regexp.MustCompile(`^\[(.+?)\]\s*$`)

	recipe := models.Recipe{
		Title:       "",
		Description: "",
		Ingredients: []models.Ingredient{},
	}

	var (
		currentHeader    string
		currentSubHeader string
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Match section header like "# Title"
		if match := headerRegex.FindStringSubmatch(line); match != nil {
			currentHeader = strings.ToLower(match[1])
			currentSubHeader = ""
			continue
		}

		// Match sub-header like "[dough]"
		if match := subHeaderRegex.FindStringSubmatch(line); match != nil {
			currentSubHeader = strings.ToLower(match[1])
			continue
		}

		// Handle line content
		switch currentHeader {
		case "title":
			recipe.Title = line

		case "description":
			if recipe.Description != "" {
				recipe.Description += " "
			}
			recipe.Description += line

		case "ingredients":
			phase := currentSubHeader
			if phase == "" {
				phase = "dough"
			}

			ing := parseIngredientLine(line, models.Phase(phase))
			if ing.IngredientName != "" {
				recipe.Ingredients = append(recipe.Ingredients, ing)
			}

		case "instructions":
			clean := strings.TrimSpace(line)
			clean = regexp.MustCompile(`^\d+[\.\)]\s*`).ReplaceAllString(clean, "")
			recipe.Instructions = append(recipe.Instructions, clean)

		// case "notes":
		// 	recipe.Notes = append(recipe.Notes, line)
		}
	}

	return recipe
}

func parseIngredientLine(line string, phase models.Phase) models.Ingredient {
	//example line = "4 Tbsp salted butter (softened in microwave for 15 seconds)"
	re := regexp.MustCompile(`([\d¼½¾⅓⅔⅛\.\s\/]+)\s*([a-zA-Z]+)?\s*(.*)`)
	matches := re.FindStringSubmatch(line)
	//[]string{
	//"4 Tbsp salted butter (softened in microwave for 15 seconds)"
	//"4"
	//"Tbsp"
	//"salted butter (softened in microwave for 15 seconds)"
	//}
	if len(matches) < 4 {
		return models.Ingredient{IngredientName: line, Phase: phase}
	}
	//so I can use matches
	qty := parseNumber(strings.TrimSpace(matches[1]))
	unit := strings.ToLower(matches[2])
	mult := unitToGrams[unit]
	grams := qty * mult

	ing := models.Ingredient{
		IngredientName: matches[3],
		Quantity: qty,
		Unit: matches[2],
		Grams: grams,
		Phase: phase,
		DensityGPerMl: utility.LookupDensity(matches[3]),
	}
	return ing
}

// Converts fractional Unicode numbers to floats
func parseNumber(s string) float64 {
	replacements := map[string]string{
		"¼": "1/4",
		"½": "1/2",
		"¾": "3/4",
		"⅓": "1/3",
		"⅔": "2/3",
		"⅛": "1/8",
	}

	for k, v := range replacements {
		s = strings.ReplaceAll(s, k, v)
	}

	parts := strings.Fields(s)
	total := 0.0
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "/") {
			frParts := strings.Split(part, "/")
			if len(frParts) == 2 {
				num, err1 := strconv.ParseFloat(frParts[0], 64)
				den, err2 := strconv.ParseFloat(frParts[1], 64)
				if err1 == nil && err2 == nil && den != 0 {
					total += num / den
					continue
				}
			}
		}
		if f, err := strconv.ParseFloat(part, 64); err == nil {
			total += f
		}
	}
	return total
}
