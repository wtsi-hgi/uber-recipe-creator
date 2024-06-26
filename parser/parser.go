package parser

import (
	"errors"
	"slices"

	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

type Version struct {
	Version   tokeniser.Token
	HashType  tokeniser.Token
	Hash      tokeniser.Token
	URL       *tokeniser.Token
	Preferred *tokeniser.Token
}

type DependsOn struct {
	Spec tokeniser.Token
	Type []tokeniser.Token
	When *tokeniser.Token
}

type Recipe struct {
	Indent      string
	Imports     [][]tokeniser.Token
	ClassName   tokeniser.Token
	Description tokeniser.Token
	Homepage    tokeniser.Token
	URL         []tokeniser.Token
	Versions    []Version
	DependsOn   []DependsOn

	Extra []tokeniser.Token
}

type Phrase struct {
	Tokens     []tokeniser.Token
	PhraseType int
}

type PhraseFunc func(*Phraser) (Phrase, PhraseFunc)

type Phraser struct {
	tokens  []tokeniser.Token
	pos     int
	lastPos int
}

const (
	PhraseTop = iota
	PhraseImport
	PhraseClass
	PhraseDescription
	PhraseHomepage
	PhraseURL
	PhraseVersions
	PhraseDependsOn
	PhraseExtra
	PhraseDone  = -1
	PhraseError = -2
)

var done = tokeniser.Token{Val: "", Type: tokeniser.TokenDone}

func Parse(input string) (*Recipe, error) {
	tokens, err := tokeniser.Tokenise(input)
	if err != nil {
		return nil, err
	}

	p := Phraser{tokens: tokens}
	var r Recipe

	state := stateStart
loop:
	for {
		var phrase Phrase
		phrase, state = state(&p)

		switch phrase.PhraseType {
		case PhraseTop:
		case PhraseImport:
			r.Imports = append(r.Imports, phrase.Tokens)
		case PhraseClass:
			r.ClassName = phrase.Tokens[2]
		case PhraseDescription:
			r.Description = phrase.Tokens[0]
		case PhraseHomepage:
			r.Homepage = phrase.Tokens[len(phrase.Tokens)-1]
		case PhraseURL:
			r.URL = phrase.Tokens
		case PhraseVersions:
			// r. Versions = append(r.Versions, phrase.Tokens)
		case PhraseDependsOn:

		case PhraseExtra:
			r.Extra = append(r.Extra, phrase.Tokens...)
		case PhraseDone:
			break loop
		case PhraseError:
			return nil, errors.New(phrase.Tokens[0].Val)
		}
	}

	return &r, nil
}

func stateStart(p *Phraser) (Phrase, PhraseFunc) {
	return Phrase{}, nil
}

func (p *Phraser) Next() tokeniser.Token {
	if p.pos >= len(p.tokens) {
		return done
	}
	char := p.tokens[p.pos]
	p.pos++
	return char
}

func (p *Phraser) Accept(types ...tokeniser.TokenType) bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	char := p.tokens[p.pos]
	if !slices.Contains(types, char.Type) {
		return false
	}
	p.pos++
	return true
}

func (p *Phraser) AcceptRun(types ...tokeniser.TokenType) tokeniser.Token {
	for p.Accept(types...) {
	}
	return p.Peek()
}

func (p *Phraser) Except(types ...tokeniser.TokenType) bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	char := p.tokens[p.pos]
	if slices.Contains(types, char.Type) {
		return false
	}
	p.pos++
	return true
}

func (p *Phraser) ExceptRun(types ...tokeniser.TokenType) tokeniser.Token {
	for p.Except(types...) {
	}
	return p.Peek()
}

func (p *Phraser) Peek() tokeniser.Token {
	if p.pos >= len(p.tokens) {
		return done
	}
	char := p.tokens[p.pos]
	return char
}
