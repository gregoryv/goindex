// Command goindex indexes a go source file into sections
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gregoryv/goindex"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s FILES...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	for _, filename := range flag.Args() {
		src, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		sections := goindex.Index(src)
		for _, s := range sections {
			fmt.Printf("%s %s\n", filename, s.String())
		}
	}
}
