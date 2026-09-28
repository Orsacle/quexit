package shell

import (
  "io"
  "os"
  "os/exec"
  "runtime"
)

type Runner interface {
    Run(line string) error
}

type ExecRunner struct {
    in io.Writer
    out io.Writer
    errOut io.Writer
}

func NewExecRunner(in io.Reader, out, errOut io.Writer) *ExecRunner {
    return &ExecRunner{in: in, out: out, errOut: errOut}
}

func (r *ExecRunner) Run(line string) error {
    name, args := systemShell()
    cmd := exec.Command(name, append(args, line)...)
    cmd.Stdin = r.in
    cmd.Stdout = r.out
    cmd.Stderr = r.errOut
    cmd.Env = os.Environ()
    return cmd.Run()
}

func systemShell() (string, []string) {
    if runtime.GOOS == "windows" {
        return "cmd", []string{"/C"}
    }
    if sh := os.Getenv("SHELL"); sh != "" {
        return sh, []string{"-c"}
    }
  return "/bin/sh", []string{"-c"}
}
