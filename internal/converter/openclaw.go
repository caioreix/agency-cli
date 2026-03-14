package converter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caioreix/agency-cli/internal/agent"
)

type openclaw struct{}

func init() {
	Register("openclaw", &openclaw{})
}

func (c *openclaw) Name() string          { return "OpenClaw" }
func (c *openclaw) Description() string   { return "~/.openclaw/agency-agents/" }
func (c *openclaw) IsProjectScoped() bool { return false }

func (c *openclaw) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	// openclaw installs globally regardless of scope
	agentDir := filepath.Join(destDir, a.Slug)
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		return nil, err
	}

	soulContent, agentsContent := splitOpenClawSections(a.Body)

	// SOUL.md — persona, tone, boundaries
	soulFile := filepath.Join(agentDir, "SOUL.md")
	if err := os.WriteFile(soulFile, []byte(soulContent), 0o644); err != nil {
		return nil, err
	}

	// AGENTS.md — mission, deliverables, workflow
	agentsFile := filepath.Join(agentDir, "AGENTS.md")
	if err := os.WriteFile(agentsFile, []byte(agentsContent), 0o644); err != nil {
		return nil, err
	}

	// IDENTITY.md
	identityFile := filepath.Join(agentDir, "IDENTITY.md")
	var identityContent string
	if a.Emoji != "" && a.Vibe != "" {
		identityContent = "# " + a.Emoji + " " + a.Name + "\n" + a.Vibe + "\n"
	} else {
		identityContent = "# " + a.Name + "\n" + a.Description + "\n"
	}
	if err := os.WriteFile(identityFile, []byte(identityContent), 0o644); err != nil {
		return nil, err
	}

	return []string{soulFile, agentsFile, identityFile}, nil
}

func splitOpenClawSections(body string) (soul, agents string) {
	soulKeywords := []string{"identity", "communication", "style", "critical rule", "rules you must follow"}
	lines := strings.Split(body, "\n")

	var soulBuilder, agentsBuilder strings.Builder
	currentTarget := "agents"
	var currentSection strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			// Flush current section
			if currentSection.Len() > 0 {
				if currentTarget == "soul" {
					soulBuilder.WriteString(currentSection.String())
				} else {
					agentsBuilder.WriteString(currentSection.String())
				}
				currentSection.Reset()
			}

			// Classify header
			headerLower := strings.ToLower(line)
			currentTarget = "agents"
			for _, kw := range soulKeywords {
				if strings.Contains(headerLower, kw) {
					currentTarget = "soul"
					break
				}
			}
		}

		currentSection.WriteString(line)
		currentSection.WriteString("\n")
	}

	// Flush final section
	if currentSection.Len() > 0 {
		if currentTarget == "soul" {
			soulBuilder.WriteString(currentSection.String())
		} else {
			agentsBuilder.WriteString(currentSection.String())
		}
	}

	return soulBuilder.String(), agentsBuilder.String()
}
