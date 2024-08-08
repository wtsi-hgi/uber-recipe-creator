package testdata

import _ "embed"

//go:embed package.txt
var TestRecipe1 string

//go:embed test.txt
var TestScript1 string
