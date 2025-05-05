package gomonorepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	stdLog = log.New(os.Stdout, "[gomono] ", log.LstdFlags)
	errLog = log.New(os.Stderr, "[gomono] ", log.LstdFlags)
)

// AppOptions hold global configurations an options that are used by most if not all commands.
type AppOptions struct {
	Root                      flags.Filename `long:"root" short:"r" description:"Root directory of the mono repo." default:"."`
	InvocationDir             flags.Filename `long:"invocationDir" description:"If invocationDir is not the root directory, invocation will be limited to the invocationDir and its subdirectories."`
	InvocationDirNotRecursive bool           `long:"non-recursive-from-invocation-dir" description:"When using invocationDir, turn off recursive invocation; limit to the invocation directory only."`
	Verbose                   bool           `long:"verbose" short:"v" description:"Enable verbose logging."`
	Parallelism               int            `long:"parallelism" short:"p" description:"Permitted parallelism for tasks that can be parallelized." default:"4"`
	Timeout                   time.Duration  `long:"timeout" short:"t" description:"Timeout for the command." default:"5m"`
	ExcludePaths              []string       `long:"excludePath" short:"x" description:"Paths to to exclude from searches (these may be regex). Note: Anything excluded by git is ignored by default." default:"node_modules" default:"vendor"`
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

// GetFocusDir returns the focus directory of the mono repo.
// The bool returned indicates if the focus directory is valid and should be respected.
func (x *AppOptions) GetFocusDir() (string, bool) {
	if x.InvocationDir == "" {
		return "", false
	}

	root := x.GetRoot()
	focus := normalizeFilePath(x.InvocationDir)
	willUseFocus := focus != root && strings.HasPrefix(focus, root)
	if x.Verbose && willUseFocus {
		x.Infof("Focusing on directory: %q", focus)
	}
	if !x.InvocationDirNotRecursive && !strings.HasSuffix(focus, "/...") {
		focus = strings.TrimSuffix(focus, "/") // just in case
		focus += "/..."
	}
	return focus, willUseFocus
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
