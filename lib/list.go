package lib

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/jnathanh/recipe/lib/model"
)

func List() (recipes []Recipe, err error) {

	err = filepath.Walk(RecipeHomePath(), func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {

			r, err := Get(path)
			if err != nil {
				return nil
			}

			recipes = append(recipes, *r)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking path: %s", err)
	}

	return
}
