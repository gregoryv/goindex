package gosort

import "os"

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
