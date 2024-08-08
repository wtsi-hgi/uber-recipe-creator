package parser

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/wtsi-hgi/uber-recipe-creator/internal/testdata"
	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

func TestParser(t *testing.T) {
	recipe, err := DoParse(testdata.TestRecipe1)
	if err != nil {
		t.Fatalf("Failed to parse test recipe: %s", err)
	}

	expected := &Recipe{
		Header: `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class Nextdenovo(
	MakefilePackage
	):
	"""NextDenovo is a string graph-based de novo assembler for long reads.
	idk
	something 
	
	
	
	hello"""

	homepage = "https://nextdenovo.readthedocs.io/en/latest/index.html"
	url = "https://github.com/Nextomics/NextDenovo/archive/refs/tags/2.5.2.tar.gz"`,
		Indent: "\t",
		Versions: []Version{
			{
				Version:  tokeniser.Token{Val: "\"2.5.2\"", Type: tokeniser.TokenString},
				HashType: tokeniser.Token{Val: "sha256", Type: tokeniser.TokenIdentifier},
				Hash:     tokeniser.Token{Val: "\"f1d07c9c362d850fd737c41e5b5be9d137b1ef3f1aec369dc73c637790611190\"", Type: tokeniser.TokenString},
			},
		},
		Depends: []Dependency{
			{
				Spec: tokeniser.Token{Val: "\"python\"", Type: tokeniser.TokenString},
				Type: []tokeniser.Token{{Val: "\"run\"", Type: tokeniser.TokenString}},
			},
			{
				Spec: tokeniser.Token{Val: "\"py-paralleltask\"", Type: 2},
				Type: []tokeniser.Token{{Val: "\"run\"", Type: tokeniser.TokenString}},
			},
			{
				Spec: tokeniser.Token{Val: "\"zlib\"", Type: tokeniser.TokenString},
				Type: []tokeniser.Token{
					{Val: "\"build\"", Type: tokeniser.TokenString},
					{Val: "\"link\"", Type: tokeniser.TokenString},
					{Val: "\"run\"", Type: tokeniser.TokenString},
				},
			},
		},
		Footer: `

	def edit(self, spec, prefix):
		makefile = FileFilter("Makefile")
		makefile.filter(r"^TOP_DIR.*", "TOP_DIR={0}".format(self.build_directory))
		runfile = FileFilter("nextDenovo")
		runfile.filter(r"^SCRIPT_PATH.*", "SCRIPT_PATH = '{0}'".format(prefix))

	def install(self, spec, prefix):
		install_tree("bin", prefix.bin)
		install("nextDenovo", prefix.bin)
		install_tree("lib", prefix.lib)`,
	}

	if !reflect.DeepEqual(recipe, expected) {
		t.Fatalf("Recipe parsed incorrectly: got %v, expected %v", recipe, expected)
	}
}
