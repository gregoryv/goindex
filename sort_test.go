package gosort

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func ExampleBlock() {
	src := []byte(`package x

// X stores info
type X struct {
}
func NewX() *X { return &X{} }

// not much here

func play() {
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
	// 5 IsConstructor: false IsType: false IsMethod: false

}

func Example_Index() {
	src := []byte(`package x

// X stores info
type X struct {
}

// not much here

func play() {
_ = "hey"
}`)

	index := Index(src)
	os.Stdout.Write(src[index[0]:index[1]])
	// output:
	// // X stores info
	// type X struct {
	// }
}

func Test_Index(t *testing.T) {
	src := []byte(`// package x does something
package x

// X stores info
type X struct {
}

// not much here

func play() {
_ = "hey"
}`)
	var from int
	var buf bytes.Buffer
	for _, to := range Index(src) {
		fmt.Fprint(&buf, string(src[from:to]))
		from = to
	}
	if got := buf.String(); got != string(src) {
		t.Error(got)
	}
}

func Example_FindConstructors() {
	src := []byte(`package x
type X struct {
}
func play() {}
func New() *X { return &X{} }
func Parse(v []byte) (*X, error) { return &X{}, nil }
func (c *X) Clone(_ interface{ Y() }) *X { t := *c; return &t}
`)

	fmt.Print(FindConstructors(src, "X"))
	// output:
	// [New Parse]
}

func Example_FindTypes() {
	src := []byte(`package x
type X struct {}
func NewX() *X { return &X{} }`)

	fmt.Print(FindTypes(src))
	// output:
	// [X]
}

func init() {
	log.SetFlags(0)
}
