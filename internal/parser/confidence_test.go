package parser

import "testing"

func TestScoreConfidence_Title(t *testing.T) {
	dto := RecipeDTO{}
	cases := []struct {
		method TitleDetectionMethod
		want   float32
	}{
		{TitleFromHeader, 1.0},
		{TitleHeuristic, 0.7},
		{TitleEmpty, 0.3},
	}
	for _, c := range cases {
		got := ScoreConfidence(dto, c.method)
		if got.Title != c.want {
			t.Errorf("method=%v: Title = %v, want %v", c.method, got.Title, c.want)
		}
	}
}

func TestScoreConfidence_Ingredients_Zero(t *testing.T) {
	dto := RecipeDTO{}
	got := ScoreConfidence(dto, TitleFromHeader)
	if got.Ingredients != 0.2 {
		t.Errorf("0 ingredients: expected 0.2, got %v", got.Ingredients)
	}
}

func TestScoreConfidence_Ingredients_VeryFew(t *testing.T) {
	dto := RecipeDTO{
		DoughIngredients: []IngredientDTO{
			{IngredientName: "flour", ParseOK: true},
		},
	}
	got := ScoreConfidence(dto, TitleFromHeader)
	if got.Ingredients != 0.5 {
		t.Errorf("1 ingredient: expected 0.5, got %v", got.Ingredients)
	}
}

func TestScoreConfidence_Ingredients_Few(t *testing.T) {
	dto := RecipeDTO{
		DoughIngredients: []IngredientDTO{
			{IngredientName: "flour", ParseOK: true},
			{IngredientName: "water", ParseOK: true},
			{IngredientName: "salt", ParseOK: true},
		},
	}
	got := ScoreConfidence(dto, TitleFromHeader)
	if got.Ingredients != 0.7 {
		t.Errorf("3 ingredients: expected 0.7, got %v", got.Ingredients)
	}
}

func TestScoreConfidence_Ingredients_ParseOKFalse(t *testing.T) {
	dto := RecipeDTO{
		DoughIngredients: []IngredientDTO{
			{IngredientName: "flour", ParseOK: true},
			{IngredientName: "water", ParseOK: true},
			{IngredientName: "salt", ParseOK: true},
			{IngredientName: "yeast", ParseOK: true},
			{IngredientName: "", ParseOK: false},
		},
	}
	got := ScoreConfidence(dto, TitleFromHeader)
	if got.Ingredients != 0.7 {
		t.Errorf("any ParseOK=false: expected 0.7, got %v", got.Ingredients)
	}
}

func TestScoreConfidence_Ingredients_UnsupportedYeastDominates(t *testing.T) {
	dto := RecipeDTO{
		DoughIngredients: []IngredientDTO{
			{IngredientName: "flour", ParseOK: true},
			{IngredientName: "water", ParseOK: true},
			{IngredientName: "salt", ParseOK: true},
			{IngredientName: "sugar", ParseOK: true},
			{IngredientName: "butter", ParseOK: true},
			{IngredientName: "active dry yeast", ParseOK: true},
		},
	}
	got := ScoreConfidence(dto, TitleFromHeader)
	if got.Ingredients != 0.4 {
		t.Errorf("unsupported yeast should dominate: expected 0.4, got %v", got.Ingredients)
	}
}

func TestScoreConfidence_Instructions(t *testing.T) {
	cases := []struct {
		n    int
		want float32
	}{
		{0, 0.3},
		{1, 0.7},
		{2, 0.7},
		{3, 1.0},
		{10, 1.0},
	}
	for _, c := range cases {
		steps := make([]string, c.n)
		dto := RecipeDTO{Instructions: steps}
		got := ScoreConfidence(dto, TitleFromHeader)
		if got.Instructions != c.want {
			t.Errorf("n=%d: Instructions = %v, want %v", c.n, got.Instructions, c.want)
		}
	}
}
