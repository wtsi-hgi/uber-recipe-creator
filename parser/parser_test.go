package parser

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/wtsi-hgi/uber-recipe-creator/internal/testdata"
	"github.com/wtsi-hgi/uber-recipe-creator/tokeniser"
)

func TestParser(t *testing.T) {
	for n, test := range [...]struct {
		input       string
		expectation *Recipe
	}{
		{
			input: testdata.TestRecipe1,
			expectation: &Recipe{
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
			},
		},
		{
			input: testdata.TestCran1,
			expectation: &Recipe{
				Header: `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class RAbcrf(RPackage):
	"""Approximate Bayesian Computation via Random Forests

	Performs Approximate Bayesian Computation (ABC) model choice and parameter inference via random forests.
  Pudlo P., Marin J.-M., Estoup A., Cornuet J.-M., Gautier M. and Robert C. P. (2016) <doi:10.1093/bioinformatics/btv684>.
  Estoup A., Raynal L., Verdu P. and Marin J.-M. <http://journal-sfds.fr/article/view/709>.
  Raynal L., Marin J.-M., Pudlo P., Ribatet M., Robert C. P. and Estoup A. (2019) <doi:10.1093/bioinformatics/bty867>.
	"""
	
	cran = "abcrf" `,
				Indent: "\t",
				Versions: []Version{
					{
						Version:  tokeniser.Token{Val: "\"1.9\"", Type: tokeniser.TokenString},
						HashType: tokeniser.Token{Val: "md5", Type: tokeniser.TokenIdentifier},
						Hash:     tokeniser.Token{Val: "\"506f4cc36ae9d66bd174f4b65f8c3bb2\"", Type: tokeniser.TokenString},
					},
				},
				Depends: []Dependency{
					{
						Spec: tokeniser.Token{Val: "\"r@3.1:\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-readr\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-mass\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-matrixstats\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-ranger\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-doparallel\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-foreach\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-stringr\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-rcpp\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-rcpparmadillo\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
				},
				Footer: ``,
			},
		},
		{
			input: testdata.TestBioc1,
			expectation: &Recipe{
				Header: `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class RArraymvout(RPackage):
    """multivariate outlier detection for expression array QA

    This package supports the application of diverse quality metrics to AffyBatch instances, summarizing these metrics via PCA, and then performing parametric outlier detection on the PCs to identify aberrant arrays with a fixed Type I error rate
    """
    
    bioc = "arrayMvout" 
    urls = ["https://www.bioconductor.org/packages/3.18/bioc/src/contrib/arrayMvout_1.60.0.tar.gz", "https://www.bioconductor.org/packages/3.18/bioc/src/contrib/Archive/arrayMvout/arrayMvout_1.60.0.tar.gz"]`,
				Indent: "    ",
				Versions: []Version{
					{
						Version:  tokeniser.Token{Val: "\"1.60.0\"", Type: tokeniser.TokenString},
						HashType: tokeniser.Token{Val: "md5", Type: tokeniser.TokenIdentifier},
						Hash:     tokeniser.Token{Val: "\"7aa46c496dbe47218ea774cb02108800\"", Type: tokeniser.TokenString},
					},
				},
				Depends: []Dependency{
					{
						Spec: tokeniser.Token{Val: "\"r@2.6:\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-parody\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-biobase\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-affy\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-mdqc\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-affycontam\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
					{
						Spec: tokeniser.Token{Val: "\"r-lumi\"", Type: tokeniser.TokenString},
						Type: []tokeniser.Token{
							{Val: "\"build\"", Type: tokeniser.TokenString},
							{Val: "\"run\"", Type: tokeniser.TokenString},
						},
					},
				},
				Footer: ``,
			},
		},
	} {
		recipe, err := DoParse(test.input)
		if err != nil {
			t.Fatalf("Test %d: failed to parse test recipe: %s", n+1, err)
		}

		if !reflect.DeepEqual(recipe, test.expectation) {
			errorPrinted := false
			if recipe.Header != test.expectation.Header {
				t.Errorf("Test %d: header:\n\tgot  %v\n\twant %v", n+1, recipe.Header, test.expectation.Header)
				errorPrinted = true
			}
			if recipe.Indent != test.expectation.Indent {
				t.Errorf("Test %d: indent:\n\tgot  %v\n\twant %v", n+1, recipe.Indent, test.expectation.Indent)
				errorPrinted = true
			}
			for i, version := range recipe.Versions {
				if !reflect.DeepEqual(version, test.expectation.Versions[i]) {
					t.Errorf("Test %d: version %d:\n\tgot  %v\n\twant %v", n+1, i+1, version, test.expectation.Versions[i])
					errorPrinted = true
				}
			}
			for i, dependency := range recipe.Depends {
				if !reflect.DeepEqual(dependency, test.expectation.Depends[i]) {
					t.Errorf("Test %d: dependency %d:\n\tgot  %v\n\twant %v", n+1, i+1, dependency, test.expectation.Depends[i])
					errorPrinted = true
				}
			}
			if recipe.Footer != test.expectation.Footer {
				t.Errorf("Test %d: footer:\n\tgot  %v\n\twant %v", n+1, recipe.Footer, test.expectation.Footer)
				errorPrinted = true
			}
			if !errorPrinted {
				t.Errorf("Test %d: recipe:\n\tgot  %v\n\twant %v", n+1, recipe, test.expectation)
			}
		}
	}
}
