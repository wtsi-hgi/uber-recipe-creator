package phraser

import (
	_ "embed"
	"testing"
)

//go:embed package.py
var testRecipe string

func TestPhraser(t *testing.T) {
	phrases, err := doPhrase(testRecipe)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if phrases[0].Type != PhraseTop {
		t.Errorf("incorrect phrase type: %q, expected %q", phrases[0].Type.String(), "PhraseTop")
	}
}
