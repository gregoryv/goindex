package goindex

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func ExampleSection_Grab() {
	src := []byte(`package x
// Greet returns a greeting
func Greet() string { return "hello" }`)

	sections := Index(src)
	os.Stdout.Write(sections[1].Grab(src))
	//output:
	// // Greet returns a greeting
	// func Greet() string { return "hello" }
}

func TestIndex(t *testing.T) {
	src, err := os.ReadFile("testdata/complex.go")
	if err != nil {
		t.Fatal(err)
	}
	sections := Index(src)

	t.Run("decoupled comment", func(t *testing.T) {
		for _, s := range sections {
			var buf bytes.Buffer
			buf.Write(src[s.From():s.To()])
			got := buf.String()
			if strings.Contains(got, "Decoupled comment") {
				if strings.Contains(got, "func") || strings.Contains(got, "type") {
					t.Log(got)
					t.Error("contains unrelated comment")
				}
			}
		}
	})

	t.Run("related comment", func(t *testing.T) {
		var buf bytes.Buffer
		for _, s := range sections {
			buf.Write(src[s.From():s.To()])
		}
		got := buf.String()
		if !strings.Contains(got, "Func comment") {
			t.Log(got)
			t.Error("missing related comment")
		}
	})

	t.Run("starts from 0", func(t *testing.T) {
		if got := sections[0].From(); got != 0 {
			t.Error(got)
		}
	})

	t.Run("ends with src len", func(t *testing.T) {
		exp := len(src)
		if got := sections[len(sections)-1].To(); got != exp {
			t.Error(got, "expected", exp)
		}
	})

	t.Run("sections are complete", func(t *testing.T) {
		for i, s := range sections[:len(sections)-1] {
			next := sections[i+1]
			t.Log(s.To(), next.From())
			if s.To() != next.From() {
				t.Errorf("missing section between %v and %v", i, i+1)
			}
		}
	})
	t.Run("multiline type comment", func(t *testing.T) {
		s := sections[4]
		got := string(s.Grab(src))
		if !strings.Contains(got, "second line") {
			t.Log(s.String())
			t.Error("missing second line")
		}
	})
	t.Run("multiline func comment", func(t *testing.T) {
		s := sections[6]
		got := string(s.Grab(src))
		if !strings.Contains(got, "second line") {
			t.Log(s.String())
			t.Error("missing second line")
		}
	})

	t.Run("equals", func(t *testing.T) {
		var buf bytes.Buffer
		for _, s := range sections {
			buf.Write(src[s.From():s.To()])
		}
		got, exp := buf.String(), string(src)
		golden.AssertEquals(t, got, exp)
	})
}
