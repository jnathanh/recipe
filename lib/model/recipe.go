package model

type Ingredient struct {
	Name string
}

type Quantity struct {
	Amount float64
	Units Unit
}

// could be a complex type (ex: 4 tortillas, we need a way to represent a simple count
type Unit string

type Portion struct {
	Quantity
	Ingredient
}

type Recipe struct {
	Name  string
	Path  string
	Yield *Portion
}
