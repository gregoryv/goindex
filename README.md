goindex - Package for indexing go files

## Quick start

    $ go install github.com/gregoryv/goindex/cmd/goindex@latest
	$ goindex testdata/complex.go
	testdata/complex.go 0 18 package testdata
    testdata/complex.go 18 31 import
    testdata/complex.go 31 54 // Decoupled comment
    testdata/complex.go 54 94 func NewBoat() *Boat 
    testdata/complex.go 94 96 
    testdata/complex.go 96 130 type Boat struct 
    testdata/complex.go 130 132 
    testdata/complex.go 132 280 func (me *Boat) Model() string 
    testdata/complex.go 280 282 
    testdata/complex.go 282 369 func DoSomething(v interface{ X() }) (interface{ S() int }, error) 
    testdata/complex.go 369 392 // Decoupled comment
