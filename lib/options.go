package lib

import (
	"log"
	"os"
)

const RecipeHomeVarName = "RECIPE_HOME"

func RecipeHomePath() string {
	path := os.Getenv(RecipeHomeVarName)

	// default to current directory
	if path == "" {

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		return dir
	}

	return path
}
