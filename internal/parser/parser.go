package parser

// Parse is the single public entry point. It runs raw recipe text through
// all six pipeline stages and returns a RecipeDTO with confidence metadata.
// Returns ErrInputTooLarge or ErrInputEmpty for invalid inputs.
// For all other inputs, returns a partial RecipeDTO rather than an error.
func Parse(input string) (RecipeDTO, error) {
	if len([]rune(input)) == 0 {
		return RecipeDTO{}, ErrInputEmpty
	}

	cleaned, err := Normalise(input)
	if err != nil {
		return RecipeDTO{}, err
	}

	if len([]rune(cleaned)) == 0 {
		return RecipeDTO{}, ErrInputEmpty
	}

	sm := DetectSections(cleaned)

	dough, other := ParseIngredients(sm.IngredientGroups)

	instructions := ParseInstructions(sm.InstructionLines)

	servings, prepTime, cookTime, additionalTime := ExtractMetadata(sm.MetadataLines)

	dto := RecipeDTO{
		Title:            sm.Title,
		Description:      sm.Description,
		Instructions:     instructions,
		DoughIngredients: dough,
		OtherIngredients: other,
		Servings:         servings,
		PrepTime:         prepTime,
		CookTime:         cookTime,
		AdditionalTime:   additionalTime,
	}

	dto.Confidence = ScoreConfidence(dto, sm.TitleDetectionMethod)

	return dto, nil
}
