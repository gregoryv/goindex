package goindex

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

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

	t.Run("equals", func(t *testing.T) {
		var buf bytes.Buffer
		for _, s := range sections {
			buf.Write(src[s.From():s.To()])
		}
		got, exp := buf.String(), string(src)
		golden.AssertEquals(t, got, exp)
	})
}
