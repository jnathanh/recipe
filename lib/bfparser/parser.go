package bfparser

import (
	"log"
	"regexp"
	"strings"

	"github.com/jnathanh/recipe/lib/model"
	"gopkg.in/russross/blackfriday.v2"
)

func Parse(node *blackfriday.Node) (model.Recipe, error) {
	walker := NewRecipeWalker()

	node.Walk(walker.Walk)

	return *walker.recipe, nil
}

func AST(b []byte) *blackfriday.Node {
	parser := blackfriday.New()

	return parser.Parse(b)
}

type RecipeSection int

const (
	Summary RecipeSection = iota
	Ingredients
	Directions
	History
	Log
)

type recipeWalker struct {
	recipe  *model.Recipe
	section RecipeSection
}

func NewRecipeWalker() recipeWalker {
	return recipeWalker{recipe: &model.Recipe{}}
}

func (w *recipeWalker) Walk(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if !entering {
		return blackfriday.GoToNext
	}

	// get name
	if w.recipe.Name == "" {

		// skip to first heading
		if n.Type != blackfriday.Heading {
			return blackfriday.GoToNext
		}

		w.recipe.Name = nextText(n)
		return blackfriday.SkipChildren
	}

	// detect new section
	if n.Type == blackfriday.Heading {
		w.section = sectionFromHeading(n)
		return blackfriday.SkipChildren
	}

	// section switch
	switch w.section {
	case Summary:
		break
	case Directions:
		break
	case Ingredients:
		break
	case History:
		break
	case Log:
		break
	default:
		log.Fatalf("unhandled section %v", w.section)
	}

	log.Println(n.String())

	return blackfriday.GoToNext
}

func sectionFromHeading(n *blackfriday.Node) RecipeSection {
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}

	t := strings.TrimSpace(strings.ToLower(reg.ReplaceAllString(nextText(n), "")))

	switch t {
	case "directions":
		return Directions
	case "ingredients":
		return Ingredients
	case "history":
		return History
	case "log":
		return Log
	default:
		log.Fatalf("unhandled section description %s", t)
	}

	return RecipeSection(-1)
}

func (w *recipeWalker) Recipe() model.Recipe {
	return *w.recipe
}

func nextText(n *blackfriday.Node) (t string) {
	n.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if len(node.Literal) > 0 {
			t = string(node.Literal)
			return blackfriday.Terminate
		}
		return blackfriday.GoToNext
	})

	return
}
