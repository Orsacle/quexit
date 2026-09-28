package shell

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const prompt = "quexit> "

type Shell struct {
	exitCommand string
	in          *bufio.Scanner
	out         io.Writer
	errOut      io.Writer
	runner      Runner
}

func New(exitCommand string, in io.Reader, out, errOut io.Writer) *Shell {
	return &Shell{
		exitCommand: exitCommand,
		in:          bufio.NewScanner(in),
		out:         out,
		errOut:      errOut,
		runner:      NewExecRunner(in, out, errOut),
	}
}

func (s *Shell) Loop() error {
	fmt.Fprintf(s.out, "type %q to quit\n", s.exitCommand)

	for {
		fmt.Fprint(s.out, prompt)

		if !s.in.Scan() {
			fmt.Fprintln(s.out)
			return s.in.Err()
		}

		line := strings.TrimSpace(s.in.Text())
		if line == "" {
			continue
		}
		if IsExit(line, s.exitCommand) {
			return nil
		}

		if err := s.runner.Run(line); err != nil {
			fmt.Fprintln(s.errOut, "quexit:", err)
		}
	}
}

func IsExit(line, exitCommand string) bool {
	return strings.TrimSpace(line) == strings.TrimSpace(exitCommand)
}
