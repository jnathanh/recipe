package lib

import (
	"log"
	"os"
)

const RecipeHomeVarName = "RECIPE_HOME"

func RecipeHomePath() string {
	path := os.Getenv(RecipeHomeVarName)

	if path == "" {
		log.Panicf("Please set the %s environment variable", RecipeHomeVarName)
	}

	return path
}
