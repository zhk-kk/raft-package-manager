package cmd

import (
	"errors"
	"flag"
	"os"
)

var (
	ErrNoSubcommandProvided    = errors.New("no subcommand provided")
	ErrUnknownSubcommand       = errors.New("unknown subcommand provided")
	ErrArgumentMustBeSpecified = errors.New("argument must be specified, but it wasn't")
	ErrExpectedPath            = errors.New("path was expected, but not received")
)

type Subcommand interface {
	Name() string
	Parse(args []string) error
}

type Resolver struct {
	subcommands []Subcommand
}

func NewResolver(subcommands []Subcommand) *Resolver {
	return &Resolver{subcommands}
}

func (r *Resolver) Parse() error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	if err := fs.Parse(os.Args); err != nil {
		return err
	}

	if err := parseSubcommandsList(fs.Args()[1:], r.subcommands); err != nil {
		return err
	}

	return nil
}

func parseSubcommandsList(args []string, subcommands []Subcommand) error {
	if len(args) == 0 {
		return ErrNoSubcommandProvided
	}
	name := args[0]

	for _, cmd := range subcommands {
		if cmd.Name() == name {
			// Handle nested subcommands.
			if n, ok := cmd.(*nested); ok {
				if err := parseSubcommandsList(args[1:], n.Subcommands); err != nil {
					return err
				}
				return nil
			}

			if err := cmd.Parse(args[1:]); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrUnknownSubcommand
}

type nested struct {
	name        string
	Subcommands []Subcommand
}

func NewNested(name string, subcommands []Subcommand) *nested {
	return &nested{name, subcommands}
}

func (n *nested) Name() string { return n.name }

// Dummy function.
func (*nested) Parse(args []string) error { return nil }
