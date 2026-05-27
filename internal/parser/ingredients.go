package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reQuantityAnchor = regexp.MustCompile(`^(\d+\.?\d*)\s*`)
	reParenthetical  = regexp.MustCompile(`\([^)]*\)`)
	reBulletPrefix   = regexp.MustCompile(`^[-*•—–]\s+`)
	reCrossRef       = regexp.MustCompile(`(?i)\s*(from the build above|see note|recipe follows|from above)\s*`)
)

var noQtyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^(.+?)\s+to\s+taste$`),
	regexp.MustCompile(`(?i)^(.+?)\s+as\s+needed$`),
	regexp.MustCompile(`(?i)^(.+?)\s+for\s+dusting$`),
	regexp.MustCompile(`(?i)^(.+?)\s+for\s+greasing$`),
	regexp.MustCompile(`(?i)^(.+?)\s+to\s+serve$`),
}

var adjectivePrefixes = []string{"bubbly", "active", "ripe", "fresh"}

// doughPhases routes IngredientGroup phases to doughIngredients.
var doughPhases = map[string]bool{
	"dough":       true,
	"":            true,
	"starter build": true,
	"levain":      true,
	"starter":     true,
	"pre-ferment": true,
	"final dough": true,
}

// ParseIngredients parses all IngredientGroups and routes each line to either
// the dough slice (bread-related phases) or the other slice (everything else).
func ParseIngredients(groups []IngredientGroup) (dough []IngredientDTO, other []IngredientDTO) {
	for _, group := range groups {
		for _, line := range group.Lines {
			dto := parseIngredientLine(line)
			if doughPhases[group.Phase] {
				dough = append(dough, dto)
			} else {
				dto.Phase = group.Phase
				other = append(other, dto)
			}
		}
	}
	return
}

func parseIngredientLine(raw string) IngredientDTO {
	dto := IngredientDTO{RawLine: raw}
	line := raw

	// 1. Strip bullet prefix
	line = reBulletPrefix.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	// 2. Yeast alternatives — take first option before "or", keep full RawLine
	if idx := strings.Index(strings.ToLower(line), " or "); idx != -1 {
		before := strings.TrimSpace(line[:idx])
		after := strings.TrimSpace(line[idx+4:])
		if reQuantityAnchor.MatchString(before) && reQuantityAnchor.MatchString(after) {
			line = before
		}
	}

	// 3. No-quantity patterns (e.g. "salt to taste")
	for _, re := range noQtyPatterns {
		m := re.FindStringSubmatch(line)
		if m != nil {
			dto.IngredientName = strings.ToLower(strings.TrimSpace(m[1]))
			dto.ParseOK = dto.IngredientName != ""
			return dto
		}
	}

	// 4. Negative quantity → ParseOK=false
	if strings.HasPrefix(line, "-") && len(line) > 1 && line[1] >= '0' && line[1] <= '9' {
		dto.ParseOK = false
		return dto
	}

	// 5. Extract quantity
	if m := reQuantityAnchor.FindStringSubmatch(line); m != nil {
		if qty, err := strconv.ParseFloat(m[1], 64); err == nil {
			dto.Quantity = qty
			line = strings.TrimSpace(line[len(m[0]):])
		}
	}

	// 6. Extract unit (KnownUnits is ordered multi-word first)
	lower := strings.ToLower(line)
	for _, u := range KnownUnits {
		if lower == u {
			dto.Unit = u
			line = ""
			break
		}
		if strings.HasPrefix(lower, u+" ") || strings.HasPrefix(lower, u+",") {
			dto.Unit = u
			line = strings.TrimSpace(line[len(u):])
			break
		}
	}

	// Strip leading comma left by "u," prefix match
	line = strings.TrimLeft(line, ", \t")

	// 7. Strip parenthetical content (e.g. "(100% hydration)", "(240 grams)")
	line = reParenthetical.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	// 8. Strip adjective prefixes (loop; handle "adj " and "adj," forms)
	for {
		stripped := false
		for _, adj := range adjectivePrefixes {
			lower2 := strings.ToLower(line)
			prefix1 := adj + " "
			prefix2 := adj + ","
			var candidate string
			if strings.HasPrefix(lower2, prefix1) {
				candidate = strings.TrimSpace(line[len(prefix1):])
			} else if strings.HasPrefix(lower2, prefix2) {
				candidate = strings.TrimSpace(line[len(prefix2):])
			} else {
				continue
			}
			// Guard: don't strip if the current line is an unsupported yeast name
			if !IsUnsupportedYeast(line) {
				line = candidate
				stripped = true
			}
			break
		}
		if !stripped {
			break
		}
	}

	// 9. Strip comma-separated notes (rest after first comma has no quantity)
	if idx := strings.Index(line, ","); idx != -1 {
		rest := strings.TrimSpace(line[idx+1:])
		if !reQuantityAnchor.MatchString(rest) {
			line = strings.TrimSpace(line[:idx])
		}
	}

	// 10. Strip cross-reference phrases
	line = reCrossRef.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	dto.IngredientName = line

	// 11. Determine ParseOK: name must be non-empty; zero qty+empty unit only ok
	// for lines matching a no-qty pattern (handled above and returned early).
	dto.ParseOK = dto.IngredientName != "" &&
		!(dto.Quantity == 0 && dto.Unit == "" && !isNoQtyLine(raw))

	return dto
}

func isNoQtyLine(raw string) bool {
	for _, re := range noQtyPatterns {
		if re.MatchString(raw) {
			return true
		}
	}
	return false
}
