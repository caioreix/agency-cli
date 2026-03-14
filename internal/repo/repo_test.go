package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caioreix/agency-cli/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheDir_IsUnderHome(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	dir, err := repo.CacheDir()
	require.NoError(t, err)

	// Must be located inside the user home directory.
	rel, err := filepath.Rel(home, dir)
	require.NoError(t, err)
	assert.False(t, filepath.IsAbs(rel), "CacheDir should be under home: %s", dir)
}

func TestCacheDir_EndsWithRepoName(t *testing.T) {
	t.Parallel()
	dir, err := repo.CacheDir()
	require.NoError(t, err)
	assert.Equal(t, "agency-agents", filepath.Base(dir))
}

func TestCacheDir_ContainsCacheSegment(t *testing.T) {
	t.Parallel()
	dir, err := repo.CacheDir()
	require.NoError(t, err)
	// Path must contain "agency-cli" as a directory segment.
	assert.Contains(t, dir, filepath.Join("agency-cli"))
}
