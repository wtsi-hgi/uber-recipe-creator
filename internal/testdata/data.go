package testdata

import _ "embed"

//go:embed package.txt
var TestRecipe1 string

//go:embed test.txt
var TestScript1 string

//go:embed cran.txt
var TestCran1 string

//go:embed bioc.txt
var TestBioc1 string

//go:embed PACKAGES
var TestPackageDB string
