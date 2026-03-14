package repo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	repoURL  = "https://github.com/msitarzewski/agency-agents.git"
	cacheDir = ".cache/agency-cli"
	repoName = "agency-agents"
)

func CacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, cacheDir, repoName), nil
}

func EnsureRepo() (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", err
	}

	if _, statErr := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(statErr) {
		return dir, cloneRepo(dir)
	}

	return dir, nil
}

func Sync() (string, int, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", 0, err
	}

	if _, statErr := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(statErr) {
		if cloneErr := cloneRepo(dir); cloneErr != nil {
			return dir, 0, cloneErr
		}
		count, countErr := commitCount(dir)
		return dir, count, countErr
	}

	beforeHash, err := currentHash(dir)
	if err != nil {
		return dir, 0, err
	}

	cmd := exec.CommandContext(context.Background(), "git", "pull", "--ff-only")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		return dir, 0, fmt.Errorf("git pull failed: %w", runErr)
	}

	afterHash, err := currentHash(dir)
	if err != nil {
		return dir, 0, err
	}

	if beforeHash == afterHash {
		return dir, 0, nil
	}

	count, err := commitsBetween(dir, beforeHash, afterHash)
	return dir, count, err
}

func cloneRepo(dir string) error {
	parent := filepath.Dir(dir)
	if err := os.MkdirAll(parent, 0o750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	cmd := exec.CommandContext(context.Background(), "git", "clone", "--depth=1", repoURL, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func currentHash(dir string) (string, error) {
	cmd := exec.CommandContext(context.Background(), "git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func commitCount(dir string) (int, error) {
	cmd := exec.CommandContext(context.Background(), "git", "rev-list", "--count", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var count int
	if _, scanErr := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &count); scanErr != nil {
		return 0, fmt.Errorf("failed to parse commit count: %w", scanErr)
	}
	return count, nil
}

func commitsBetween(dir, from, to string) (int, error) {
	cmd := exec.CommandContext( //nolint:gosec // G204: git hashes are controlled
		context.Background(), "git", "rev-list", "--count", from+".."+to)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var count int
	if _, scanErr := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &count); scanErr != nil {
		return 0, fmt.Errorf("failed to parse commit count: %w", scanErr)
	}
	return count, nil
}
