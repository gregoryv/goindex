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

	for i, b := range ParseBlocks(src) {
		fmt.Printf("%v IsConstructor: %5v ", i, b.IsConstructor("X"))
		fmt.Printf("IsType: %5v ", b.IsType("X"))
		fmt.Printf("IsMethod: %5v", b.IsMethod("X"))
		fmt.Println()
	}

	// output:
	// 0 IsConstructor: false IsType: false IsMethod: false
	// 1 IsConstructor: false IsType:  true IsMethod: false
	// 2 IsConstructor: false IsType: false IsMethod: false
	// 3 IsConstructor:  true IsType: false IsMethod: false
	// 4 IsConstructor: false IsType: false IsMethod: false
	// 5 IsConstructor: false IsType: false IsMethod:  true
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
