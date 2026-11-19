package handlers

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/BreadBrand/breadmachine/internal/parser"
	"github.com/BreadBrand/breadmachine/models"
)

const (
	maxMetaMinutes    = 7 * 24 * 60
	maxServingsLength = 100
	maxAuthorLength   = parser.MaxAuthorLength
)

func validateMeta(meta models.Meta) error {
	times := []struct {
		name    string
		minutes int
	}{
		{"prepTime", meta.PrepTime},
		{"cookTime", meta.CookTime},
		{"additionalTime", meta.AdditionalTime},
	}
	for _, t := range times {
		if t.minutes < 0 || t.minutes > maxMetaMinutes {
			return fmt.Errorf("%s must be between 0 and %d minutes", t.name, maxMetaMinutes)
		}
	}
	if utf8.RuneCountInString(meta.Servings) > maxServingsLength {
		return errors.New("servings is too long")
	}
	if utf8.RuneCountInString(meta.Author) > maxAuthorLength {
		return errors.New("author is too long")
	}
	return nil
}
