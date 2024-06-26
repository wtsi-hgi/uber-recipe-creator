package tokeniser

import (
	"errors"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Val  string
	Type tokenType
}

type Tokeniser struct {
	input   string
	pos     int
	lastPos int
}

type tokenFunc func(*Tokeniser) (Token, tokenFunc)

type tokenType int

func (t tokenType) String() string {
	switch t {
	case TokenIdentifier:
		return "Identifier"
	case TokenNumber:
		return "Number"
	case TokenString:
		return "String"
	case TokenKeyword:
		return "Keyword"
	case TokenWhitespace:
		return "Whitespace"
	case TokenNewline:
		return "Newline"
	case TokenComment:
		return "Comment"
	case TokenOperator:
		return "Operator"
	case TokenDelimiter:
		return "Delimiter"
	case TokenDone:
		return "Done"
	case TokenError:
		return "Error"
	default:
		return "unknown"
	}
}

const (
	TokenIdentifier tokenType = iota
	TokenNumber
	TokenString
	TokenKeyword
	TokenWhitespace
	TokenNewline
	TokenComment
	TokenOperator
	TokenDelimiter
	TokenDone  tokenType = -1
	TokenError tokenType = -2
)

const whiteSpace = " \t\f"
const newLine = "\n"
const decimal = "0123456789"
const hexadecimal = "0123456789abcdefABCDEF"
const octal = "01234567"
const binary = "01"

var keywords = [...]string{"False", "await", "else", "import", "pass", "None", "break", "except", "in", "raise", "True", "class", "finally", "is", "return", "and", "continue", "for", "lambda", "try", "as", "def", "from", "nonlocal", "while", "assert", "del", "global", "not", "with", "async", "elif", "if", "or", "yield"}

var id_start = []*unicode.RangeTable{unicode.Other_ID_Start, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl}
var id_continue = append(id_start, unicode.Other_ID_Continue, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc)

func Tokenise(input string) ([]Token, error) {
	state := stateStart
	t := Tokeniser{input: input}
	var tokens []Token
	for {
		var token Token
		token, state = state(&t)
		if token.Type == TokenDone {
			return tokens, nil
		} else if token.Type == TokenError {
			return nil, errors.New(token.Val)
		}
		tokens = append(tokens, token)
	}
}

func stateStart(t *Tokeniser) (Token, tokenFunc) {
	if t.Accept(whiteSpace) {
		t.AcceptRun(whiteSpace)
		return Token{t.Get(), TokenWhitespace}, stateStart
	}
	if t.Accept(newLine) {
		t.AcceptRun(newLine)
		return Token{t.Get(), TokenNewline}, stateStart
	}
	if c := t.Peek(); t.Accept("rRuUfFbB") {
		return t.possibleString(c)
	}
	if t.idStart() {
		return t.identifier()
	}
	if t.Accept("0") {
		return t.number()
	}
	if t.Accept(decimal) {
		return t.decimalNumber()
	}
	if t.Accept(".") {
		if strings.ContainsRune(decimal, t.Peek()) {
			return t.float()
		}
		return Token{t.Get(), TokenDelimiter}, stateStart
	}
	if c := t.Peek(); t.Accept("\"'") {
		return t.string(c)
	}
	if t.Accept("#") {
		t.ExceptRun("\n")
		return Token{t.Get(), TokenComment}, stateStart
	}

	return t.operator()
}

func stateDone(t *Tokeniser) (Token, tokenFunc) {
	return Token{"", TokenDone}, stateDone
}
func (t *Tokeniser) stateError(err string) (Token, tokenFunc) {
	return Token{err, TokenError}, stateDone
}

// Accept consumes the next rune if it's in the input string.
// Input string is a set of (unordered) characters that are accepted.
func (t *Tokeniser) Accept(input string) bool {
	if t.pos >= len(t.input) {
		return false
	}
	char, size := utf8.DecodeRuneInString(t.input[t.pos:])
	if !strings.ContainsRune(input, char) {
		return false
	}
	t.pos += size
	return true
}

func (t *Tokeniser) AcceptRun(input string) rune {
	for t.Accept(input) {
	}
	return t.Peek()
}

func (t *Tokeniser) Except(input string) bool {
	if t.pos >= len(t.input) {
		return false
	}
	char, size := utf8.DecodeRuneInString(t.input[t.pos:])
	if strings.ContainsRune(input, char) {
		return false
	}
	t.pos += size
	return true
}

func (t *Tokeniser) ExceptRun(input string) rune {
	for t.Except(input) {
	}
	return t.Peek()
}

func (t *Tokeniser) Peek() rune {
	if t.pos >= len(t.input) {
		return 0
	}
	char, _ := utf8.DecodeRuneInString(t.input[t.pos:])
	return char
}

func (t *Tokeniser) Get() string {
	lastPos := t.lastPos
	t.lastPos = t.pos
	return t.input[lastPos:t.pos]
}

