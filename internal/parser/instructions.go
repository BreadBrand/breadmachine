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
// instruction lines produced by Stage 2. Each non-blank line is its own step
// because Stage 2 has already segmented the instruction block into one logical
// paragraph per line. Disclaimer lines (e.g. "Keep in mind…", "Your oven may…")
// are excluded wherever they appear, not only when trailing.
func ParseInstructions(lines []string) []string {
	var steps []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Drop disclaimers regardless of position. Stage 2 keeps each disclaimer
		// on its own line, so this filters them without affecting real steps.
		if reDisclaimer.MatchString(line) {
			continue
		}

		if reNumberedStep.MatchString(line) {
			line = reNumberedStep.ReplaceAllString(line, "")
		} else if reBulletStep.MatchString(line) {
			line = reBulletStep.ReplaceAllString(line, "")
		}

		line = strings.TrimSpace(line)
		if line != "" {
			steps = append(steps, line)
		}
	}

	return steps
}
