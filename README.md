goindex - Package for indexing go files

This module provides tools to manipulate go source files by indexing
and extracting sections of code. This is useful when you are dealing
with large files and want to extract type related sections into a
separate file.

![](goindex.gif)

## Quick start

    $ go install github.com/gregoryv/goindex/cmd/...@latest
	
Index contents of a go file

	$ index complex.go
    complex.go 0 18 package testdata
    complex.go 18 31 import
    complex.go 31 54 // Decoupled comment
    complex.go 54 96 func NewBoat() *Boat
    complex.go 96 132 type Boat struct
    complex.go 132 282 func (me *Boat) Model() string
    complex.go 282 369 func DoSomething(v interface{ X() }) (interface{ S() int }, error)
    complex.go 369 392 // Decoupled comment


then grab Boat related sections using a combination of grep and grab

```shell
$ goindex complex.go | grep Boat | grab
func NewBoat() *Boat {
        return &Boat{}
}

type Boat struct {
        model string
}

// Func comment
func (b *Boat) Model() string {
        if b.model == "" {
                return fmt.Sprintf("%s", "unknown")
        }
        // Inline comment
        return b.model
}
```	
