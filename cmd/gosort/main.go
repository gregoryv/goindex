// Command gosort sorts go file content in
// - Constructor (funcs starting with New or Parse, returning a public type)
// - Type
// - Methods
// - Funcs
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gregoryv/gosort"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Println("Usage: gosort [OPTIONS] FILENAME")
		fmt.Println("Options")
		flag.PrintDefaults()
	}

	var (
		writeFile = flag.Bool("w", false, "writes result inplace")
	)
	flag.Parse()

	files := flag.Args()
	gosort.SetDebugOutput(os.Stderr)

	for _, filename := range files {
		var dst bytes.Buffer
		cmd := gosort.New(&dst, load(filename))
		_ = cmd.Run()

		if !*writeFile {
			io.Copy(os.Stdout, &dst)
			continue
		}

		os.WriteFile(filename, dst.Bytes(), 0644)
	}
}

func load(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
