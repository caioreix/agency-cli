package repo

import (
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

	if _, err := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(err) {
		return dir, cloneRepo(dir)
	}

	return dir, nil
}

func Sync() (string, int, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", 0, err
	}

	if _, err := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(err) {
		if err := cloneRepo(dir); err != nil {
			return dir, 0, err
		}
		count, err := commitCount(dir)
		return dir, count, err
	}

	beforeHash, err := currentHash(dir)
	if err != nil {
		return dir, 0, err
	}

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return dir, 0, fmt.Errorf("git pull failed: %w", err)
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
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth=1", repoURL, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func currentHash(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func commitCount(dir string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var count int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &count)
	return count, nil
}

func commitsBetween(dir, from, to string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", from+".."+to)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var count int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &count)
	return count, nil
}
