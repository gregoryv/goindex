package gosort

import (
	"go/scanner"
	"go/token"
	"io"
)

type Block struct {
	src []byte
}

func (me *Block) WriteTo(w io.Writer) {
	w.Write(me.src)
}

func (me *Block) IsConstructor(typeName string) bool {
	c := FindConstructors(me.src, typeName)
	return len(c) > 0
}

func (me *Block) IsMethod(typeName string) bool {
	// todo
	return false
}

func (me *Block) IsType(typeName string) bool {
	for _, t := range FindTypes(me.src) {
		if t == typeName {
			return true
		}
	}
	return false
}

// Index returns an index of the Go src types and funcs
func Index(src []byte) []int {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	var index []int
	var comment int
loop:
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		j := fset.Position(pos).Offset

		switch tok {
		case token.COMMENT:
			comment = j // store for later to include in type or func
			continue loop
		case token.STRING, token.IDENT:
			j = fset.Position(pos).Offset + len(lit)
		case token.TYPE, token.FUNC:
		default:
			j = fset.Position(pos).Offset + len(tok.String())
		}

		//log.Println("i:", i, "j:", j)
		switch tok {
		case token.TYPE, token.FUNC:
			// add start of
			if comment > 0 {
				j = comment
			}
			index = append(index, j)
			index = append(index, moveToEnd(fset, s))
		}
		comment = -1
	}
	return index
}

func moveToEnd(fset *token.FileSet, s scanner.Scanner) int {
	var left int

	for {
		pos, tok, _ := s.Scan()
		switch tok {
		case token.LBRACE:
			left++

		case token.RBRACE:
			left--
			if left == 0 {
				return fset.Position(pos).Offset + len(tok.String())

			}
		}
	}
}

// FindConstructors returns a list of constructor functions for the
// given type name
func FindConstructors(src []byte, typeName string) []string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, 0)

	result := make([]string, 0)
loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break loop
		case token.FUNC:
			_, tok, name := s.Scan()
			if tok == token.LPAREN { // ie. method
				continue loop
			}

			// move to end of arguments
			moveToReturns(s)
			// parse return args
		inner:
			for {
				_, tok, lit := s.Scan()
				switch tok {
				case token.EOF, token.LBRACE:
					break inner
				case token.IDENT:
					if lit == typeName {
						result = append(result, name)
					}
				}
			}
		}
	}
	return result
}

func moveToReturns(s scanner.Scanner) {
	var left int
loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.LPAREN:
			left++

		case token.RPAREN:
			left--
			if left == 0 {
				break loop
			}
		}
	}
}

// FindTypes returns all named types in the given Go source
func FindTypes(src []byte) []string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, 0)

	result := make([]string, 0)
loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break loop
		case token.TYPE:
			_, _, name := s.Scan()
			result = append(result, name)
		}
	}
	return result
}
