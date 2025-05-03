package gomonorepo

import (
	"context"

	"github.com/jessevdk/go-flags"
)

type Commander interface {
	RunCommand(ctx context.Context, options *GlobalOptions) error
}

// An EmbeddedCommand should be embedded in a struct to provide the Command interface.
// This allows for a simple way to define a command with a command name,
// short and long descriptions, without requiring the implementer to
// define the Command interface methods.
type EmbeddedCommand struct {
	CmdName string
	Aliases []string // (optional) add any number of additional command aliases.
	Short   string
	Long    string
}

// Execute is a no-op implementation of the Commander interface.
func (x *EmbeddedCommand) Execute([]string) error {
	return nil
}

// CommandName returns the name of the command.
// This interface exists to automatically build commands.
func (x *EmbeddedCommand) CommandName() string {
	return x.CmdName
}

// GetAliases returns the aliases of the command.
// This interface exists to automatically build commands.
func (x *EmbeddedCommand) GetAliases() []string {
	return x.Aliases
}

// ShortDesc returns the short description of the command.
// This interface exists to automatically build commands.
func (x *EmbeddedCommand) ShortDesc() string {
	return x.Short
}

// LongDesc returns the long description of the command.
// If the Long description is not set, it will return Short description.
// This interface exists to automatically build commands.
func (x *EmbeddedCommand) LongDesc() string {
	if x.Long != "" {
		return x.Long
	}
	return x.Short
}

// CommandSubGroup is a convenience struct that can be used to define a group of
// sub-commands that are related to each other.
type CommandSubGroup struct {
	EmbeddedCommand
	Commands []Command
}

// SubCommands returns the sub-commands of the group as a list.
func (x *CommandSubGroup) SubCommands() []Command {
	return x.Commands
}

// Command is an interface that represents a command that can be executed.
// This interface is used to automatically build commands.
type Command interface {
	Commander
	flags.Commander
	CommandName() string
	GetAliases() []string
	ShortDesc() string
	LongDesc() string
}

// subCommander is an interface that represents a command with sub-commands.
type subCommander interface {
	SubCommands() []Command
}

// AddCommander is an interface that represents a parser or command that sub-commands can be added to.
// This is used to resolve a tree of commands.
type AddCommander interface {
	AddCommand(command, shortDescription, longDescription string, data any) (
		*flags.Command,
		error,
	)
}

// AddCommand resolves a command tree to the parent command.
func AddCommand(parent AddCommander, command Command) error {
	cmd, err := parent.AddCommand(command.CommandName(), command.ShortDesc(), command.LongDesc(), command)
	if err != nil {
		return err
	}
	cmd.Aliases = command.GetAliases()
	if subCommander, ok := command.(subCommander); ok {
		for _, subCommand := range subCommander.SubCommands() {
			err := AddCommand(cmd, subCommand)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
