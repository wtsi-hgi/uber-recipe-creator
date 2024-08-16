package recipe

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/wtsi-hgi/uber-recipe-creator/parser"
)

type Recipe struct {
	Header       string
	Indent       string
	Versions     []Version
	Dependencies []DependsOn
	Footer       string
}

type Version struct {
	Version string
	Extra   map[string]string
}

type DependsOn struct {
	Spec Spec
	Type []string
	When string
}

type Spec struct {
	Name     string
	Version  string
	Variants []string
}

type Header struct {
	ClassName   string
	PackageName string
	Repo        string
	URLType     string
	URLs        []string
}

type Package struct {
	Name    string
	Version string
	Depends []Dependency
	MD5sum  string
}

type Dependency struct {
	Name    string
	Version VersionRange
}

type VersionRange struct {
	Min, Max string
}

func (v VersionRange) String() string {
	var output string
	if v.Min != "" {
		output = v.Min + " <="
	}
	if v.Max != "" || v.Min != "" {
		output += " v "
	}
	if v.Max != "" {
		output += "< " + v.Max
	}
	return output
}

const headerTemplate = "header.tmpl"

var dependencyPattern = regexp.MustCompile(`^([^ (]+) *(\((>=|<=|>|<|==) *([^),]+)(, *(>=|<=|>|<) *([^)]+))?\))?`)

func New(name, repo, urlType string, urls ...string) (*Recipe, error) {
	header := Header{PackageName: name, Repo: repo, URLs: urls, URLType: urlType}
	header.ClassName = "R" + strcase.ToCamel(strings.ReplaceAll(name, "-", " "))
	tmpl, err := template.New(headerTemplate).ParseFiles(headerTemplate)
	if err != nil {
		return nil, err
	}
	result := bytes.NewBufferString("")
	err = tmpl.Execute(result, header)
	if err != nil {
		return nil, err
	}
	return &Recipe{
		Header: result.String(),
		Indent: "\t",
	}, nil
}

func CRANDatabase() (string, error) {
	response, err := http.Get("https://cran.r-project.org/src/contrib/PACKAGES")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var buf strings.Builder
	io.Copy(&buf, response.Body)
	return buf.String(), nil
}

func parseCranDatabase(data string) ([]Package, error) {
	packageList := strings.Split(data, "\n\n")
	var packages []Package
	for _, pkg := range packageList {
		pkg = strings.ReplaceAll(pkg, "\n        ", " ")
		lines := strings.Split(pkg, "\n")
		packageData := make(map[string]string)
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ": ", 2)
			packageData[parts[0]] = parts[1]
		}
		ver := packageData["Version"]
		deps, err := objectifyDependencies(splitString(packageData["Depends"], ", "), splitString(packageData["Imports"], ", "))
		if err != nil {
			return nil, err
		}
		packages = append(packages, Package{
			Name:    packageData["Package"],
			Version: ver,
			Depends: deps,
			MD5sum:  packageData["MD5sum"],
		})
	}
	return packages, nil
}

func (p Package) createRecipe() (*Recipe, error) {
	recipe, err := New(p.Name, "cran", "")
	if err != nil {
		return nil, err
	}
	recipe.Versions = append(recipe.Versions, Version{Version: p.Version})
	for _, dep := range p.Depends {
		ver := dep.versionToSpack()
		recipe.Dependencies = append(recipe.Dependencies, DependsOn{
			Spec: Spec{Name: dep.Name, Version: ver},
			Type: []string{"build", "run"},
		})
	}
	return recipe, nil
}

func objectifyDependencies(depends, imports []string) ([]Dependency, error) {
	stringDeps := append(depends, imports...)
	if len(stringDeps) == 0 {
		return nil, nil
	}
	var deps []Dependency
	for _, stringDep := range stringDeps {
		var dep Dependency
		matches := dependencyPattern.FindStringSubmatch(stringDep)
		if len(matches) == 0 {
			return nil, fmt.Errorf("could not parse dependency: %q", stringDep)
		}
		dep.Name = matches[1]
		if matches[3] != "" {
			dep.setDepVersion(matches[3], matches[4])
		}
		if matches[6] != "" {
			dep.setDepVersion(matches[6], matches[7])
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func (d *Dependency) setDepVersion(cmp, ver string) {
	if strings.HasPrefix(cmp, ">") {
		d.Version.Min = ver
	} else if strings.HasPrefix(cmp, "<") {
		d.Version.Max = ver
	} else if strings.HasPrefix(cmp, "=") {
		d.Version.Min = ver
		d.Version.Max = d.Version.Min
	}
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var lastPos int
	var brackets int
	var skip int
	var parts []string
	for i, char := range s {
		if skip > 0 {
			skip--
			continue
		}
		switch char {
		case '(':
			brackets++
		case ')':
			if brackets == 0 {
				return nil
			}
			brackets--
		default:
			if brackets == 0 && strings.HasPrefix(s[i:], sep) {
				parts = append(parts, s[lastPos:i])
				lastPos = i + len(sep)
				skip = len(sep) - 1
			}
		}
	}
	parts = append(parts, s[lastPos:])
	return parts
}

func (d Dependency) versionToSpack() string {
	if d.Version.Min != "" && d.Version.Max != "" {
		return fmt.Sprintf("@%s:%s", d.Version.Min, d.Version.Max)
	} else if d.Version.Min != "" {
		return fmt.Sprintf("@%s:", d.Version.Min)
	} else if d.Version.Max != "" {
		return fmt.Sprintf("@:%s", d.Version.Max)
	}
	return ""
}

// func parseRecipeFile(r io.Reader) (Recipe, error) {
// 	var recipeBytes []byte
// 	r.Read(recipeBytes)
// 	return parseRecipe(string(recipeBytes))
// }

func parseRecipe(r string) (Recipe, error) {
	var recipe Recipe
	recipeData, err := parser.DoParse(r)
	if err != nil {
		return Recipe{}, err
	}
	recipe.Header = recipeData.Header
	recipe.Indent = recipeData.Indent
	for _, v := range recipeData.Versions {
		version := Version{
			Version: v.Version.Val,
			Extra: map[string]string{
				v.HashType.Val: v.Hash.Val,
			},
		}
		if v.URLType != nil {
			version.Extra[v.URLType.Val] = v.URL.Val
		}
		if v.Preferred != nil {
			version.Extra["preferred"] = v.Preferred.Val
		}
		recipe.Versions = append(recipe.Versions, version)
	}
	for _, d := range recipeData.Depends {
		spec := strings.Split(d.Spec.Val, "@")
		var name, version string
		if len(spec) > 2 {
			return Recipe{}, fmt.Errorf("invalid spec: %q", d.Spec.Val)
		} else if len(spec) > 1 {
			name = spec[0]
			version = "@" + spec[1]
		} else {
			name = spec[0]
			version = ""
		}
		depends := DependsOn{
			Spec: Spec{
				Name:    name,
				Version: version,
			},
		}
		var types []string

		for _, t := range d.Type {
			types = append(types, t.Val)
		}
		depends.Type = types
		if d.When != nil {
			depends.When = d.When.Val
		}
		recipe.Dependencies = append(recipe.Dependencies, depends)
	}
	return recipe, nil
}
