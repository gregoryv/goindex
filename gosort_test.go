package gosort

import (
	"bytes"
	"os"
	"testing"
)

func ExampleGoSort_Run() {
	cmd := New(os.Stdout, []byte(`package x
type Car struct{}
func NewCar() *Car { return &Car{} }
`))

	_ = cmd.Run()
	// output:
	// package x
	// func NewCar() *Car { return &Car{} }
	// type Car struct{}
}

func TestGosort(t *testing.T) {
	src, err := os.ReadFile("block.go")
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cmd := New(&buf, src)
	SetDebugOutput(os.Stdout)

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}
