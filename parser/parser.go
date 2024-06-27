package parser

import "github.com/wtsi-hgi/uber-recipe-creator/tokeniser"

type Import struct {
	From   tokeniser.Token
	Import tokeniser.Token
}

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
