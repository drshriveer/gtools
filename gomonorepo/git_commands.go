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
		return nil, fmt.Errorf("failed to run git diff: %w\n%s", err, stderr.String())
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

func listChangedFiles(ctx context.Context, parent string) ([]string, error) {
	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)
	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", parent)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
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

// getCurrentBranch returns the current branch name.
func getCurrentBranch(ctx context.Context) (remote string, branch string, err error) {
	// git rev-parse --verify HEAD <-- prints out the current hash
	// What about unstaged?
	// git rev-parse --abbrev-ref --symbolic-full-name @{upstream} <-- prints out the upstream branch OR fatals with "fatal: no upstream configured ..."
	//

	stdout, done := GetBuffer()
	defer done(stdout)
	stderr, done := GetBuffer()
	defer done(stderr)
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

	cmd = exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to get upstream: %w\n%s", err, stderr.String())
	}

	remote = strings.TrimSpace(stdout.String())
	if strings.HasPrefix(remote, "fatal: no") {
		remote = ""
	}

	return remote, branch, nil
}
