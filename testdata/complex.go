package testdata

import "fmt"

// Decoupled comment
// second line

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
func (b *Boat) Model() string {
	if b.model == "" {
		return fmt.Sprintf("%s", "unknown")
	}
	// Inline comment
	return b.model
}

func DoSomething(v interface{ X() }) (interface{ S() int }, error) {
	return nil, nil
}

// Decoupled comment
