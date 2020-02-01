package model

import (
	"net/url"
	"time"
)

type Ingredient struct {
	Name string
}

type Quantity struct {
	Amount float64
	Units Unit
}

// could be a complex type (ex: 4 tortillas, we need a way to represent a simple count
type Unit string

const (
	Minute Unit = "minute"
)

type Portion struct {
	Quantity
	Ingredient
}

type Images struct {
	Featured   *url.URL
	Additional []url.URL
}

type Recipe struct {
	Name            string
	Path            string
	Yield           *Portion
	PreparationTime time.Duration
	CookTime        time.Duration
	TotalTime       time.Duration
	Images          Images
	Source          *url.URL
}
