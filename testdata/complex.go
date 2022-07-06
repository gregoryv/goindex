package testdata

import "fmt"

// Decoupled comment

func NewBoat() *Boat {
	return &Boat{}
}

// Type comment
// second line
type Boat struct {
	model string
}

// Func comment
// second line
func (me *Boat) Model() string {
	if me.model == "" {
		return fmt.Sprintf("%s", "unknown")
	}
	// Inline comment
	return me.model
}

func DoSomething(v interface{ X() }) (interface{ S() int }, error) {
	return nil, nil
}

// Decoupled comment
