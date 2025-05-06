package gomonorepo

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// getIgnoredDirectories calls `git status --porcelain --ignored` to get a list directories
// that we should not recurse down.
func getIgnoredDirectories(
	ctx context.Context,
	rootDir string,
) ([]string, error) {
	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain", "--ignored")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get ignored files: %w\n%s", err, stderr.String())
	}

	result := make([]string, 0, 8)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && strings.HasPrefix(line, "!! ") {
			temp := strings.TrimPrefix(line, "!! ")
			temp = filepath.Join(rootDir, temp)
			result = append(result, temp)
		}
	}

	return result, nil
}

func listChangedFiles(ctx context.Context, parent string, patched bool) ([]string, error) {
	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)
	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", parent)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if !patched && strings.HasPrefix(stderr.String(), "fatal: ambiguous argument '"+parent+"'") {
			err = tryFetchParentRevision(ctx, parent)
			if err != nil {
				return nil, err
			}
			return listChangedFiles(ctx, parent, true)
		}
		return nil, fmt.Errorf("failed to run git diff: %w\n%s", err, stderr.String())
	}

	result := make([]string, 0, 8)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

func tryFetchParentRevision(ctx context.Context, parent string) error {
	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)

	// git rev-list --count origin/main..$(git branch --show-current)
	cmd := exec.CommandContext(ctx, "git", "fetch", parent)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err == nil {
		return nil
	}
	remote, branch, ok := strings.Cut(parent, "/")
	if !ok {
		return fmt.Errorf("failed to fetch parent revision: %w\n%s", err, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	cmd = exec.CommandContext(ctx, "git", "fetch", remote, branch)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to fetch parent revision: %w\n%s", err, stderr.String())
	}
	return nil
}

// getCurrentBranch returns the current branch name.
func getCurrentBranch(ctx context.Context) (remote string, branch string, err error) {
	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)

	// Get the current branch / revision name excluding the remote:
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to get branch name: %w\n%s", err, stderr.String())
	}
	branch = strings.TrimSpace(stdout.String())

	stdout.Reset()
	stderr.Reset()

	// Get the branch name AND its upstream remote (e.g. origin/main):
	cmd = exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		// noop. there is no remote.
		return "", branch, nil
	}
	remote = strings.TrimSuffix(strings.TrimSpace(stdout.String()), "/"+branch)
	return remote, branch, nil
}
