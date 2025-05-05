package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"

	"github.com/drshriveer/gtools/gomonorepo"
)

func main() {
	opts := &gomonorepo.AppOptions{}

	// set up a new parser with the updated cfg struct.
	parser := flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash)
	cmds := []gomonorepo.Command{
		gomonorepo.DetectedChangesCommand,
		gomonorepo.TestModulesCommand,
		gomonorepo.LintModulesCommand,
		gomonorepo.FormatModulesCommand,
		gomonorepo.GenerateModulesCommand,
		gomonorepo.TidyModulesCommand,
		gomonorepo.ListModulesCommand,
	}
	var err error
	for _, cmd := range cmds {
		err = gomonorepo.AddCommand(parser, cmd)
		if err != nil {
			panic(err)
		}
	}
	// CommandHandler allows us to pre-process anything ahead of executing a command,
	// and also lets us invoke our own command function for the command we are running,
	// which is nicer for context propagation and the like.
	parser.CommandHandler = func(cmd flags.Commander, _ []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
		ctx, cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		defer cancel()

		return cmd.(gomonorepo.Commander).RunCommand(ctx, opts)
	}

	_, err = parser.Parse()
	if fErr, ok := err.(*flags.Error); ok && fErr.Type == flags.ErrHelp {
		parser.WriteHelp(os.Stdout)
	} else if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}
