package parser

import (
	"strings"
	"testing"
)

func TestParseInstructions_NumberedList(t *testing.T) {
	lines := []string{"1. Mix flour and water.", "2. Knead for 10 minutes.", "3. Let rise."}
	got := ParseInstructions(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	if got[0] != "Mix flour and water." {
		t.Errorf("step prefix not stripped: %q", got[0])
	}
}

func TestParseInstructions_BulletList(t *testing.T) {
	lines := []string{"- Mix ingredients.", "- Knead.", "- Proof."}
	got := ParseInstructions(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d: %v", len(got), got)
	}
}

func TestParseInstructions_NumberedListReset(t *testing.T) {
	lines := []string{
		"1. Make the dough.",
		"2. Refrigerate.",
		"* Fold once.",
		"* Fold twice.",
		"3. Shape.",
	}
	got := ParseInstructions(lines)
	if len(got) < 5 {
		t.Errorf("expected at least 5 instructions, got %d", len(got))
	}
}

func TestParseInstructions_BlankLineParagraphs(t *testing.T) {
	lines := []string{"Mix flour and water.", "", "Knead the dough.", "", "Let rise."}
	got := ParseInstructions(lines)
	if len(got) != 3 {
		t.Errorf("expected 3 paragraph instructions, got %d", len(got))
	}
}

func TestParseInstructions_DisclaimerExcluded(t *testing.T) {
	lines := []string{
		"1. Bake at 200C for 25 minutes.",
		"Keep in mind that conditions in each kitchen vary.",
		"Your oven may differ.",
	}
	got := ParseInstructions(lines)
	for _, s := range got {
		if strings.Contains(s, "Keep in mind") || strings.Contains(s, "Your oven") {
			t.Errorf("disclaimer should be excluded: %q", s)
		}
	}
}

func TestParseInstructions_NoEmptyStrings(t *testing.T) {
	lines := []string{"1. Step one.", "", "2. Step two.", ""}
	got := ParseInstructions(lines)
	for _, s := range got {
		if s == "" {
			t.Error("output slice must not contain empty strings")
		}
	}
}

func TestParseInstructions_NoteIncluded(t *testing.T) {
	lines := []string{"1. Mix.", "Note: Dough should be sticky.", "2. Shape."}
	got := ParseInstructions(lines)
	hasNote := false
	for _, s := range got {
		if strings.Contains(s, "Note:") {
			hasNote = true
		}
	}
	if !hasNote {
		t.Error("inline Note: paragraph should be included in instructions")
	}
}
