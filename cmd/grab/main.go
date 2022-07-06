// Command grab extracts sections from a file by byte range
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s\n", os.Args[0])
		fmt.Fprint(w, `
grab command reads input i the form

FILE1 FROM TO
FILE1 FROM TO
FILE2 FROM TO

and writes out the sections from those files to stdout.

FROM and TO are the byte index in each file.
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetFlags(0)

	s := bufio.NewScanner(os.Stdin)
	var g Grabber
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) < 3 {
			log.Println(line)
			log.Fatal("line should be: FILE FROM TO")
		}
		if err := g.Grab(os.Stdout, f[0], f[1], f[2]); err != nil {
			log.Fatal(err)
		}
	}
}

type Grabber struct {
	filename string
	src      []byte
}

func (me *Grabber) Grab(w io.Writer, filename, sfrom, sto string) error {
	if me.filename != filename {
		src, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		me.src = src
		me.filename = filename
	}
	from, err := strconv.Atoi(sfrom)
	if err != nil {
		return err
	}
	to, err := strconv.Atoi(sto)
	if err != nil {
		return err
	}
	w.Write(me.src[from:to])
	return nil
}
