package gomonorepo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/semaphore"

	"github.com/drshriveer/gtools/gsync"
)

func normalizeFilePath(filename flags.Filename) string {
	s, err := filepath.Abs(string(filename))
	if err != nil {
		panic(err)
	}
	return s
}

type commandResult struct {
	cmd       string
	output    *bytes.Buffer
	err       error
	succeeded bool
}

func (cr *commandResult) Print() {
	_, _ = fmt.Fprintln(os.Stdout, "→", cr.cmd)
	_, _ = cr.output.WriteTo(os.Stdout)
	PutBuffer(cr.output)
}
func (cr *commandResult) join(other *commandResult) {
	cr.output.WriteString("→ " + other.cmd + "\n")
	_, _ = other.output.WriteTo(cr.output)
	cr.err = errors.Join(cr.err, other.err)
	cr.succeeded = cr.succeeded && other.succeeded
}

func invokeOnElement[T any](
	ctx context.Context,
	opts *AppOptions,
	mods []T,
	f func(ctx context.Context, m T) (commandResult, error),
) (success bool, err error) {
	success = true // start this way

	if opts.Parallelism == 1 {
		var cr commandResult
		for _, m := range mods {
			cr, err = f(ctx, m)
			if err != nil {
				return false, err
			}
			success = successAgg(cr, success)
		}
		return success, nil
	}

	sem := semaphore.NewWeighted(int64(opts.Parallelism))
	executor, done := gsync.NewExecutor(ctx, successAgg, success)
	defer done()

	for _, m := range mods {
		err = executor.AddTask(func(ctx context.Context) (r commandResult, err error) {
			err = sem.Acquire(ctx, 1)
			if err != nil {
				return r, err
			}
			defer sem.Release(1)
			return f(ctx, m)
		})
		if err != nil {
			return false, err
		}
	}
	return executor.WaitAndResult()
}

func successAgg(cr commandResult, success bool) bool {
	cr.Print()
	return cr.succeeded && success
}

func runCommand(ctx context.Context, args []string) commandResult {
	cr := commandResult{}
	cr.output, _ = GetBuffer()
	//nolint:gosec // G204 shouldn't be an issue here.
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = cr.output
	cmd.Stderr = cr.output
	cr.cmd = buildCommandName(cmd)
	cr.err = cmd.Run()
	cr.succeeded = cr.err == nil
	return cr
}

func buildCommandName(cmd *exec.Cmd) string {
	sb, done := GetBuffer()
	defer done(sb)
	for i, a := range cmd.Args {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(a)
	}

	return sb.String()
}

func ensureRecursivePath(path string) string {
	if strings.HasSuffix(path, "/...") {
		return path
	}
	path = strings.TrimSuffix(path, "/") // just in case
	path += "/..."
	return path
}
