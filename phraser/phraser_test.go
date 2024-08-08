package phraser

import (
	"reflect"
	"testing"

	"github.com/wtsi-hgi/uber-recipe-creator/internal/testdata"
	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

func TestPhraser(t *testing.T) {
	phrases, err := DoPhrase(testdata.TestRecipe1)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if phrases[0].Type != PhraseTop {
		t.Errorf("incorrect phrase type: %q, expected %q", phrases[0].Type.String(), "PhraseTop")
	}

	expected := []Phrase{
		{
			Tokens: []tokeniser.Token{
				{Val: "# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other", Type: tokeniser.TokenComment},
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "# Spack Project Developers. See the top-level COPYRIGHT file for details.", Type: tokeniser.TokenComment},
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "#", Type: tokeniser.TokenComment},
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "# SPDX-License-Identifier: (Apache-2.0 OR MIT)", Type: tokeniser.TokenComment},
				{Val: "\n\n", Type: tokeniser.TokenNewline},
			},
			Type: PhraseTop,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "from", Type: tokeniser.TokenKeyword},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "spack", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "package", Type: tokeniser.TokenIdentifier},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "import", Type: tokeniser.TokenKeyword},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "*", Type: tokeniser.TokenOperator},
			},
			Type: PhraseImport,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n\n", Type: tokeniser.TokenNewline},
				{Val: "class", Type: tokeniser.TokenKeyword},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "Nextdenovo", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\n\t", Type: tokeniser.TokenWhitespace},
				{Val: "MakefilePackage", Type: tokeniser.TokenIdentifier},
				{Val: "\n\t", Type: tokeniser.TokenWhitespace},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ":", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseClass,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "\"\"\"NextDenovo is a string graph-based de novo assembler for long reads.\n\tidk\n\tsomething \n\t\n\t\n\t\n\thello\"\"\"", Type: tokeniser.TokenString},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "homepage", Type: tokeniser.TokenIdentifier},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"https://nextdenovo.readthedocs.io/en/latest/index.html\"", Type: tokeniser.TokenString},
			},
			Type: PhraseHomepage,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "url", Type: tokeniser.TokenIdentifier},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"https://github.com/Nextomics/NextDenovo/archive/refs/tags/2.5.2.tar.gz\"", Type: tokeniser.TokenString},
			},
			Type: PhraseURL,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "version", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"2.5.2\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "sha256", Type: tokeniser.TokenIdentifier},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: "\"f1d07c9c362d850fd737c41e5b5be9d137b1ef3f1aec369dc73c637790611190\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseVersion,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "depends_on", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"python\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "type", Type: tokeniser.TokenIdentifier},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: "\"run\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseDependsOn,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "depends_on", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"py-paralleltask\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "type", Type: tokeniser.TokenIdentifier},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: "\"run\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseDependsOn,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "depends_on", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"zlib\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "type", Type: tokeniser.TokenIdentifier},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"build\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"link\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"run\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseDependsOn,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "def", Type: tokeniser.TokenKeyword},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "edit", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "self", Type: tokeniser.TokenIdentifier},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "spec", Type: tokeniser.TokenIdentifier},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ":", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "makefile", Type: tokeniser.TokenIdentifier},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "FileFilter", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"Makefile\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "makefile", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "filter", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "r\"^TOP_DIR.*\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"TOP_DIR={0}\"", Type: tokeniser.TokenString},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "format", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "self", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "build_directory", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "runfile", Type: tokeniser.TokenIdentifier},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "=", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "FileFilter", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"nextDenovo\"", Type: tokeniser.TokenString},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "runfile", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "filter", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "r\"^SCRIPT_PATH.*\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "\"SCRIPT_PATH = '{0}'\"", Type: tokeniser.TokenString},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "format", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n\n", Type: tokeniser.TokenNewline},
				{Val: "\t", Type: tokeniser.TokenWhitespace},
				{Val: "def", Type: tokeniser.TokenKeyword},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "install", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "self", Type: tokeniser.TokenIdentifier},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "spec", Type: tokeniser.TokenIdentifier},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
				{Val: ":", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "install_tree", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"bin\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "bin", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "install", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"nextDenovo\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "bin", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
		{
			Tokens: []tokeniser.Token{
				{Val: "\n", Type: tokeniser.TokenNewline},
				{Val: "\t\t", Type: tokeniser.TokenWhitespace},
				{Val: "install_tree", Type: tokeniser.TokenIdentifier},
				{Val: "(", Type: tokeniser.TokenDelimiter},
				{Val: "\"lib\"", Type: tokeniser.TokenString},
				{Val: ",", Type: tokeniser.TokenDelimiter},
				{Val: " ", Type: tokeniser.TokenWhitespace},
				{Val: "prefix", Type: tokeniser.TokenIdentifier},
				{Val: ".", Type: tokeniser.TokenDelimiter},
				{Val: "lib", Type: tokeniser.TokenIdentifier},
				{Val: ")", Type: tokeniser.TokenDelimiter},
			},
			Type: PhraseExtra,
		},
	}

	if !reflect.DeepEqual(phrases, expected) {
		for i, phrase := range phrases {
			if !reflect.DeepEqual(phrase, expected[i]) {
				if phrase.Type != expected[i].Type {
					t.Errorf("incorrect phrase type: %s, expected %s", phrase.Type, expected[i].Type)
				} else {
					for j, token := range phrase.Tokens {
						if token != expected[i].Tokens[j] {
							t.Errorf("incorrect token: %q, expected %q", token.Val, expected[i].Tokens[j].Val)
							break
						}
					}
				}
			}
		}
	}
}
