package parser

import "testing"

func TestParseMinutes(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		// plain minutes
		{"20 mins", 20},
		{"30 minutes", 30},
		{"45 min", 45},
		// plain hours
		{"1 hour", 60},
		{"2 hours", 120},
		{"2 hrs", 120},
		// combined
		{"1 hour 30 minutes", 90},
		{"2 hours 15 mins", 135},
		// site artifacts — concatenated with no space
		{"10minutes minutes", 10},
		{"22minutes minutes", 22},
		{"2hours hrs", 120},
		// range — take the first number
		{"15-17 minutes", 15},
		// empty / unknown
		{"", 0},
		{"overnight", 0},
	}
	for _, c := range cases {
		got := parseMinutes(c.in)
		if got != c.want {
			t.Errorf("parseMinutes(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestExtractMetadata_PrepTime(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"Prep Time: 20 mins", 20},
		{"Prep: 30 minutes", 30},
		{"Preparation Time: 1 hour", 60},
	}
	for _, c := range cases {
		_, prep, _, _ := ExtractMetadata([]string{c.in})
		if prep != c.want {
			t.Errorf("ExtractMetadata(%q) prepTime = %d, want %d", c.in, prep, c.want)
		}
	}
}

func TestExtractMetadata_CookTime(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"Cook Time: 15 mins", 15},
		{"Baking Time: 55 minutes", 55},
		{"Bake Time: 30 mins", 30},
	}
	for _, c := range cases {
		_, _, cook, _ := ExtractMetadata([]string{c.in})
		if cook != c.want {
			t.Errorf("cookTime = %d, want %d", cook, c.want)
		}
	}
}

func TestExtractMetadata_Servings(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Servings: 12", "12"},
		{"Yield: 1 loaf", "1 loaf"},
		{"Makes 8 rolls", "8 rolls"},
	}
	for _, c := range cases {
		servings, _, _, _ := ExtractMetadata([]string{c.in})
		if servings != c.want {
			t.Errorf("servings = %q, want %q", servings, c.want)
		}
	}
}

func TestExtractMetadata_AdditionalTime(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"Additional Time: 25 mins", 25},
		{"Rise Time: 2 hours", 120},
		{"Rest Time: 30 mins", 30},
	}
	for _, c := range cases {
		_, _, _, additional := ExtractMetadata([]string{c.in})
		if additional != c.want {
			t.Errorf("additionalTime = %d, want %d", additional, c.want)
		}
	}
}

func TestExtractMetadata_Missing(t *testing.T) {
	servings, prep, cook, additional := ExtractMetadata([]string{"500g flour", "200g water"})
	if servings != "" || prep != 0 || cook != 0 || additional != 0 {
		t.Errorf("expected all zero for non-metadata lines, got servings=%q prep=%d cook=%d additional=%d",
			servings, prep, cook, additional)
	}
}
