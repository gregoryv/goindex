//go:build steps
// +build steps

// steps are used in the pipeline, written in Go to make it more
// platform independent
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func main() {
	tmp := os.TempDir()
	steps := map[string]func(){
		"build": func() {
			sh("go", "build", "-o", path.Join(tmp, "goindex"), "./cmd/goindex")
			sh("go", "build", "-o", path.Join(tmp, "grab"), "./cmd/grab")
		},
		"test": func() {
			sh("go", "test", "./...")
		},
		"dist": func() {
			os.MkdirAll("dist", 0755)
			notes := "dist/release_notes.txt"
			os.WriteFile(notes, releaseNotes(), 0644)
		},
		"clear": func() {
			os.RemoveAll("dist")
		},
	}

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s STEPS\n\nSteps\n", os.Args[0])
		// add steps in prefered order of execution
		fmt.Fprintln(w, "\tbuild")
		fmt.Fprintln(w, "\ttest")
		fmt.Fprintln(w, "\tdist")
		fmt.Fprintln(w, "\tclear")
		flag.PrintDefaults()
	}
	flag.Parse()

	// run all specified target steps
	for _, step := range os.Args[1:] {
		if fn, found := steps[step]; !found {
			fmt.Fprint(os.Stderr, "unknown step:", step)
			os.Exit(1)
		} else {
			fn()
		}
	}
}

func sh(app string, args ...string) {
	c := exec.Command(app, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		os.Exit(1)
	}
}

// reelaseNotes returns last changelog section
func releaseNotes() []byte {
	h2 := []byte("## [")
	from := bytes.Index(changelog, h2)
	to := bytes.Index(changelog[from+len(h2):], h2)
	if to == -1 { // only one
		return changelog[from:]
	}
	return changelog[from : to+from]
}

//go:embed changelog.md
var changelog []byte
