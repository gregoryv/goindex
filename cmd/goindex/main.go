// Command goindex indexes a go source file into sections
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gregoryv/goindex"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s FILES...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetFlags(0)

	for _, filename := range flag.Args() {
		src, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		sections := goindex.Index(src)
		for _, s := range sections {
			fmt.Printf("%s %s\n", filename, s.String())
		}
	}
}
