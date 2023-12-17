// Command goto opens file on specific line
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gotoi [OPTIONS] INDEX")
		flag.PrintDefaults()
	}
	filename := flag.String("f", ".index", "goindex file")
	flag.Parse()

	args := flag.Args()

	index := make(map[int]bool)
	for _, a := range args {
		i, err := strconv.Atoi(a)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		index[i] = true
	}
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	in := bufio.NewScanner(bytes.NewReader(data))
	var i int
	for in.Scan() {
		i++
		line := in.Text()
		part := strings.SplitN(line, " ", 5)
		if len(part) < 5 {
			fmt.Fprintln(os.Stderr, line)
			fmt.Fprintln(os.Stderr, "line should be: FILE FROM TO LINE ...")
			continue
		}
		switch {
		case len(args) == 0:
			fmt.Printf("%v %s %s\n", i, part[0], part[4])

		case index[i]:
			exec.Command("emacsclient", "-n", "+"+part[3], part[0]).Run()
		}
	}
}
