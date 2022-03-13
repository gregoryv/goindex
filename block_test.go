package gosort

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleBlock() {
	src := []byte(`package x

// X stores info
type X struct {
}
func NewX() *X { return &X{} }

// not much here

func (x *X) play() {
_ = "hey"
}`)

	blocks := ParseBlocks(src)
	fmt.Println(blocks[3].String())
	// output:
	// X NewX
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
