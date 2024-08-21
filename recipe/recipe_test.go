package recipe

import (
	"reflect"
	"testing"

	"github.com/wtsi-hgi/uber-recipe-creator/internal/testdata"
)

func TestRecipe(t *testing.T) {
	t.Run("New multiple urls", func(t *testing.T) {
		name := "test-recipe"
		repo := "cran"
		urlType := "urls"
		urls := []string{"https://test.com/test-recipe", "https://test2.com/test-recipe2"}
		r, err := New(name, repo, urlType, urls...)
		if err != nil {
			t.Fatal(err)
		}
		expected := `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class RTestRecipe(RPackage):
	urls = ["https://test.com/test-recipe", "https://test2.com/test-recipe2"]
	cran = "test-recipe"`

		if r.Header != expected {
			t.Fatalf("Header incorrect, expected:\n%q\n\ngot:\n%q", expected, r.Header)
		}
	})

	t.Run("New one url", func(t *testing.T) {
		name := "test-recipe2"
		repo := "bioc"
		urlType := "git"
		urls := []string{"https://test.com/test-recipe"}
		r, err := New(name, repo, urlType, urls...)
		if err != nil {
			t.Fatal(err)
		}
		expected := `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class RTestRecipe2(RPackage):
	git = "https://test.com/test-recipe"
	bioc = "test-recipe2"`

		if r.Header != expected {
			t.Fatalf("Header incorrect, expected %s, got %s", expected, r.Header)
		}
	})
}

func TestCRANDatabase(t *testing.T) {
	t.Run("Fetches and parses CRAN database correctly", func(t *testing.T) {
		db, err := CRANDatabase()
		if err != nil {
			t.Fatal(err)
		}
		_, err = parseCranDatabase(db)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Parses custom database correctly", func(t *testing.T) {
		r := `Package: A3
Version: 1.0-56
Depends: R (>= 2.15.0, < 5.0.0), xtable, pbapply
Suggests: randomForest, e1071
License: GPL (>= 2)
MD5sum: 027ebdd8affce8f0effaecfcd5f5ade2
NeedsCompilation: no

Package: AalenJohansen
Version: v583.1.0
Suggests: knitr, rmarkdown
License: GPL (>= 2)
MD5sum: d7eb2a6275daa6af43bf8a980398b312
NeedsCompilation: no

Package: AATtools
Version: 0.0_1
Depends: R (>= 3.6.0)
Imports: magrittr, dplyr, doParallel, foreach
License: GPL-3
MD5sum: bc59207786e9bc49167fd7d8af246b1c
NeedsCompilation: no

Package: ABACUS
Version: 2024-06-05
Depends: R (< 4.3.0)
Imports: ggplot2 (>= 3.1.0), shiny (>= 1.3.1),
Suggests: rmarkdown (>= 1.13), knitr (>= 1.22)
License: GPL-3
MD5sum: 50c54c4da09307cb95a70aaaa54b9fbd
NeedsCompilation: no

Package: abasequence
Version: 4
License: GPL-3
MD5sum: 1392d909eb0f65be94fd4160a371ae21
NeedsCompilation: no

Package: abbreviate
Version: develop
Suggests: testthat (>= 3.0.0)
License: GPL-3
MD5sum: 37285eddefb6b0fce95783bf21b32999
NeedsCompilation: no
`
		db, err := parseCranDatabase(r)
		if err != nil {
			t.Fatal(err)
		}

		expected := []Package{
			{
				Name:    "A3",
				Version: "1.0-56",
				Depends: []Dependency{
					{
						Name:    "R",
						Version: VersionRange{Min: "2.15.0", Max: "5.0.0"},
					},
					{
						Name:    "xtable",
						Version: VersionRange{Min: "", Max: ""},
					},
					{
						Name:    "pbapply",
						Version: VersionRange{Min: "", Max: ""},
					},
				},
				MD5sum: "027ebdd8affce8f0effaecfcd5f5ade2",
			},
			{
				Name:    "AalenJohansen",
				Version: "v583.1.0",
				Depends: []Dependency(nil),
				MD5sum:  "d7eb2a6275daa6af43bf8a980398b312",
			},
			{
				Name:    "AATtools",
				Version: "0.0_1",
				Depends: []Dependency{
					{
						Name:    "R",
						Version: VersionRange{Min: "3.6.0", Max: ""},
					},
					{
						Name:    "magrittr",
						Version: VersionRange{Min: "", Max: ""},
					},
					{
						Name:    "dplyr",
						Version: VersionRange{Min: "", Max: ""},
					},
					{
						Name:    "doParallel",
						Version: VersionRange{Min: "", Max: ""},
					},
					{
						Name:    "foreach",
						Version: VersionRange{Min: "", Max: ""},
					},
				},
				MD5sum: "bc59207786e9bc49167fd7d8af246b1c",
			},
			{
				Name:    "ABACUS",
				Version: "2024-06-05",
				Depends: []Dependency{
					{
						Name:    "R",
						Version: VersionRange{Min: "", Max: "4.3.0"},
					},
					{
						Name:    "ggplot2",
						Version: VersionRange{Min: "3.1.0", Max: ""},
					},
					{
						Name:    "shiny",
						Version: VersionRange{Min: "1.3.1", Max: ""},
					},
				},
				MD5sum: "50c54c4da09307cb95a70aaaa54b9fbd",
			},
			{
				Name:    "abasequence",
				Version: "4",
				Depends: []Dependency(nil),
				MD5sum:  "1392d909eb0f65be94fd4160a371ae21",
			},
			{
				Name:    "abbreviate",
				Version: "develop",
				Depends: []Dependency(nil),
				MD5sum:  "37285eddefb6b0fce95783bf21b32999",
			},
		}

		if !reflect.DeepEqual(db, expected) {
			for i, p := range db {
				if !reflect.DeepEqual(p, expected[i]) {
					t.Fatalf("Package incorrect, expected \n%+v, got \n%+v", expected[i], p)
				}
			}
		}
	})
}

