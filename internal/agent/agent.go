package agent

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var agentDirs = []string{
	"design", "engineering", "game-development", "marketing",
	"paid-media", "sales", "product", "project-management",
	"testing", "support", "spatial-computing", "specialized",
}

type Agent struct {
	Name        string
	Description string
	Color       string
	Emoji       string
	Vibe        string
	Tools       string
	Category    string
	Slug        string
	Body        string
	FilePath    string
}

func Slugify(name string) string {
	s := strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func Parse(filePath, category string) (*Agent, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Check first line is ---
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, fmt.Errorf("not an agent file (no frontmatter): %s", filePath)
	}

	// Parse frontmatter
	fields := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		if idx := strings.Index(line, ": "); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+2:])
			// Strip surrounding quotes (single or double)
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			fields[key] = val
		}
	}

	// Read body
	var body strings.Builder
	for scanner.Scan() {
		body.WriteString(scanner.Text())
		body.WriteString("\n")
	}

	name := fields["name"]
	if name == "" {
		return nil, fmt.Errorf("agent has no name: %s", filePath)
	}

	return &Agent{
		Name:        name,
		Description: fields["description"],
		Color:       fields["color"],
		Emoji:       fields["emoji"],
		Vibe:        fields["vibe"],
		Tools:       fields["tools"],
		Category:    category,
		Slug:        Slugify(name),
		Body:        body.String(),
		FilePath:    filePath,
	}, nil
}

func ListAll(repoDir string) ([]*Agent, error) {
	var agents []*Agent

	for _, dir := range agentDirs {
		dirPath := filepath.Join(repoDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			agent, err := Parse(filepath.Join(dirPath, entry.Name()), dir)
			if err != nil {
				continue // skip non-agent files
			}

			agents = append(agents, agent)
		}
	}

	return agents, nil
}

func FindBySlug(repoDir, slug string) (*Agent, error) {
	agents, err := ListAll(repoDir)
	if err != nil {
		return nil, err
	}

	for _, a := range agents {
		if a.Slug == slug {
			return a, nil
		}
	}

	return nil, fmt.Errorf("agent not found: %s", slug)
}
