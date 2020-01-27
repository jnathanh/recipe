package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/jnathanh/recipe/lib/model"
	"gopkg.in/russross/blackfriday.v2"
)

func Get(path string) (r *Recipe, err error) {
	r = &Recipe{Path: path}

	// get file bytes
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return r, fmt.Errorf("error reading file '%s': %s", path, err)
	}

	// ensure it is a recipe
	if !isRecipe(bytes) {
		return nil, fmt.Errorf("%s is not a recipe", path)
	}

	// parse to recipe object
	parser := blackfriday.New()
	md := parser.Parse(bytes)
	parsed, err := newRecipeFromBlackFridayAST(*r, md)

	return &parsed, err
}

func isRecipe(bytes []byte) bool {
	trimmed := strings.TrimSpace(string(bytes))

	if trimmed == "" {
		return false
	}

	if trimmed[0] == '#' {
		return true
	}

	return false
}

func newRecipeFromBlackFridayAST(r Recipe, node *blackfriday.Node) (Recipe, error) {
	var ok bool

	r.Name, ok = getName(node)
	if !ok {
		return r, errors.New("could not parse the recipe name")
	}

	return r, nil
}

func getName(recipeDoc *blackfriday.Node) (name string, ok bool) {
	if recipeDoc == nil {
		return
	}

	// get first heading
	firstHeading := recipeDoc.FirstChild
	if firstHeading == nil && firstHeading.Type != blackfriday.Heading {
		return
	}

	// get heading text
	text := firstHeading.FirstChild
	if text == nil {
		return
	}

	return string(text.Literal), true
}
