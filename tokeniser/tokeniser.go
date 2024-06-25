package tokeniser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Val  string
	Type int
}

type Tokeniser struct {
	input   string
	pos     int
	lastPos int
}

type tokenFunc func(*Tokeniser) (Token, tokenFunc)

const (
	TokenIdentifier = iota
	TokenNumber
	TokenString
	TokenKeyword
	TokenWhitespace
	TokenNewline
	TokenComment
	TokenDocstring
	TokenOperator
	TokenDone  = -1
	TokenError = -2
)

const whiteSpace = " \t\f"
const newLine = "\n"
const decimal = "0123456789"
const hexadecimal = "0123456789abcdefABCDEF"
const octal = "01234567"
const binary = "01"

var id_start = []*unicode.RangeTable{unicode.Other_ID_Start, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl}
var id_continue = append(id_start, unicode.Other_ID_Continue, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc)

func Tokenise(input string) []Token {
	state := stateStart
	t := Tokeniser{input: input}
	var tokens []Token
	for {
		var token Token
		token, state = state(&t)
		if token.Type < 0 {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
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
	if t.idStart() {
		t.idContinue()
		return Token{t.Get(), TokenIdentifier}, stateStart
	}
	if t.Accept(decimal[1:]) {
		return t.decimalNumber()
	}
	if t.Accept("0") {
		return t.number()
	}
	if t.Accept(".") {
		return t.float()
	}
	if t.Accept("\"") {
		return t.string()
	}

	return stateDone(t)
}

func stateDone(t *Tokeniser) (Token, tokenFunc) {
	return Token{"", TokenDone}, stateDone
}
func (t *Tokeniser) stateError(err string) (Token, tokenFunc) {
	return Token{err, TokenError}, stateDone
}

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
	if !t.Accept(digits) || !t.acceptNumeric(digits) {
		return t.stateError("bad number")
	}
	if digits == "0" && t.Accept(".") {
		return t.float()
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
	if t.Accept("eE") {
		t.Accept("+-")
		if !t.acceptNumeric(decimal) {
			return t.stateError("bad exponent")
		}
	}
	return Token{t.Get(), TokenNumber}, stateStart
}

func (t *Tokeniser) string() (Token, tokenFunc) {
loop:
	for {
		c := t.ExceptRun("\\\n\"")
		fmt.Println(c)
		switch c {
		case '\\':
			t.Next()
			t.Next()
		case '\n':
			return t.stateError("newline in string")
		case '"':
			t.Next()
			break loop
		default:
			return t.stateError("eof")
		}
	}
	return Token{t.Get(), TokenString}, stateStart
}
