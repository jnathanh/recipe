package bfparser_test

import (
	"testing"
	"time"

	"github.com/jnathanh/recipe/lib/bfparser"
	"github.com/stretchr/testify/assert"
)

const testRecipe = `
## Falafel

**yield:** 12 falafel balls

**prep time:** 20 MINUTES

**cook time:** 30 MINUTES

**total time:** 50 MINUTES

<img src="https://cookieandkate.com/images/2018/05/baked-falafel-recipe.jpg" alt="img" style="zoom:25%;" />

**source:** https://cookieandkate.com/crispy-falafel-recipe/#tasty-recipes-22901

### Ingredients:

- ¼ cup + 1 tablespoon extra-virgin olive oil
- 1 cup dried (uncooked/raw) chickpeas, rinsed, picked over and soaked for *at least* 4 hours and up to 24 hours in the refrigerator
- ½ cup roughly chopped red onion (about ½ small red onion)
- ½ cup packed fresh parsley (mostly leaves but small stems are ok)
- ½ cup packed fresh cilantro (mostly leaves but small stems are ok)
- 4 cloves garlic, quartered
- 1 teaspoon fine sea salt
- ½ teaspoon (about 25 twists) freshly ground black pepper
- ½ teaspoon ground cumin
- ¼ teaspoon ground cinnamon

### Directions:

1. With an oven rack in the middle position, preheat oven to 375 degrees Fahrenheit. Pour ¼ cup of the olive oil into a large, rimmed baking sheet and turn until the pan is evenly coated.
2. In a food processor, combine the soaked and drained chickpeas, onion, parsley, cilantro, garlic, salt, pepper, cumin, cinnamon, and the remaining 1 tablespoon of olive oil. Process until smooth, about 1 minute.
3. Using your hands, scoop out about 2 tablespoons of the mixture at a time. Shape the falafel into small patties, about 2 inches wide and ½ inch thick. Place each falafel on your oiled pan.
4. Bake for 25 to 30 minutes, carefully flipping the falafels halfway through baking, until the falafels are deeply golden on both sides. These falafels keep well in the refrigerator for up to 4 days, or in the freezer for several months.

### History

Falafel is a traditionally Arab food. The word falafel may descend from the Arabic word *falāfil,* a plural of the word *filfil,* meaning “pepper.” These fried vegetarian fritters are often served along with [hummus](https://toriavey.com/toris-kitchen/classic-hummus/) and [tahini sauce](https://toriavey.com/toris-kitchen/tahini-sauce/) (known as a “falafel plate.”) They’re also great served with [toum](https://toriavey.com/toris-kitchen/toum-middle-eastern-garlic-sauce/), a Middle Eastern garlic sauce. So just what is the history of this tasty little fritter? According to [The Encyclopedia of Jewish Food](http://www.amazon.com/Encyclopedia-Jewish-Food-Gil-Marks/dp/0470391308?tag=theshiintheki-20) by Gil Marks, “The first known appearance of legume fritters (aka falafel) in the Middle East appears to be in Egypt, where they were made from dried white fava beans (*ful nabed*) and called *tamiya/ta-amia* (from the Arabic for ‘nourishment’); these fritters were a light green color inside. Many attribute *tamiya* to the Copts of Egypt, who practiced one of the earliest forms of Christianity. They believed that the original state of humankind was vegetarian and, therefore, mandated numerous days of eating only vegan food, including *tamiya*.”

source: https://toriavey.com/toris-kitchen/falafel/

### Log
`

func TestParse(t *testing.T) {
	ast := bfparser.AST([]byte(testRecipe))
	w := bfparser.NewRecipeWalker()
	ast.Walk(w.Walk)
	r := w.Recipe()

	is := assert.New(t)

	is.Equal("Falafel", r.Name)
	if is.NotNil(r.Yield) {
		is.Equal("12 falafel balls", r.Yield.Name) // TODO: need some way to distinguish form in the model
	}
	is.NotNil(r.Images.Featured)
	is.Equal(30 * time.Minute, r.CookTime)
	is.Equal(20 * time.Minute, r.PreparationTime)
	is.Equal(50 * time.Minute, r.TotalTime)
	is.Equal("https://cookieandkate.com/crispy-falafel-recipe/#tasty-recipes-22901", r.Source.String())

	for r.Ingredients
}
