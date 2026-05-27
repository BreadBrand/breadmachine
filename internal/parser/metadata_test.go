package parser

import "testing"

func TestExtractMetadata_PrepTime(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Prep Time: 20 mins", "20 mins"},
		{"Prep: 30 minutes", "30 minutes"},
		{"Preparation Time: 1 hour", "1 hour"},
	}
	for _, c := range cases {
		_, prep, _, _ := ExtractMetadata([]string{c.in})
		if prep != c.want {
			t.Errorf("ExtractMetadata(%q) prepTime = %q, want %q", c.in, prep, c.want)
		}
	}
}

func TestExtractMetadata_CookTime(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Cook Time: 15 mins", "15 mins"},
		{"Baking Time: 55 minutes", "55 minutes"},
		{"Bake Time: 30 mins", "30 mins"},
	}
	for _, c := range cases {
		_, _, cook, _ := ExtractMetadata([]string{c.in})
		if cook != c.want {
			t.Errorf("cookTime = %q, want %q", cook, c.want)
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
	cases := []struct{ in, want string }{
		{"Additional Time: 25 mins", "25 mins"},
		{"Rise Time: 2 hours", "2 hours"},
		{"Rest Time: 30 mins", "30 mins"},
	}
	for _, c := range cases {
		_, _, _, additional := ExtractMetadata([]string{c.in})
		if additional != c.want {
			t.Errorf("additionalTime = %q, want %q", additional, c.want)
		}
	}
}

func TestExtractMetadata_Missing(t *testing.T) {
	servings, prep, cook, additional := ExtractMetadata([]string{"500g flour", "200g water"})
	if servings != "" || prep != "" || cook != "" || additional != "" {
		t.Errorf("expected all empty for non-metadata lines")
	}
}
