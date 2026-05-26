package parser

import (
	"strings"
	"testing"
)

func TestDetectSections_KeywordMidSentence_DoesNotTrigger(t *testing.T) {
	// "ingredients" mid-sentence must not trigger section detection
	input := "This recipe uses the best ingredients I have ever tasted in my life"
	sm := DetectSections(input)
	if len(sm.IngredientGroups) > 0 && len(sm.IngredientGroups[0].Lines) > 0 {
		t.Error("mid-sentence 'ingredients' should not trigger section detection")
	}
}

func TestDetectSections_KeywordOnOwnLine_Triggers(t *testing.T) {
	input := "My Great Bread\n\nIngredients\n500g flour\n200g water\n\nDirections\nMix together."
	sm := DetectSections(input)
	if len(sm.IngredientGroups) == 0 {
		t.Fatal("expected ingredient group, got none")
	}
	if len(sm.IngredientGroups[0].Lines) != 2 {
		t.Errorf("expected 2 ingredient lines, got %d", len(sm.IngredientGroups[0].Lines))
	}
	if len(sm.InstructionLines) == 0 {
		t.Error("expected instruction lines")
	}
}

func TestDetectSections_KeywordWithColon_Triggers(t *testing.T) {
	input := "My Bread\n\nIngredients:\n400g flour\n\nInstructions:\nMix well."
	sm := DetectSections(input)
	if len(sm.IngredientGroups) == 0 || len(sm.IngredientGroups[0].Lines) == 0 {
		t.Error("keyword with colon should trigger section detection")
	}
}

func TestDetectSections_TitleFromHeader(t *testing.T) {
	input := "Roasted Garlic Focaccia\n\nIngredients\n500g flour"
	sm := DetectSections(input)
	if sm.Title != "Roasted Garlic Focaccia" {
		t.Errorf("expected title 'Roasted Garlic Focaccia', got %q", sm.Title)
	}
	if sm.TitleDetectionMethod != TitleHeuristic {
		t.Errorf("expected TitleHeuristic, got %v", sm.TitleDetectionMethod)
	}
}

func TestDetectSections_TitleEmpty_WhenFirstLineIsIngredient(t *testing.T) {
	input := "500g flour\n200g water\n\nInstructions\nMix."
	sm := DetectSections(input)
	if sm.TitleDetectionMethod != TitleEmpty {
		t.Errorf("first line starting with digit should not be title; got method %v, title %q",
			sm.TitleDetectionMethod, sm.Title)
	}
}

func TestDetectSections_SubsectionHeaders(t *testing.T) {
	input := "My Focaccia\n\nIngredients\nFor the dough –\n200g flour\n50g water\n\nFor the pesto –\n30g olive oil\n20g parmesan\n\nDirections\nMix."
	sm := DetectSections(input)
	if len(sm.IngredientGroups) < 2 {
		t.Fatalf("expected 2 ingredient groups, got %d", len(sm.IngredientGroups))
	}
	if sm.IngredientGroups[0].Phase != "dough" {
		t.Errorf("expected phase 'dough', got %q", sm.IngredientGroups[0].Phase)
	}
	if sm.IngredientGroups[1].Phase != "pesto" {
		t.Errorf("expected phase 'pesto', got %q", sm.IngredientGroups[1].Phase)
	}
}

func TestDetectSections_EmptyGroupDiscarded(t *testing.T) {
	// Subsection header with no lines before next header → discarded
	input := "Bread\n\nIngredients\nFor the dough –\nFor the topping –\n20g cheese\n\nDirections\nMix."
	sm := DetectSections(input)
	for _, g := range sm.IngredientGroups {
		if len(g.Lines) == 0 {
			t.Errorf("empty group with phase %q should have been discarded", g.Phase)
		}
	}
}

func TestDetectSections_BakersPercentageBlockSkipped(t *testing.T) {
	input := "Bread\n\nIngredients\n500g flour\n300g water\n\nFinal Dough Baker's Percentages\n100% all purpose flour\n65% liquid\n\nInstructions\nMix."
	sm := DetectSections(input)
	for _, g := range sm.IngredientGroups {
		for _, line := range g.Lines {
			if strings.Contains(line, "%") && strings.Contains(strings.ToLower(line), "flour") {
				t.Errorf("baker's percentage line should be skipped: %q", line)
			}
		}
	}
}

func TestDetectSections_NoLineBreaks(t *testing.T) {
	// Fewer than 3 newlines → NoLineBreaks = true
	input := "title\ningredients\n500g flour"
	sm := DetectSections(input)
	if !sm.NoLineBreaks {
		t.Error("expected NoLineBreaks=true for input with < 3 newlines")
	}
}

func TestDetectSections_InstructionSubHeaderPrepended(t *testing.T) {
	input := "Bread\n\nIngredients\n500g flour\n\nInstructions\nMix the Dough\nCombine flour and water in a bowl."
	sm := DetectSections(input)
	if len(sm.InstructionLines) == 0 {
		t.Fatal("expected instruction lines")
	}
	if !strings.HasPrefix(sm.InstructionLines[0], "Mix the Dough:") {
		t.Errorf("sub-header not prepended; got %q", sm.InstructionLines[0])
	}
}

func TestDetectSections_DescriptionCappedAt2000(t *testing.T) {
	longDesc := strings.Repeat("This is a description sentence. ", 80) // >2000 chars
	input := "My Bread\n\n" + longDesc + "\n\nIngredients\n500g flour"
	sm := DetectSections(input)
	if len(sm.Description) > 2000 {
		t.Errorf("description exceeds 2000 chars: %d", len(sm.Description))
	}
}
