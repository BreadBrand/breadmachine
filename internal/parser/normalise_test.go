package parser

import (
	"strings"
	"testing"
)

func TestNormalise_TooLarge(t *testing.T) {
	input := strings.Repeat("a", 10001)
	_, err := Normalise(input)
	if err != ErrInputTooLarge {
		t.Fatalf("expected ErrInputTooLarge, got %v", err)
	}
}

func TestNormalise_ExactLimit(t *testing.T) {
	input := strings.Repeat("a", 10000)
	_, err := Normalise(input)
	if err != nil {
		t.Fatalf("10000 runes should be accepted, got %v", err)
	}
}

func TestNormalise_HTMLTags(t *testing.T) {
	got, _ := Normalise("<b>flour</b> and <em>water</em>")
	if got != "flour and water" {
		t.Errorf("got %q", got)
	}
}

func TestNormalise_HTMLEntities(t *testing.T) {
	cases := []struct{ in, want string }{
		{"1 &amp; 2", "1 & 2"},
		{"&frac12; cup", "1/2 cup"},
		{"&frac14; tsp", "1/4 tsp"},
		{"&frac34; cup", "3/4 cup"},
		{"&nbsp;flour", " flour"},
	}
	for _, c := range cases {
		got, _ := Normalise(c.in)
		if got != c.want {
			t.Errorf("Normalise(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalise_MarkdownStripping(t *testing.T) {
	cases := []struct{ in, want string }{
		{"**flour**", "flour"},
		{"__flour__", "flour"},
		{"*flour*", "flour"},
		{"_flour_", "flour"},
		{"# Ingredients", "Ingredients"},
		{"## Method", "Method"},
		{"[Bake](https://example.com)", "Bake"},
		{"`code`", "code"},
	}
	for _, c := range cases {
		got, _ := Normalise(c.in)
		if got != c.want {
			t.Errorf("Normalise(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalise_TrailingAsterisks(t *testing.T) {
	cases := []struct{ in, want string }{
		{"sourdough starter**", "sourdough starter"},
		{"starter*", "starter"},
		{"flour***", "flour"},
	}
	for _, c := range cases {
		got, _ := Normalise(c.in)
		if got != c.want {
			t.Errorf("Normalise(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalise_CurlyQuotes(t *testing.T) {
	got, _ := Normalise("‘starter’") // left/right single quotes
	if got != "'starter'" {
		t.Errorf("got %q", got)
	}
	got2, _ := Normalise("“fresh bread”") // left/right double quotes
	if got2 != `"fresh bread"` {
		t.Errorf("got %q", got2)
	}
}

func TestNormalise_UnicodeFractions(t *testing.T) {
	cases := []struct{ in, want string }{
		{"½ cup", "1/2 cup"},
		{"¼ tsp", "1/4 tsp"},
		{"¾ cup", "3/4 cup"},
		{"⅓ cup", "1/3 cup"},
		{"⅔ cup", "2/3 cup"},
	}
	for _, c := range cases {
		got, _ := Normalise(c.in)
		if got != c.want {
			t.Errorf("Normalise(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalise_MixedNumbers(t *testing.T) {
	cases := []struct{ in, want string }{
		{"2 3/4 cups", "2.75 cups"},
		{"1 1/2 tsp", "1.5 tsp"},
		{"1 1/4 cup water", "1.25 cup water"},
	}
	for _, c := range cases {
		got, _ := Normalise(c.in)
		if got != c.want {
			t.Errorf("Normalise(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalise_MixedNumbers_DoesNotAlterType00Flour(t *testing.T) {
	// "00" in "type 00 flour" must not be combined with anything
	got, _ := Normalise("500g type 00 flour")
	if got != "500g type 00 flour" {
		t.Errorf("type 00 flour was altered: got %q", got)
	}
}

func TestNormalise_AttributionLines(t *testing.T) {
	cases := []string{
		"Submitted by Terri McCarrell",
		"Tested by Allrecipes Test Kitchen",
		"Reviewed by staff",
		"Recipe by John",
		"Author: Emilie Raffa",
	}
	for _, line := range cases {
		got, _ := Normalise(line)
		if strings.TrimSpace(got) != "" {
			t.Errorf("attribution line not stripped: Normalise(%q) = %q", line, got)
		}
	}
}

func TestNormalise_SourceLines(t *testing.T) {
	cases := []string{
		"Find it online: https://www.theclevercarrot.com/recipe",
		"Source: The Joy of Cooking",
		"Originally published at breadblog.com",
		"https://example.com/recipe",
	}
	for _, line := range cases {
		got, _ := Normalise(line)
		if strings.TrimSpace(got) != "" {
			t.Errorf("source line not stripped: Normalise(%q) = %q", line, got)
		}
	}
}

func TestNormalise_ScalingArtifacts(t *testing.T) {
	got, _ := Normalise("Yield: 1 loaf 1x")
	if strings.Contains(got, "1x") {
		t.Errorf("scaling artifact not stripped: got %q", got)
	}
}

func TestNormalise_NutritionBlock(t *testing.T) {
	input := "Directions\nMix everything.\n\nNutrition Facts\nCalories 200\nProtein 5g\n\nNotes\nStore in fridge."
	got, _ := Normalise(input)
	if strings.Contains(got, "Calories") || strings.Contains(got, "Protein") {
		t.Errorf("nutrition block not stripped: got %q", got)
	}
	if !strings.Contains(got, "Mix everything") {
		t.Errorf("content before nutrition block was stripped: got %q", got)
	}
}

func TestNormalise_BrowserPrintArtifacts(t *testing.T) {
	input := "5/26/26, 1:40 PM Basic All Purpose Sourdough Sandwich Bread\nabout:blank 1/3\n500g flour\nabout:blank 2/3"
	got, _ := Normalise(input)
	if strings.Contains(got, "about:blank") {
		t.Errorf("about:blank not stripped: got %q", got)
	}
	if strings.Contains(got, "1:40 PM") {
		t.Errorf("browser header not stripped: got %q", got)
	}
	if !strings.Contains(got, "500g flour") {
		t.Errorf("real content was stripped: got %q", got)
	}
}

func TestNormalise_InlineMetadataExpansion(t *testing.T) {
	input := "Prep Time: 12 hours Cook Time: 50 minutes Yield: 1 loaf"
	got, _ := Normalise(input)
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 3 {
		t.Errorf("inline metadata not expanded into separate lines: got %d lines:\n%s", len(lines), got)
	}
}

func TestNormalise_CollapseBlankLines(t *testing.T) {
	input := "line1\n\n\n\nline2"
	got, _ := Normalise(input)
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("triple blank line not collapsed: got %q", got)
	}
}

func TestNormalise_DegreeSymbolPreserved(t *testing.T) {
	got, _ := Normalise("Bake at 220°C")
	if !strings.Contains(got, "°") {
		t.Errorf("degree symbol was stripped: got %q", got)
	}
}
