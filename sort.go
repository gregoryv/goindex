package gosort

import (
	"go/scanner"
	"go/token"
)

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
