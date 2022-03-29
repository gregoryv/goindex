package testdata

// This is an example file for indexing

func NewCar() *Car {
	return &Car{}
}

type Car struct {
	model string
}

// Model returns the car model
func (me *Car) Model() string {
	if me.model == "" {
		return "unknown"
	}
	// Inline comment
	return me.model
}
