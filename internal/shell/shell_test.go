package shell

import (
	"bytes"
	"strings"
	"testing"
)

type fakeRunner struct {
	lines []string
}

func (f *fakeRunner) Run(line string) error {
	f.lines = append(f.lines, line)
	return nil
}

func TestIsExit(t *testing.T) {
	cases := []struct {
		line, exit string
		want       bool
	}{
		{"exit", "exit", true},
		{"  bye  ", "bye", true},
		{"exit", "bye", false},
		{"EXIT", "exit", false},
		{"exit now", "exit", false},
	}
	for _, c := range cases {
		if got := IsExit(c.line, c.exit); got != c.want {
			t.Errorf("IsExit(%q, %q) = %v, want %v", c.line, c.exit, got, c.want)
		}
	}
}

func TestLoopStopsOnExitCommand(t *testing.T) {
	in := strings.NewReader("echo a\n\nexit\nbye\necho b\n")
	var out bytes.Buffer
	fr := &fakeRunner{}

	s := New("bye", in, &out, &out)
	s.runner = fr

	if err := s.Loop(); err != nil {
		t.Fatal(err)
	}
	want := []string{"echo a", "exit"}
	if strings.Join(fr.lines, "|") != strings.Join(want, "|") {
		t.Fatalf("ran %v, want %v", fr.lines, want)
	}
}
