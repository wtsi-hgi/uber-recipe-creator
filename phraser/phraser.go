package phraser

import (
	"errors"
	"slices"

	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

type Phrase struct {
	Tokens []tokeniser.Token
	Type   PhraseType
}

type PhraseType int

func (p PhraseType) String() string {
	switch p {
	case PhraseTop:
		return "Top"
	case PhraseImport:
		return "Import"
	case PhraseClass:
		return "Class"
	case PhraseDescription:
		return "Description"
	case PhraseHomepage:
		return "Homepage"
	case PhraseURL:
		return "URL"
	case PhraseVersions:
		return "Versions"
	case PhraseDependsOn:
		return "DependsOn"
	case PhraseExtra:
		return "Extra"
	case PhraseDone:
		return "Done"
	case PhraseError:
		return "Error"
	default:
		return "unknown"
	}
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

func doPhrase(input string) ([]Phrase, error) {
	tokens, err := tokeniser.Tokenise(input)
	if err != nil {
		return nil, err
	}

	p := Phraser{tokens: tokens}
	var phrases []Phrase

	state := stateStart
	for {
		var phrase Phrase
		phrase, state = state(&p)

		if phrase.Type == PhraseDone {
			return phrases, nil
		} else if phrase.Type == PhraseError {
			return nil, errors.New(phrase.Tokens[len(phrase.Tokens)-1].Val)
		}
		phrases = append(phrases, phrase)
	}
}

func stateStart(p *Phraser) (Phrase, PhraseFunc) {
	p.AcceptRun(tokeniser.TokenNewline)
	p.Get()
	if p.Accept(tokeniser.TokenComment) {
		p.AcceptRun(tokeniser.TokenComment, tokeniser.TokenNewline)
		return Phrase{p.Get(), PhraseTop}, stateStart
	}
	if c := p.Peek(); p.Accept(tokeniser.TokenKeyword) {
		return p.importOrClass(c)
	}
	return Phrase{}, nil
}

func stateMain(p *Phraser) (Phrase, PhraseFunc) {
	p.AcceptRun(tokeniser.TokenNewline)
	p.Get()
	if p.Accept(tokeniser.TokenComment) {
		p.AcceptRun(tokeniser.TokenComment, tokeniser.TokenNewline)
		return Phrase{p.Get(), PhraseTop}, stateMain
	}
	if c := p.Peek(); p.Accept(tokeniser.TokenIdentifier) {
		return p.identifier(c)
	}

	return Phrase{}, nil
}

func stateDone(p *Phraser) (Phrase, PhraseFunc) {
	return Phrase{[]tokeniser.Token{}, PhraseDone}, stateDone
}

func stateError(p *Phraser) (Phrase, PhraseFunc) {
	return Phrase{p.tokens, PhraseError}, stateDone
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

func (p *Phraser) Get() []tokeniser.Token {
	lastPos := p.lastPos
	p.lastPos = p.pos
	return p.tokens[lastPos:p.pos]
}

func (p *Phraser) importOrClass(c tokeniser.Token) (Phrase, PhraseFunc) {
	p.ExceptRun(tokeniser.TokenNewline)
	if c.Val == "import" || c.Val == "from" {
		return Phrase{p.Get(), PhraseImport}, stateStart
	}
	if c.Val == "class" {
		return Phrase{p.Get(), PhraseClass}, stateMain
	}
	return Phrase{p.Get(), PhraseError}, stateError
}

func (p *Phraser) identifier(c tokeniser.Token) (Phrase, PhraseFunc) {
	p.ExceptRun(tokeniser.TokenNewline)
	if c.Val == "description" {
		return Phrase{p.Get(), PhraseDescription}, stateMain
	}
	if c.Val == "homepage" {
		return Phrase{p.Get(), PhraseHomepage}, stateMain
	}
	if c.Val == "url" {
		return Phrase{p.Get(), PhraseURL}, stateMain
	}
	if c.Val == "versions" {
		return Phrase{p.Get(), PhraseVersions}, stateMain
	}
	if c.Val == "depends_on" {
		return Phrase{p.Get(), PhraseDependsOn}, stateMain
	}
	return Phrase{p.Get(), PhraseError}, stateError
}
