package lib

import (
	"fmt"
	"io/ioutil"

	. "github.com/jnathanh/recipe/lib/model"
	"gopkg.in/russross/blackfriday.v2"
)

func Get(path string) (r Recipe, err error) {
	r.Path = path

	// get file bytes
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return r, fmt.Errorf("error reading file '%s': %s", path, err)
	}

	// parse to recipe object
	parser := blackfriday.New()
	md := parser.Parse(bytes)
	r = newRecipeFromBlackFridayAST(r, md)

	return
}

func newRecipeFromBlackFridayAST(r Recipe, node *blackfriday.Node) Recipe {
	// heading > text
	r.Name = string(node.FirstChild.FirstChild.Literal)

	return r
}
