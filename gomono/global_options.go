package gomono

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	stdLog = log.New(os.Stdout, "[gomono] ", log.LstdFlags)
	errLog = log.New(os.Stderr, "[gomono] ", log.LstdFlags)
)

type GlobalOptions struct {
	Root        flags.Filename `long:"root" short:"r" description:"Root directory of the mono repo." default:"."`
	Verbose     bool           `long:"verbose" short:"v" description:"Enable verbose logging."`
	Parallelism int            `long:"parallelism" short:"p" description:"Permitted parallelism for tasks that can be parallelized." default:"4"`
	Timeout     time.Duration  `long:"timeout" short:"t" description:"Timeout for the command." default:"5m"`
}

// GetRoot returns the root directory of the mono repo.
func (x *GlobalOptions) GetRoot() string {
	return normalizeFilePath(x.Root)
}

// Infof prints a message to stdout, using the logger, which will
// be formatted with a [gomono] prefix and a timestamp.
func (x *GlobalOptions) Infof(format string, args ...any) {
	stdLog.Printf(format, args...)
}

// Printf prints directly to stdout, without using the logger which
// adds a prefix and timestamp.
// This should be used for continuing input, or when we want to pipe
// the exact output to stdout.
func (x *GlobalOptions) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

// Errorf prints a message to stderr, using the logger, which will
// be formatted with a [gomono] prefix and a timestamp.
func (x *GlobalOptions) Errorf(format string, args ...any) {
	errLog.Fatalf(format, args...)
}

// Pipeln pipes a string directly to stdout.
func (x *GlobalOptions) Pipeln(s string) {
	_, _ = fmt.Fprintln(os.Stdout, s)
}
