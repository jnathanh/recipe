package lib

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jnathanh/recipe/lib/bfparser"
	. "github.com/jnathanh/recipe/lib/model"
	"gopkg.in/russross/blackfriday.v2"
)

func Get(path string) (r *Recipe, err error) {

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
	parsed, err := bfparser.Parse(md)

	parsed.Path = path

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

