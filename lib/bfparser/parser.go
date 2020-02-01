package bfparser

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

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
		n.Walk(w.walkSummaryNode)
		return blackfriday.SkipChildren
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

func (w *recipeWalker) walkSummaryNode(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if !entering {
		return blackfriday.GoToNext
	}

	if n.Type == blackfriday.Paragraph {
		text := allChildText(n)

		// featured image
		if strings.Contains(text, "img") {
			r := regexp.MustCompile(`src=["']([^'"]+)["']`)
			matches := r.FindStringSubmatch(text)
			if len(matches) != 2 {
				log.Println("could not parse image tag")
			}
			u, err := url.Parse(matches[1])
			if err != nil {
				log.Println("could not parse image uri")
			}
			w.recipe.Images.Featured = u
			return blackfriday.SkipChildren
		}

		tokens := strings.Split(text, ":")
		if len(tokens) > 2 {
			tokens = []string{tokens[0], strings.Join(tokens[1:], ":")}
		}

		// handle non-key-value lines
		if len(tokens) < 2 {
			log.Printf("Summary line does not match key - value format or a handled bespoke format: '%s'", text)
			return blackfriday.SkipChildren
		}


		key := strings.ToLower(strings.TrimSpace(tokens[0]))
		value := strings.TrimSpace(tokens[1])

		if key == "yield" {
			w.recipe.Yield = &model.Portion{
				Ingredient: model.Ingredient{Name: value},
			}
			return blackfriday.SkipChildren
		}

		if strings.Contains(key, "prep") {
			d, err := parseDuration(value)
			if err != nil {
				log.Printf("Cannot parse value for %s: '%s'; err: %s", key, value, err)
			}
			w.recipe.PreparationTime = d
			return blackfriday.SkipChildren
		}

		if strings.Contains(key, "cook") {
			d, err := parseDuration(value)
			if err != nil {
				log.Printf("Cannot parse value for %s: '%s'; err: %s", key, value, err)
			}
			w.recipe.CookTime = d
			return blackfriday.SkipChildren
		}

		if strings.Contains(key, "total") {
			d, err := parseDuration(value)
			if err != nil {
				log.Printf("Cannot parse value for %s: '%s'; err: %s", key, value, err)
			}
			w.recipe.TotalTime = d
			return blackfriday.SkipChildren
		}

		if strings.Contains(key, "source") {
			u, err := url.Parse(value)
			if err != nil {
				log.Printf("could not parse the 'source' url: %s", value)
			}
			w.recipe.Source = u
			return blackfriday.SkipChildren
		}

		return blackfriday.SkipChildren
	}

	log.Println(n.String())

	return blackfriday.GoToNext
}

func parseDuration(s string) (t time.Duration, err error) {

	// substitutes labels with known time unit labels
	rules := [][]string{
		{"minutes", "m"},
		{"minute", "m"},
		{"min", "m"},
	}

	normalized := strings.ToLower(s)
	for _, set := range rules {
		normalized = strings.ReplaceAll(normalized, set[0], set[1])
	}

	// remove all spaces
	var re = regexp.MustCompile(`\s`)
	normalized = re.ReplaceAllString(normalized, "")

	return time.ParseDuration(normalized)
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

func allChildText(n *blackfriday.Node) (t string) {
	n.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if len(node.Literal) > 0 {
			t = t + string(node.Literal)
		}

		// don't ascend to parent
		if n.Next == nil && n.FirstChild == nil {
			return blackfriday.Terminate
		}

		return blackfriday.GoToNext
	})

	return
}
