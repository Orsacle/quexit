package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Orsacle/quexit/internal/config"
	"github.com/Orsacle/quexit/internal/shell"
)

func Run(args []string, in io.Reader, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("quexit", flag.ContinueOnError)
	fs.SetOutput(errOut)

	cnfg := fs.String("cnfg", "", "set the string that is treated as exit")
	show := fs.Bool("show", false, "print the currently configured exit string")

	if err := fs.Parse(args); err != nil {
		return err
	}

	store, err := config.NewStore()
	if err != nil {
		return err
	}

	cnfgSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "cnfg" {
			cnfgSet = true
		}
	})

	switch {
	case cnfgSet:
		value := strings.TrimSpace(*cnfg)
		if value == "" {
			return errors.New("exit string must not be empty")
		}
		if err := store.Save(config.Config{ExitCommand: value}); err != nil {
			return err
		}
		fmt.Fprintf(out, "exit command set to %q (%s)\n", value, store.Path())
		return nil

	case *show:
		cfg, err := store.Load()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, cfg.ExitCommand)
		return nil
	}

	cfg, err := store.Load()
	if err != nil {
		return err
	}

	return shell.New(cfg.ExitCommand, in, out, errOut).Loop()
}
