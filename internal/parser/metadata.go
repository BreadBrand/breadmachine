package parser

import (
	"regexp"
	"strconv"
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

var (
	reTimeHours   = regexp.MustCompile(`(?i)(\d+)\s*h(?:r|rs|our|ours)?`)
	reTimeMinutes = regexp.MustCompile(`(?i)(\d+)\s*m(?:in|ins|inute|inutes)?`)
	reTimeRange   = regexp.MustCompile(`(\d+)-\d+`)
)

// parseMinutes converts a raw time string to an integer number of minutes.
// Ranges like "15-17 minutes" use the lower bound. Returns 0 if no
// recognisable time value is found.
func parseMinutes(s string) int {
	s = reTimeRange.ReplaceAllString(s, "$1")
	h, m := 0, 0
	if mh := reTimeHours.FindStringSubmatch(s); mh != nil {
		h, _ = strconv.Atoi(mh[1])
	}
	if mm := reTimeMinutes.FindStringSubmatch(s); mm != nil {
		m, _ = strconv.Atoi(mm[1])
	}
	return h*60 + m
}

// ExtractMetadata parses metadata lines from Stage 2 into individual fields.
// Time fields are normalised to integer minutes; servings is returned as a raw string.
func ExtractMetadata(lines []string) (servings string, prepTime, cookTime, additionalTime int) {
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
				if prepTime == 0 {
					prepTime = parseMinutes(value)
				}
			case "cookTime":
				if cookTime == 0 {
					cookTime = parseMinutes(value)
				}
			case "additionalTime":
				if additionalTime == 0 {
					additionalTime = parseMinutes(value)
				}
			}
			break
		}
	}
	return
}
