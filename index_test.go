package goindex

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func ExampleIndex() {
	src := []byte(`package x
// Greet returns a greeting
func Greet() string { return "hello" }`)

	for _, s := range Index(src) {
		fmt.Println(s.String())
	}
	// output:
	// 0 10 1 package x
	// 10 76 3 func Greet() string
}

func ExampleSection_Grab() {
	src := []byte(`package x
// Greet returns a greeting
func Greet() string { return "hello" }`)

	sections := Index(src)
	os.Stdout.Write(sections[1].Grab(src))
	//output:
	// // Greet returns a greeting
	// func Greet() string { return "hello" }
}

func TestIndex_generateComment(t *testing.T) {
	src := []byte("package x\n//go:generate ...")
	sections := Index(src)
	if err := grabIs(src, sections); err != nil {
		t.Error(err)
	}
}

func grabIs(src []byte, sections []Section) error {
	var buf bytes.Buffer
	for _, s := range sections {
		buf.Write(s.Grab(src))
	}
	if got := buf.String(); got != string(src) {
		return fmt.Errorf("FAIL\n%s", got)
	}
	return nil
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
			name:   "type func",
			expLen: 3,
			src: `package x

type f func()
type x struct{}
`,
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
