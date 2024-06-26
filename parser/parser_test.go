package parser

import (
	_ "embed"
	"testing"
)

//go:embed package.py
var testRecipe string

func TestParser(t *testing.T) {
	recipe, err := Parse(testRecipe)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if recipe.Indent != "    " {
		t.Errorf("incorrect indent: %q", recipe.Indent)
	}
}