func TestVersionSplit(t *testing.T) {
	deps, err := objectifyDependencies([]string{"R (>= 3.6.0, < 4.0.0)"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	dep := deps[0]
	verMin := "3.6.0"
	verMax := "4.0.0"
	if !reflect.DeepEqual(dep, Dependency{Name: "R", Version: VersionRange{Min: verMin, Max: verMax}}) {
		t.Fatalf("Dependency incorrect, expected %+v, got %+v", Dependency{Name: "R", Version: VersionRange{Min: verMin, Max: verMax}}, dep)
	}
}

func TestCreateRecipe(t *testing.T) {
	p := Package{
		Name:    "A3",
		Version: "1.0-56",
		Depends: []Dependency{
			{
				Name:    "R",
				Version: VersionRange{Min: "2.15.0", Max: "5.0.0"},
			},
			{
				Name:    "xtable",
				Version: VersionRange{Min: "0.7.5", Max: ""},
			},
			{
				Name:    "pbapply",
				Version: VersionRange{Min: "", Max: "48.1"},
			},
		},
		MD5sum: "027ebdd8affce8f0effaecfcd5f5ade2",
	}
	r, err := p.createRecipe()
	if err != nil {
		t.Fatal(err)
	}

	expected := &Recipe{
		Header: `# Copyright 2013-2023 Lawrence Livermore National Security, LLC and other
# Spack Project Developers. See the top-level COPYRIGHT file for details.
#
# SPDX-License-Identifier: (Apache-2.0 OR MIT)

from spack.package import *


class RA3(RPackage):
	
	cran = "A3"`,
		Indent: "\t",
		Versions: []Version{
			{
				Version: "1.0-56",
				Extra:   map[string]string(nil),
			},
		},
		Dependencies: []DependsOn{
			{
				Spec: Spec{
					Name:     "R",
					Version:  "@2.15.0:5.0.0",
					Variants: []string(nil),
				},
				Type: []string{"build", "run"},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "xtable",
					Version:  "@0.7.5:",
					Variants: []string(nil),
				},
				Type: []string{"build", "run"},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "pbapply",
					Version:  "@:48.1",
					Variants: []string(nil),
				},
				Type: []string{"build", "run"},
				When: "",
			},
		},
		Footer: "",
	}
	if !reflect.DeepEqual(r, expected) {
		if r.Header != expected.Header {
			t.Fatalf("Header incorrect, expected %q, got %q", expected.Header, r.Header)
		}
		for i, v := range r.Versions {
			if !reflect.DeepEqual(v, expected.Versions[i]) {
				t.Fatalf("Version incorrect, expected %+v, got %+v", expected.Versions[i], v)
			}
		}
		for i, d := range r.Dependencies {
			if !reflect.DeepEqual(d, expected.Dependencies[i]) {
				t.Fatalf("Dependency incorrect, expected %+v, got %+v", expected.Dependencies[i], d)
			}
		}
		if r.Footer != expected.Footer {
			t.Fatalf("Footer incorrect, expected %q, got %q", expected.Footer, r.Footer)
		}
	}

}

func TestReadRecipe(t *testing.T) {
	r := testdata.TestCran1
	parsed, err := parseRecipe(r, "abcrf")
	if err != nil {
		t.Fatal(err)
	}

	expected := Recipe{
		Name: "abcrf",
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
				Version: "\"1.9\"",
				Extra: map[string]string{
					"md5": "\"506f4cc36ae9d66bd174f4b65f8c3bb2\"",
				},
			},
		},
		Dependencies: []DependsOn{
			{
				Spec: Spec{
					Name:     "\"r",
					Version:  "@3.1:\"",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-readr\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-mass\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-matrixstats\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-ranger\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-doparallel\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-foreach\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-stringr\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-rcpp\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
			{
				Spec: Spec{
					Name:     "\"r-rcpparmadillo\"",
					Version:  "",
					Variants: []string(nil),
				},
				Type: []string{
					"\"build\"",
					"\"run\"",
				},
				When: "",
			},
		},
		Footer: "",
	}

	if !reflect.DeepEqual(parsed, expected) {
		if parsed.Header != expected.Header {
			t.Fatalf("Header incorrect, expected %q, got %q", expected.Header, parsed.Header)
		}
		if parsed.Indent != expected.Indent {
			t.Fatalf("Indent incorrect, expected %q, got %q", expected.Indent, parsed.Indent)
		}
		for i, v := range parsed.Versions {
			if !reflect.DeepEqual(v, expected.Versions[i]) {
				t.Fatalf("Version incorrect, expected %+v, got %+v", expected.Versions[i], v)
			}
		}
		for i, d := range parsed.Dependencies {
			if !reflect.DeepEqual(d, expected.Dependencies[i]) {
				t.Fatalf("Dependency incorrect, expected %+v, got %+v", expected.Dependencies[i], d)
			}
		}
		if parsed.Footer != expected.Footer {
			t.Fatalf("Footer incorrect, expected %q, got %q", expected.Footer, parsed.Footer)
		}
	}
}
