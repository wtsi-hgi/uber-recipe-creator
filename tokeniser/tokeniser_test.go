package tokeniser

import (
	"reflect"
	"testing"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Tokens []Token
	}{
		{"", nil},
		{"a", []Token{{"a", TokenIdentifier}}},
		{"abc", []Token{{"abc", TokenIdentifier}}},
		{"077e010", []Token{{"077e010", TokenNumber}}},
		{"\"hello, world\"", []Token{{"\"hello, world\"", TokenString}}},
	} {
		tokens := Tokenise(test.Input)
		if !reflect.DeepEqual(tokens, test.Tokens) {
			t.Errorf("Test %d: got %v, want %v", n, tokens, test.Tokens)
		}
	}
}
