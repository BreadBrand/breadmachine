package parser

// ScoreConfidence scores the parsed RecipeDTO on a 0.0–1.0 scale per field.
// Scores below 0.6 are flagged in the review modal.
// Multiple conditions on the same field use min() — pessimistic stacking.
func ScoreConfidence(dto RecipeDTO, method TitleDetectionMethod) ConfidenceMeta {
	return ConfidenceMeta{
		Title:        scoreTitle(method),
		Ingredients:  scoreIngredients(dto),
		Instructions: scoreInstructions(dto),
	}
}

func scoreTitle(method TitleDetectionMethod) float32 {
	switch method {
	case TitleFromHeader:
		return 1.0
	case TitleHeuristic:
		return 0.7
	default:
		return 0.3
	}
}

func scoreIngredients(dto RecipeDTO) float32 {
	total := len(dto.DoughIngredients) + len(dto.OtherIngredients)
	if total == 0 {
		return 0.2
	}

	score := float32(1.0)

	switch {
	case total < 3:
		score = min(score, 0.5)
	case total < 5:
		score = min(score, 0.7)
	}

	all := make([]IngredientDTO, 0, len(dto.DoughIngredients)+len(dto.OtherIngredients))
	all = append(all, dto.DoughIngredients...)
	all = append(all, dto.OtherIngredients...)
	for _, ing := range all {
		if !ing.ParseOK {
			score = min(score, 0.7)
		}
		if IsUnsupportedYeast(ing.IngredientName) {
			score = min(score, 0.4)
		}
	}

	return score
}

func scoreInstructions(dto RecipeDTO) float32 {
	n := len(dto.Instructions)
	switch {
	case n >= 3:
		return 1.0
	case n >= 1:
		return 0.7
	default:
		return 0.3
	}
}
