package testdata

// This is an example file for indexing

func NewBoat() *Boat {
	return &Boat{}
}

type Boat struct {
	model string
}

// Model returns the car model
func (me *Boat) Model() string {
	if me.model == "" {
		return "unknown"
	}
	// Inline comment
	return me.model
}

func DoSomething(v interface{ X() }) (interface{ S() int }, error) {
	return nil, nil
}
