goindex - Package for indexing go files

This module provides tools to manipulate go source files by indexing
and extracting sections of code.


## Quick start

    $ go install github.com/gregoryv/goindex/cmd/goindex@latest
    $ go install github.com/gregoryv/goindex/cmd/gograb@latest	
	
Index contents of a go file

	$ goindex complex.go
    complex.go 0 18 package testdata
    complex.go 18 31 import
    complex.go 31 54 // Decoupled comment
    complex.go 54 96 func NewBoat() *Boat
    complex.go 96 132 type Boat struct
    complex.go 132 282 func (me *Boat) Model() string
    complex.go 282 369 func DoSomething(v interface{ X() }) (interface{ S() int }, error)
    complex.go 369 392 // Decoupled comment


Grab Boat related sections

```shell
$ goindex complex.go | grep Boat | gograb
func NewBoat() *Boat {
        return &Boat{}
}

type Boat struct {
        model string
}

// Func comment
func (me *Boat) Model() string {
        if me.model == "" {
                return fmt.Sprintf("%s", "unknown")
        }
        // Inline comment
        return me.model
}
```	
