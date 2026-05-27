package handlers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/RedBrand88/breadmachine/models"
	"github.com/RedBrand88/breadmachine/utility"
)

// ParseRecipeText converts a raw recipe blob into a structured Recipe.
// This is the legacy markdown-based parser; new code should use internal/parser.Parse.
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
		}
	}

	return recipe
}

func parseIngredientLine(line string, phase models.Phase) models.Ingredient {
	re := regexp.MustCompile(`([\d¼½¾⅓⅔⅛\.\s\/]+)\s*([a-zA-Z]+)?\s*(.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 4 {
		return models.Ingredient{IngredientName: line, Phase: phase}
	}
	qty := parseNumber(strings.TrimSpace(matches[1]))
	unit := strings.ToLower(matches[2])
	mult := unitToGrams[unit]
	grams := qty * mult

	ing := models.Ingredient{
		IngredientName: matches[3],
		Quantity:       qty,
		Unit:           matches[2],
		Grams:          grams,
		Phase:          phase,
		DensityGPerMl:  utility.LookupDensity(matches[3]),
	}
	return ing
}

// parseNumber converts fractional Unicode numbers to floats.
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
