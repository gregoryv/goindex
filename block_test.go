package gosort

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func TestSortBlocks(t *testing.T) {
	src := []byte(`package x

// X stores info
type X struct {
}
func NewX() *X { return &X{} }

// not much here

func (x *X) play() {
_ = "hey"
}

func sum(a, b int) (int, error) { return nil }
func NewRun() {} // not a constructor
`)

	blocks := ParseBlocks(src)
	SortBlocks(blocks[1:], "X")

	var buf bytes.Buffer
	for _, b := range blocks {
		if b.decl == DeclEmpty {
			continue
		}
		b.WriteTo(&buf)
		fmt.Fprint(&buf, "\n\n")
	}
	golden.Assert(t, strings.TrimSpace(buf.String()))
}

func ExampleBlock() {
	src := []byte(`package x

// X stores info
type X struct {
}
func NewX() *X { return &X{} }

// not much here

func (x *X) play() {
_ = "hey"
}

func sum(a, b int) (int, error) { return nil }
func NewRun() {} // not a constructor
`)

	blocks := ParseBlocks(src)
	fmt.Println(blocks[3].String())
	fmt.Println(blocks[5].String())
	fmt.Println(blocks[7].String())
	// output:
	// Constructor X NewX
	// Method X play
	// Func error sum
}

func TestBlock_WriteTo(t *testing.T) {
	src := []byte(`// package x does something
package x

// X stores info
type X struct {
}

// not much here

func play() {
_ = "hey"
}`)
	var buf bytes.Buffer
	for _, b := range ParseBlocks(src) {
		b.WriteTo(&buf)
	}
	if got := buf.String(); got != string(src) {
		t.Error(got)
	}
}

func TestDeclaration(t *testing.T) {
	if got := Declaration(-10).String(); got == "" {
		t.Fail()
	}
}
