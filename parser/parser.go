package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wtsi-hgi/uber-recipe-creator/phraser"
	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

type Recipe struct {
	Header   string
	Indent   string
	Versions []Version
	Depends  []Dependency
	Footer   string
}

type Version struct {
	Version   tokeniser.Token
	HashType  tokeniser.Token
	Hash      tokeniser.Token
	URLType   *tokeniser.Token
	URL       *tokeniser.Token
	Preferred *tokeniser.Token
}

type Dependency struct {
	Spec tokeniser.Token
	Type []tokeniser.Token
	When *tokeniser.Token
}

func DoParse(input string) (*Recipe, error) {
	phrases, err := phraser.DoPhrase(input)
	if err != nil {
		return nil, err
	}

	var header, footer, indent strings.Builder
	var seenVersionOrDepends bool
	var versions []Version
	var depends []Dependency

	for _, phrase := range phrases {
		switch phrase.Type {
		case phraser.PhraseVersion:
			seenVersionOrDepends = true
			setIndent(&phrase, &indent)
			version, err := parseVersion(phrase)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version: %w", err)
			}
			versions = append(versions, version)
		case phraser.PhraseDependsOn:
			seenVersionOrDepends = true
			setIndent(&phrase, &indent)
			dependency, err := parseDependency(phrase)
			if err != nil {
				return nil, fmt.Errorf("failed to parse depends_on: %w", err)
			}
			depends = append(depends, dependency)
		default:
			sb := &header
			if seenVersionOrDepends {
				sb = &footer
			}
			joinTokens(phrase.Tokens, sb)
		}
	}

	recipe := Recipe{
		Header:   header.String(),
		Footer:   footer.String(),
		Indent:   indent.String(),
		Versions: versions,
		Depends:  depends,
	}

	return &recipe, nil
}

func joinTokens(phrase []tokeniser.Token, sb *strings.Builder) {
	for _, token := range phrase {
		sb.WriteString(token.Val)
	}
}

func setIndent(phrase *phraser.Phrase, indent *strings.Builder) {
	for {
		if phrase.Tokens[0].Type != tokeniser.TokenNewline {
			break
		}
		phrase.Tokens = phrase.Tokens[1:]
	}
	var i int
	for i = range phrase.Tokens {
		if phrase.Tokens[i].Type != tokeniser.TokenWhitespace {
			break
		}
	}
	if indent.Len() == 0 {
		joinTokens(phrase.Tokens[:i], indent)
	}

	phrase.Tokens = phrase.Tokens[i:]
}

func parseVersion(phrase phraser.Phrase) (Version, error) {
	var v Version
	p := &phraser.Phraser{Tokens: phrase.Tokens}
	p.Next()
	p.AcceptRun(tokeniser.TokenWhitespace)
	if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "("}) {
		return v, errors.New("expected '('")
	}
	p.AcceptRun(tokeniser.TokenWhitespace)
	v.Version = p.Next()
	if v.Version.Type != tokeniser.TokenString {
		return v, errors.New("expected string")
	}
	for {
		p.AcceptRun(tokeniser.TokenWhitespace)
		nextToken := p.Next()
		if nextToken == (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ")"}) {
			break
		}
		if nextToken != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ","}) {
			return v, errors.New("expected ','")
		}
		p.AcceptRun(tokeniser.TokenWhitespace)
		switch i := p.Next(); i {
		case tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "sha256"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "md5"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "sha1"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "sha224"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "sha384"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "sha512"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "commit"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "tag"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "branch"}:
			v.HashType = i
			p.AcceptRun(tokeniser.TokenWhitespace)
			if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "="}) {
				return v, errors.New("expected '='")
			}
			p.AcceptRun(tokeniser.TokenWhitespace)
			v.Hash = p.Next()
			if v.Hash.Type != tokeniser.TokenString {
				return v, errors.New("expected string")
			}
		case tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "url"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "svn"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "hg"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "cvs"},
			tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "git"}:
			v.URLType = &i
			p.AcceptRun(tokeniser.TokenWhitespace)
			if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "="}) {
				return v, errors.New("expected '='")
			}
			p.AcceptRun(tokeniser.TokenWhitespace)
			url := p.Next()
			if url.Type != tokeniser.TokenString {
				return v, errors.New("expected string")
			}
			v.URL = &url
		case tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "preferred"}:
			p.AcceptRun(tokeniser.TokenWhitespace)
			if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "="}) {
				return v, errors.New("expected '='")
			}
			p.AcceptRun(tokeniser.TokenWhitespace)
			preferred := p.Next()
			if preferred.Type != tokeniser.TokenKeyword {
				return v, errors.New("expected keyword")
			}
			if preferred.Val != "True" && preferred.Val != "False" {
				return v, errors.New("expected 'True' or 'False'")
			}
			v.Preferred = &preferred
		}
	}
	return v, nil
}

func parseDependency(phrase phraser.Phrase) (Dependency, error) {
	var d Dependency
	p := &phraser.Phraser{Tokens: phrase.Tokens}
	p.Next()
	p.AcceptRun(tokeniser.TokenWhitespace)
	if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "("}) {
		return d, errors.New("expected '('")
	}
	p.AcceptRun(tokeniser.TokenWhitespace)
	d.Spec = p.Next()
	if d.Spec.Type != tokeniser.TokenString {
		return d, errors.New("expected string")
	}
	for {
		p.AcceptRun(tokeniser.TokenWhitespace)
		nextToken := p.Next()
		if nextToken == (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ")"}) {
			break
		}
		if nextToken != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ","}) {
			return d, errors.New("expected ','")
		}
		p.AcceptRun(tokeniser.TokenWhitespace)
		switch i := p.Next(); i {
		case tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "type"}:
			p.AcceptRun(tokeniser.TokenWhitespace)
			if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "="}) {
				return d, errors.New("expected '='")
			}
			p.AcceptRun(tokeniser.TokenWhitespace)
			bracketOrType := p.Next()
			if bracketOrType.Type == tokeniser.TokenString {
				d.Type = []tokeniser.Token{bracketOrType}
			} else if bracketOrType == (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "("}) {
				d.Type = append(d.Type, p.Next())
				for {
					if len(d.Type) > 3 {
						return d, errors.New("too many arguments")
					}
					nextToken = p.Next()
					if nextToken == (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ")"}) {
						break
					}
					if nextToken != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: ","}) {
						return d, errors.New("expected ','")
					}
					p.AcceptRun(tokeniser.TokenWhitespace)
					d.Type = append(d.Type, p.Next())
				}
			}
		case tokeniser.Token{Type: tokeniser.TokenIdentifier, Val: "when"}:
			p.AcceptRun(tokeniser.TokenWhitespace)
			if p.Next() != (tokeniser.Token{Type: tokeniser.TokenDelimiter, Val: "="}) {
				return d, errors.New("expected '='")
			}
			p.AcceptRun(tokeniser.TokenWhitespace)
			when := p.Next()
			if when.Type != tokeniser.TokenString {
				return d, errors.New("expected string")
			}
			d.When = &when
		}
	}
	return d, nil
}