func (t *Tokeniser) Next() rune {
	if t.pos >= len(t.input) {
		return 0
	}
	char, size := utf8.DecodeRuneInString(t.input[t.pos:])
	t.pos += size
	return char
}

func (t *Tokeniser) idStart() bool {
	char := t.Peek()

	if !unicode.In(char, id_start...) {
		return false
	}
	t.Next()
	return true
}

func (t *Tokeniser) idContinue() {
	for {
		char := t.Peek()
		if !unicode.In(char, id_continue...) {
			break
		}
		t.Next()
	}
}

func (t *Tokeniser) decimalNumber() (Token, tokenFunc) {
	if !t.acceptNumeric(decimal) {
		return t.stateError("bad number")
	}
	if t.Accept(".") {
		return t.float()
	}
	if !t.exponent() {
		return t.stateError("invalid exponent")
	}
	return Token{t.Get(), TokenNumber}, stateStart
}

func (t *Tokeniser) number() (Token, tokenFunc) {
	digits := "0"
	if t.Accept("bB") {
		digits = binary
	} else if t.Accept("oO") {
		digits = octal
	} else if t.Accept("xX") {
		digits = hexadecimal
	}
	t.Accept("_")
	if (!t.Accept(digits) && digits != "0") || !t.acceptNumeric(digits) {
		return t.stateError("bad number")
	}
	if digits == "0" {
		if t.Accept(".") {
			return t.float()
		}
		if t.acceptNumeric(decimal) {
			switch t.Peek() {
			case 'e', 'E':
				t.exponent()
			}
		}
	}

	return Token{t.Get(), TokenNumber}, stateStart
}

func (t *Tokeniser) acceptNumeric(digits string) bool {
	t.AcceptRun(digits)
	for t.Accept("_") {
		if !t.Accept(digits) {
			return false
		}
		t.AcceptRun(digits)
	}
	return true
}

func (t *Tokeniser) float() (Token, tokenFunc) {
	if !t.acceptNumeric(decimal) {
		return t.stateError("bad float")
	}
	if !t.exponent() {
		return t.stateError("bad exponent")
	}
	return Token{t.Get(), TokenNumber}, stateStart
}

func (t *Tokeniser) exponent() bool {
	if t.Accept("eE") {
		t.Accept("+-")
		if !t.acceptNumeric(decimal) {
			return false
		}
	}
	t.Accept("jJ")
	return true
}

func (t *Tokeniser) operator() (Token, tokenFunc) {
	var tokenType = TokenOperator
	switch c := t.Peek(); c {
	case 0:
		return stateDone(t)
	case '+', '&', '|', '^', '%', '@':
		t.Next()
		if t.Accept("=") {
			tokenType = TokenDelimiter
		}
	case '-':
		t.Next()
		if t.Accept(">=") {
			tokenType = TokenDelimiter
		}
	case '*', '/', '<', '>':
		t.Next()
		db := t.Accept(string(c))
		if t.Accept("=") && (db || c == '*' || c == '/') {
			tokenType = TokenDelimiter
		}
	case ':', '=':
		t.Next()
		if !t.Accept("=") {
			tokenType = TokenDelimiter
		}
	case '!':
		t.Next()
		if !t.Accept("=") {
			return t.stateError("invalid operator")
		}
	case '~':
		t.Next()

	default:
		if !t.Accept("()[]{},.;") {
			return t.stateError("invalid delimiter")
		}

		tokenType = TokenDelimiter
	}
	return Token{t.Get(), tokenType}, stateStart
}

func (t *Tokeniser) string(d rune) (Token, tokenFunc) {
	delim := string(d)
	long := false
	if t.Accept(delim) {
		if t.Accept(delim) {
			long = true
		} else {
			return Token{t.Get(), TokenString}, stateStart
		}
	}
	checked := "\\\n" + delim
loop:
	for {
		c := t.ExceptRun(checked)
		switch c {
		case '\\':
			t.Next()
			t.Next()
		case '\n':
			if !long {
				return t.stateError("newline in string")
			}
			t.Next()
		case d:
			t.Next()
			if long && (!t.Accept(delim) || !t.Accept(delim)) {
				continue
			}
			break loop
		default:
			return t.stateError("eof")
		}
	}
	return Token{t.Get(), TokenString}, stateStart
}

func (t *Tokeniser) identifier() (Token, tokenFunc) {
	t.idContinue()
	identifier := t.Get()
	tokenType := TokenIdentifier
	if slices.Contains(keywords[:], identifier) {
		tokenType = TokenKeyword
	}
	return Token{identifier, tokenType}, stateStart
}

func (t *Tokeniser) possibleString(c rune) (Token, tokenFunc) {
	switch c {
	case 'r', 'R':
		t.Accept("bBfF")
	case 'u', 'U':
	case 'f', 'F':
		t.Accept("rR")
	case 'b', 'B':
		t.Accept("rR")
	}
	if c := t.Peek(); t.Accept("\"'") {
		return t.string(c)
	}
	return t.identifier()
}
