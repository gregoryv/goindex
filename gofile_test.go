package gosort

import (
	"bytes"
	"os"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_indexFile(t *testing.T) {
	src, err := os.ReadFile("testdata/example.go")
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	ParseGoFile(src).WriteTo(&buf)

	got, exp := buf.String(), string(src)
	if got != exp {
		golden.AssertEquals(t, got, exp)
	}

}

func TestSection_String(t *testing.T) {
	s := Section{
		position: 1,
		ident:    "x",
	}
	if got := s.String(); got == "" {
		t.Fail()
	}
}

func TestSection_IsMethod(t *testing.T) {
	src := []byte(`
type Car struct{}
func (c *Car) Model() {}`)
	sections := ParseSource(src)
	if s := sections[0]; s.IsMethod() {
		t.Error(s.String())
	}
	if s := sections[1]; !s.IsMethod() {
		t.Error(s.String())
	}
}

func TestSection_IsFunc(t *testing.T) {
	src := []byte(`
type Car struct{}
func Model(c *Car) {}`)
	sections := ParseSource(src)
	if s := sections[0]; s.IsFunc() {
		t.Error(s.String())
	}
	if s := sections[1]; !s.IsFunc() {
		t.Error(s.String())
	}
}

func TestSection_ReceiverType(t *testing.T) {
	src := []byte(`
func sum() {}
func (c *Car) Model() {}
func (*Car) Model1() {}
func (Car) Model2() {}
`)
	sections := ParseSource(src)
	if s := sections[0]; s.ReceiverType() != "" {
		t.Error("got ", s.ReceiverType())
	}
	for _, s := range sections[1:] {
		if s.ReceiverType() != "Car" {
			t.Error(s.ReceiverType())
		}
	}
}

func TestSection_FuncName(t *testing.T) {
	src := []byte(`
type Car struct{}
func Model() {}
func (c *Car) Model() {}
func (*Car) Model() {}
func (Car) Model() {}
`)
	sections := ParseSource(src)
	if s := sections[0]; s.FuncName() != "" {
		t.Error("got ", s.FuncName())
	}
	for _, s := range sections[1:] {
		if got := s.FuncName(); got != "Model" {
			t.Error(got)
		}
	}
}

func TestSection_IsComment(t *testing.T) {
	src := []byte(`
// a comment

/* also a
comment
*/
`)
	sections := ParseSource(src)
	if len(sections) != 2 {
		t.Errorf("got %v sections", len(sections))
	}
	for _, s := range sections {
		if !s.IsComment() {
			t.Error(s.String())
		}
	}
}
