package parser

import "testing"

func doughLines(lines ...string) []IngredientGroup {
	return []IngredientGroup{{Phase: "dough", Lines: lines}}
}

func TestParseIngredients_Integer(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("200 g flour"))
	if dough[0].Quantity != 200 || dough[0].Unit != "g" || dough[0].IngredientName != "flour" {
		t.Errorf("got qty=%v unit=%q name=%q", dough[0].Quantity, dough[0].Unit, dough[0].IngredientName)
	}
}

func TestParseIngredients_Decimal(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("2.5 tsp salt"))
	if dough[0].Quantity != 2.5 {
		t.Errorf("expected 2.5, got %v", dough[0].Quantity)
	}
}

func TestParseIngredients_Fraction(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("0.5 tsp salt"))
	if dough[0].Quantity != 0.5 {
		t.Errorf("expected 0.5, got %v", dough[0].Quantity)
	}
}

func TestParseIngredients_MetricFirst_DualMeasurement(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("500 g (4 cups) all purpose flour"))
	if dough[0].Quantity != 500 || dough[0].Unit != "g" {
		t.Errorf("got qty=%v unit=%q", dough[0].Quantity, dough[0].Unit)
	}
	if dough[0].IngredientName != "all purpose flour" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
}

func TestParseIngredients_VolumeFirst_DualMeasurement(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("1 cup (240 grams) water"))
	if dough[0].Quantity != 1 || dough[0].Unit != "cup" {
		t.Errorf("got qty=%v unit=%q", dough[0].Quantity, dough[0].Unit)
	}
	if dough[0].IngredientName != "water" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
}

func TestParseIngredients_WholeItem(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("1 bell pepper, cut into chunks"))
	if dough[0].Unit != "" {
		t.Errorf("whole item should have empty unit, got %q", dough[0].Unit)
	}
	if dough[0].IngredientName != "bell pepper" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
}

func TestParseIngredients_NoQuantityPattern_ToTaste(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("Salt to taste"))
	if dough[0].Quantity != 0 {
		t.Errorf("expected quantity=0, got %v", dough[0].Quantity)
	}
	if dough[0].IngredientName != "salt" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
	if !dough[0].ParseOK {
		t.Error("to-taste pattern should set ParseOK=true")
	}
}

func TestParseIngredients_ParentheticalStripped(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("50 g bubbly, active sourdough starter (100% hydration)"))
	if dough[0].IngredientName != "sourdough starter" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
	if !dough[0].ParseOK {
		t.Error("expected ParseOK=true")
	}
}

func TestParseIngredients_CommaNoteStripped(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("1 cup water, warmed to 100-110 degrees F"))
	if dough[0].IngredientName != "water" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
	if !contains(dough[0].RawLine, "warmed") {
		t.Error("comma note should be preserved in RawLine")
	}
}

func TestParseIngredients_CrossReferenceStripped(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("100 grams ripe sourdough starter from the build above"))
	if dough[0].IngredientName != "sourdough starter" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
}

func TestParseIngredients_YeastAlternatives_TakeFirst(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("2.5g instant dry yeast or 3g active dry yeast or 7.5g fresh yeast"))
	if dough[0].Quantity != 2.5 {
		t.Errorf("expected 2.5, got %v", dough[0].Quantity)
	}
	if dough[0].IngredientName != "instant dry yeast" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
	if !contains(dough[0].RawLine, "active dry yeast") {
		t.Error("alternatives should be preserved in RawLine")
	}
}

func TestParseIngredients_UnsupportedYeast_ParseOKTrue_ConfidenceLow(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("1 tablespoon active dry yeast"))
	if !dough[0].ParseOK {
		t.Error("unsupported yeast should still have ParseOK=true (the line parsed fine)")
	}
	if dough[0].IngredientName != "active dry yeast" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
}

func TestParseIngredients_SourdoughDiscard_Supported(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("0.5 cup sourdough discard"))
	if IsUnsupportedYeast(dough[0].IngredientName) {
		t.Error("sourdough discard should not be classified as unsupported yeast")
	}
}

func TestParseIngredients_NegativeQuantity_ParseOKFalse(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("-2 cups flour"))
	if dough[0].ParseOK {
		t.Error("negative quantity should set ParseOK=false")
	}
}

func TestParseIngredients_Type00Flour_ParsedCorrectly(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("500 g type 00 flour"))
	if dough[0].IngredientName != "type 00 flour" {
		t.Errorf("got name=%q", dough[0].IngredientName)
	}
	if dough[0].Quantity != 500 {
		t.Errorf("got qty=%v", dough[0].Quantity)
	}
}

func TestParseIngredients_ArrayRouting_StarterBuild_ToDough(t *testing.T) {
	groups := []IngredientGroup{
		{Phase: "starter build", Lines: []string{"30 g sourdough starter", "35 g flour"}},
		{Phase: "pesto", Lines: []string{"30 g olive oil"}},
	}
	dough, other := ParseIngredients(groups)
	if len(dough) != 2 {
		t.Errorf("starter build should route to doughIngredients, got %d", len(dough))
	}
	if len(other) != 1 {
		t.Errorf("pesto should route to otherIngredients, got %d", len(other))
	}
	if other[0].Phase != "pesto" {
		t.Errorf("expected phase 'pesto', got %q", other[0].Phase)
	}
}

func TestParseIngredients_BulletStripped(t *testing.T) {
	for _, bullet := range []string{"- 200g flour", "* 200g flour", "• 200g flour", "— 200g flour"} {
		dough, _ := ParseIngredients(doughLines(bullet))
		if dough[0].IngredientName != "flour" {
			t.Errorf("bullet not stripped for %q: got name %q", bullet, dough[0].IngredientName)
		}
	}
}

func TestParseIngredients_ParseOKFalse_EmptyName(t *testing.T) {
	dough, _ := ParseIngredients(doughLines("100 g"))
	if dough[0].ParseOK {
		t.Error("empty ingredient name should set ParseOK=false")
	}
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
