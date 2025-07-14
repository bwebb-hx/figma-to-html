package claude_code

// package to call the claude code cli from go

import (
	"fmt"
	"os/exec"
	"strings"
)

type PromptOps struct {
	AllowedTools []string
	AllowEdits   bool
}

const FigmaMCP = "mcp__figma-mcp-1"

var DefaultPromptOps = PromptOps{
	AllowedTools: []string{FigmaMCP},
	AllowEdits:   true,
}

func ConfirmClaudeIsInstalled() error {
	_, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude code CLI is not installed")
	}
	return nil
}

func CheckMCPs() string {
	cmd := exec.Command("claude", "mcp", "list")
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func MCPExists(mcpName string) bool {
	mcpString := CheckMCPs()
	return strings.Contains(mcpString, mcpName)
}

func Prompt(prompt string, ops PromptOps) (string, error) {
	args := []string{"-p", fmt.Sprintf("\"%s\"", prompt)}

	if len(ops.AllowedTools) > 0 {
		args = append(args, "--allowedTools")
		args = append(args, ops.AllowedTools...)
	}
	if ops.AllowEdits {
		args = append(args, "--permission-mode", "acceptEdits")
	}

	cmd := exec.Command("claude", args...)

	fmt.Println("executing claude command:", cmd.Args)

	// execute the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute claude command: %w", err)
	}

	fmt.Println("claude command output:", string(output))

	return string(output), nil
}
