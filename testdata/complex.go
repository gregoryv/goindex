package testdata

// Decoupled comment

func NewBoat() *Boat {
	return &Boat{}
}

type Boat struct {
	model string
}

// Model returns the model or unknown
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
