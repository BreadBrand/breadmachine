package parser

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrInputTooLarge = errors.New("INPUT_TOO_LARGE")
	ErrInputEmpty    = errors.New("INPUT_EMPTY")
)

var (
	reHTMLTag           = regexp.MustCompile(`<[^>]+>`)
	reMarkdownBold      = regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`)
	reMarkdownItalic    = regexp.MustCompile(`\*(.+?)\*|_(.+?)_`)
	reMarkdownHeading   = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reMarkdownLink      = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reMarkdownCode      = regexp.MustCompile("`(.+?)`")
	reTrailingAsterisks = regexp.MustCompile(`\*+$`)
	reMixedNumber       = regexp.MustCompile(`\b(\d+)\s+(\d+/\d+)\b`)
	reScalingArtifact   = regexp.MustCompile(`\b\d+x\b`)
	reBrowserHeader     = regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2,4},\s+\d{1,2}:\d{2}\s+(AM|PM)\s+`)
	rePageFraction      = regexp.MustCompile(`^about:blank\s+\d+/\d+$`)
	reMetadataPair      = regexp.MustCompile(`(?i)(Prep Time|Cook Time|Total Time|Additional Time|Servings|Yield|Category|Cuisine|Diet|Author):\s*[^\s]`)
)

var unicodeFractions = map[string]string{
	"½": "1/2", "¼": "1/4", "¾": "3/4",
	"⅓": "1/3", "⅔": "2/3", "⅛": "1/8",
}

var attributionPrefixes = []string{
	"submitted by", "tested by", "reviewed by", "recipe by", "author:",
}

var sourcePrefixes = []string{
	"find it online:", "source:", "originally published at",
	"http://", "https://",
}

var knownSectionHeaders = []string{"ingredients", "directions", "instructions", "method", "steps", "notes"}

var namedFractions = map[string]string{
	"&frac12;": "1/2",
	"&frac14;": "1/4",
	"&frac34;": "3/4",
	"&frac13;": "1/3",
	"&frac23;": "2/3",
}

// Normalise cleans raw recipe text. Returns ErrInputTooLarge if input exceeds
// 10,000 runes. All other operations are best-effort; partial results are preferred.
func Normalise(input string) (string, error) {
	if utf8.RuneCountInString(input) > 10000 {
		return "", ErrInputTooLarge
	}

	if strings.TrimSpace(input) == "" {
		return "", ErrInputEmpty
	}

	// Curly quotes first — all later regex patterns use straight ASCII quotes.
	input = normalizeCurlyQuotes(input)

	// HTML
	input = reHTMLTag.ReplaceAllString(input, "")
	// Replace &nbsp; with a regular space before html.UnescapeString so it
	// doesn't become U+00A0, which strings.TrimSpace would then strip.
	input = strings.ReplaceAll(input, "&nbsp;", " ")
	for entity, replacement := range namedFractions {
		input = strings.ReplaceAll(input, entity, replacement)
	}
	input = html.UnescapeString(input)

	// Markdown
	input = reMarkdownLink.ReplaceAllString(input, "$1")
	input = reMarkdownBold.ReplaceAllString(input, "$1$2")
	input = reMarkdownItalic.ReplaceAllString(input, "$1$2")
	input = reMarkdownHeading.ReplaceAllString(input, "")
	input = reMarkdownCode.ReplaceAllString(input, "$1")

	// Unicode fractions → ASCII
	for uc, ascii := range unicodeFractions {
		input = strings.ReplaceAll(input, uc, ascii)
	}

	// Mixed numbers: "2 3/4" → "2.75" (word-boundary anchored)
	input = reMixedNumber.ReplaceAllStringFunc(input, func(s string) string {
		m := reMixedNumber.FindStringSubmatch(s)
		whole, _ := strconv.ParseFloat(m[1], 64)
		parts := strings.Split(m[2], "/")
		num, _ := strconv.ParseFloat(parts[0], 64)
		den, _ := strconv.ParseFloat(parts[1], 64)
		if den == 0 {
			return s
		}
		return fmt.Sprintf("%g", whole+num/den)
	})

	// Process line by line
	lines := strings.Split(input, "\n")
	var out []string
	inNutritionBlock := false
	consecutiveBlanks := 0

	for _, line := range lines {
		// Trim only trailing whitespace per line; preserve any leading space
		// that may have come from entities like &nbsp;.
		line = strings.TrimRight(line, " \t\r")
		// Strip trailing asterisks (end of line).
		line = reTrailingAsterisks.ReplaceAllString(line, "")

		if line == "" {
			consecutiveBlanks++
			if consecutiveBlanks == 1 {
				out = append(out, "")
			}
			continue
		}
		consecutiveBlanks = 0

		// Use a trimmed version for comparisons; preserve original line for output.
		trimmedLine := strings.TrimSpace(line)
		lower := strings.ToLower(trimmedLine)

		// Browser print artifacts
		if lower == "about:blank" ||
			rePageFraction.MatchString(lower) ||
			reBrowserHeader.MatchString(trimmedLine) {
			continue
		}

		// Nutrition block
		if strings.HasPrefix(lower, "nutrition facts") || strings.HasPrefix(lower, "nutrition information") {
			inNutritionBlock = true
			continue
		}
		if inNutritionBlock {
			// Exit block at known section header
			if isKnownSectionHeader(lower) {
				inNutritionBlock = false
				out = append(out, line)
			}
			continue
		}

		// Attribution and source lines
		isAttribution := false
		for _, prefix := range attributionPrefixes {
			if strings.HasPrefix(lower, prefix) {
				isAttribution = true
				break
			}
		}
		if isAttribution {
			continue
		}

		isSource := false
		for _, prefix := range sourcePrefixes {
			if strings.HasPrefix(lower, prefix) {
				isSource = true
				break
			}
		}
		if isSource {
			continue
		}

		// Strip scaling artifacts adjacent to yield/serving context
		line = reScalingArtifact.ReplaceAllString(line, "")
		line = strings.TrimRight(line, " \t")

		// Expand inline metadata strings: "Prep Time: 20 mins Cook Time: 15 mins" → 2 lines
		expanded := expandInlineMetadata(line)
		out = append(out, expanded...)
	}

	return strings.Trim(strings.Join(out, "\n"), "\n"), nil
}

func normalizeCurlyQuotes(s string) string {
	s = strings.ReplaceAll(s, "‘", "'") // left single
	s = strings.ReplaceAll(s, "’", "'") // right single
	s = strings.ReplaceAll(s, "“", `"`) // left double
	s = strings.ReplaceAll(s, "”", `"`) // right double
	return s
}

func isKnownSectionHeader(lower string) bool {
	for _, kw := range knownSectionHeaders {
		if lower == kw || lower == kw+":" {
			return true
		}
	}
	return false
}

// expandInlineMetadata splits a single line containing multiple "Key: Value" pairs
// (e.g. "Prep Time: 20 mins Cook Time: 15 mins") into separate lines.
func expandInlineMetadata(line string) []string {
	matches := reMetadataPair.FindAllStringIndex(line, -1)
	if len(matches) < 2 {
		return []string{line}
	}
	var result []string
	for i, m := range matches {
		start := m[0]
		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(line)
		}
		part := strings.TrimSpace(line[start:end])
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
