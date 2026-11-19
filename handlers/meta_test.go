package handlers

import (
	"strings"
	"testing"

	"github.com/BreadBrand/breadmachine/models"
)

func TestValidateMeta_Accepts(t *testing.T) {
	cases := map[string]models.Meta{
		"empty":            {},
		"typical":          {Servings: "5", PrepTime: 120, CookTime: 150, AdditionalTime: 30},
		"at the time cap":  {PrepTime: maxMetaMinutes, CookTime: maxMetaMinutes, AdditionalTime: maxMetaMinutes},
		"servings at cap":  {Servings: strings.Repeat("a", maxServingsLength)},
		"servings as text": {Servings: "2 loaves"},
	}
	for name, meta := range cases {
		if err := validateMeta(meta); err != nil {
			t.Errorf("%s: unexpected error: %v", name, err)
		}
	}
}

func TestValidateMeta_Rejects(t *testing.T) {
	cases := map[string]models.Meta{
		"negative prep":       {PrepTime: -1},
		"negative cook":       {CookTime: -30},
		"negative additional": {AdditionalTime: -5},
		"prep over cap":       {PrepTime: maxMetaMinutes + 1},
		"cook over cap":       {CookTime: maxMetaMinutes + 1},
		"additional over cap": {AdditionalTime: maxMetaMinutes + 1},
		"servings too long":   {Servings: strings.Repeat("a", maxServingsLength+1)},
	}
	for name, meta := range cases {
		if err := validateMeta(meta); err == nil {
			t.Errorf("%s: expected an error, got nil", name)
		}
	}
}
