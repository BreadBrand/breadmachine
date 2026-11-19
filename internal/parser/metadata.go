package parser

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var metadataPatterns = []struct {
	field   string
	pattern *regexp.Regexp
}{
	{"servings", regexp.MustCompile(`(?i)^(?:Servings?|Serves?)\b\s*:?\s*([^\s:].*)$`)},
	{"servings", regexp.MustCompile(`(?i)^(?:Makes?|Yield)\b\s*:?\s*([^\s:].*)$`)},
	{"prepTime", regexp.MustCompile(`(?i)^(?:Prep(?:aration)?\s*Time?)\b\s*:?\s*([^\s:].*)$`)},
	{"prepTime", regexp.MustCompile(`(?i)^Prep\b\s*:?\s*([^\s:].*)$`)},
	{"cookTime", regexp.MustCompile(`(?i)^(?:Cook(?:ing)?\s*Time?|Bak(?:ing|e)\s*Time?)\b\s*:?\s*([^\s:].*)$`)},
	{"additionalTime", regexp.MustCompile(`(?i)^(?:Additional|Rise|Rest|Chill)\s*Time?\b\s*:?\s*([^\s:].*)$`)},
}

const timeAmount = `(\d+/\d+|\d+(?:\.\d+)?)`

var (
	reTimeHours   = regexp.MustCompile(`(?i)` + timeAmount + `\s*h(?:r|rs|our|ours)?`)
	reTimeMinutes = regexp.MustCompile(`(?i)` + timeAmount + `\s*m(?:in|ins|inute|inutes)?`)
	reTimeRange   = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*[-–]\s*\d+(?:\.\d+)?`)
)

// parseMinutes converts a raw time string to an integer number of minutes.
// Ranges like "15-17 minutes" use the lower bound. Returns 0 if no
// recognisable time value is found.
func parseMinutes(s string) int {
	s = reTimeRange.ReplaceAllString(s, "$1")
	hours, minutes := 0.0, 0.0
	if mh := reTimeHours.FindStringSubmatch(s); mh != nil {
		hours = parseTimeAmount(mh[1])
	}
	if mm := reTimeMinutes.FindStringSubmatch(s); mm != nil {
		minutes = parseTimeAmount(mm[1])
	}
	return int(math.Round(hours*60 + minutes))
}

func parseTimeAmount(token string) float64 {
	if numerator, denominator, isFraction := strings.Cut(token, "/"); isFraction {
		n, _ := strconv.ParseFloat(numerator, 64)
		d, _ := strconv.ParseFloat(denominator, 64)
		if d == 0 {
			return 0
		}
		return n / d
	}
	value, _ := strconv.ParseFloat(token, 64)
	return value
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
