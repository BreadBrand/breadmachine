package parser

import (
	"regexp"
	"strings"
)

var (
	reNumberedStep = regexp.MustCompile(`^\d+[\.\)]\s+`)
	reBulletStep   = regexp.MustCompile(`^[-*•]\s+`)
	reDisclaimer   = regexp.MustCompile(`(?i)^(keep in mind|conditions may vary|your oven may|results may vary)`)
)

// ParseInstructions extracts an ordered list of instruction strings from raw
// instruction lines produced by Stage 2. Trailing disclaimer paragraphs are excluded.
func ParseInstructions(lines []string) []string {
	var steps []string
	var pendingParagraph strings.Builder

	flushParagraph := func() {
		s := strings.TrimSpace(pendingParagraph.String())
		if s != "" && !reDisclaimer.MatchString(s) {
			steps = append(steps, s)
		}
		pendingParagraph.Reset()
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			flushParagraph()
			continue
		}

		if reNumberedStep.MatchString(line) {
			flushParagraph()
			step := reNumberedStep.ReplaceAllString(line, "")
			if step != "" {
				steps = append(steps, step)
			}
			continue
		}

		if reBulletStep.MatchString(line) {
			flushParagraph()
			step := reBulletStep.ReplaceAllString(line, "")
			if step != "" {
				steps = append(steps, step)
			}
			continue
		}

		// Sub-header line ending with ":" set by Stage 2 — add as-is
		if strings.HasSuffix(line, ":") && len(strings.Fields(line)) <= 7 {
			flushParagraph()
			steps = append(steps, line)
			continue
		}

		// Accumulate as paragraph
		if pendingParagraph.Len() > 0 {
			pendingParagraph.WriteString(" ")
		}
		pendingParagraph.WriteString(line)
	}
	flushParagraph()

	// Strip trailing disclaimers
	for len(steps) > 0 && reDisclaimer.MatchString(steps[len(steps)-1]) {
		steps = steps[:len(steps)-1]
	}

	return steps
}
