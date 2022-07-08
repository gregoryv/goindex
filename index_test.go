package goindex

import (
	"fmt"
	"os"
	"testing"
)

func xExampleSection_Grab() {
	src := []byte(`package x
// Greet returns a greeting
func Greet() string { return "hello" }`)

	sections := Index(src)
	os.Stdout.Write(sections[1].Grab(src))
	//output:
	// // Greet returns a greeting
	// func Greet() string { return "hello" }
}

func TestIndex(t *testing.T) {
	cases := []struct {
		name   string
		expLen int // number of sections
		src    string
	}{
		{
			expLen: 1,
			src:    `package x`,
		},
		{
			name:   "almost empty",
			expLen: 2,
			src: `package x

// one comment`,
		},
		{
			expLen: 2,
			src: `package x

// String returns ...
func String() string { return "v" }`,
		},

		{
			expLen: 3,
			src: `package x

import _ "embed"

// String returns ...
func String() string { return "v" }`,
		},
		{
			expLen: 4,
			src: `package testdata

import "fmt"

// Decoupled comment
// second line

func NewBoat() *Boat {
        return &Boat{}
}
`,
		},
		{
			expLen: 5,
			src: `package testdata

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
`,
		},
		{
			expLen: 6,
			src: `package testdata

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
func (me *Boat) Model() string {
        if me.model == "" {
                return fmt.Sprintf("%s", "unknown")
        }
        // Inline comment
        return me.model
}
`,
		},
		{
			expLen: 4,
			src: `package testdata

import "fmt"

func DoSomething(v interface{ X() }) (interface{ S() int }, error) {
        return nil, nil
}

// Decoupled comment
`,
		},
	}
	for _, c := range cases { // todo use all
		t.Run(c.name, func(t *testing.T) {
			s := Index([]byte(c.src))
			if got := len(s); got != c.expLen {
				t.Errorf(fmt.Sprintf(`
--------------------------------------------
%s
--------------------------------------------
got %v sections, expected %v
`, c.src, got, c.expLen))
				for i, s := range s {
					t.Logf("%v) %q", i+1, s.Grab([]byte(c.src)))
				}
			}
		})
	}
}
