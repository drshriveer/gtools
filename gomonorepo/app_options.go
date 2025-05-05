package gomonorepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	stdLog = log.New(os.Stdout, "[gomono] ", log.LstdFlags)
	errLog = log.New(os.Stderr, "[gomono] ", log.LstdFlags)
)

// AppOptions hold global configurations an options that are used by most if not all commands.
type AppOptions struct {
	Root         flags.Filename `long:"root" short:"r" description:"Root directory of the mono repo." default:"."`
	Verbose      bool           `long:"verbose" short:"v" description:"Enable verbose logging."`
	Parallelism  int            `long:"parallelism" short:"p" description:"Permitted parallelism for tasks that can be parallelized." default:"4"`
	Timeout      time.Duration  `long:"timeout" short:"t" description:"Timeout for the command." default:"5m"`
	ExcludePaths []string       `long:"excludePath" short:"x" description:"Paths to to exclude from searches (these may be regex). Note: Anything excluded by git is ignored by default." default:"node_modules" default:"vendor"`
}

// ExcludePathPatterns compiles and returns ExcludePaths into regexes.
func (x *AppOptions) ExcludePathPatterns(ctx context.Context) (res Patterns, excludeDirs []string, errs error) {
	res = make([]*regexp.Regexp, len(x.ExcludePaths))
	var err error
	for i, expStr := range x.ExcludePaths {
		res[i], err = regexp.Compile(expStr)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("could not compile exclude path pattern %q: %w", expStr, err))
		}
	}
	dirs, err := getIgnoredDirectories(ctx, x.GetRoot())
	if err != nil {
		errs = errors.Join(errs, err)
	}
	return res, dirs, errs
}

// GetRoot returns the root directory of the mono repo.
func (x *AppOptions) GetRoot() string {
	return normalizeFilePath(x.Root)
}

// Infof prints a message to stdout, using the logger, which will
// be formatted with a [gomono] prefix and a timestamp.
func (x *AppOptions) Infof(format string, args ...any) {
	stdLog.Printf(format, args...)
}

// Printf prints directly to stdout, without using the logger which
// adds a prefix and timestamp.
// This should be used for continuing input, or when we want to pipe
// the exact output to stdout.
func (x *AppOptions) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

// Errorf prints a message to stderr, using the logger, which will
// be formatted with a [gomono] prefix and a timestamp.
func (x *AppOptions) Errorf(format string, args ...any) {
	errLog.Fatalf(format, args...)
}

// Pipeln pipes a string directly to stdout.
func (x *AppOptions) Pipeln(s string) {
	_, _ = fmt.Fprintln(os.Stdout, s)
}
