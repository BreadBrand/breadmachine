package parser

import (
	"regexp"
	"strings"
)

var metadataPatterns = []struct {
	field   string
	pattern *regexp.Regexp
}{
	{"servings", regexp.MustCompile(`(?i)^(?:Servings?|Serves?)\s*:?\s*(.+)$`)},
	{"servings", regexp.MustCompile(`(?i)^(?:Makes?|Yield)\s*:?\s*(.+)$`)},
	{"prepTime", regexp.MustCompile(`(?i)^(?:Prep(?:aration)?\s*Time?)\s*:?\s*(.+)$`)},
	{"prepTime", regexp.MustCompile(`(?i)^Prep\s*:?\s*(.+)$`)},
	{"cookTime", regexp.MustCompile(`(?i)^(?:Cook(?:ing)?\s*Time?|Bak(?:ing|e)\s*Time?)\s*:?\s*(.+)$`)},
	{"additionalTime", regexp.MustCompile(`(?i)^(?:Additional|Rise|Rest|Chill)\s*Time?\s*:?\s*(.+)$`)},
}

// ExtractMetadata parses metadata lines from Stage 2 into individual fields.
// All fields are returned as raw strings; no duration normalisation is applied.
func ExtractMetadata(lines []string) (servings, prepTime, cookTime, additionalTime string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for _, p := range metadataPatterns {
			m := p.pattern.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			value := strings.TrimSpace(m[1])
			switch p.field {
			case "servings":
				if servings == "" {
					servings = value
				}
			case "prepTime":
				if prepTime == "" {
					prepTime = value
				}
			case "cookTime":
				if cookTime == "" {
					cookTime = value
				}
			case "additionalTime":
				if additionalTime == "" {
					additionalTime = value
				}
			}
			break
		}
	}
	return
}
