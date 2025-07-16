package claude_code

// package to call the claude code cli from go

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/utils"
)

var (
	TotalCostUsd float64 = 0
)

type PromptOps struct {
	AllowedTools []string
	AllowEdits   bool
	Continue     bool
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

type ClaudeJSONResponse struct {
	Result       string  `json:"result"`
	Subtype      string  `json:"subtype"`
	IsError      bool    `json:"is_error"`
	TotalCostUsd float64 `json:"total_cost_usd"`
}

func (r ClaudeJSONResponse) String() string {
	s := ""
	if r.IsError {
		s += utils.Colors.Error("ERROR\n")
	}
	s += fmt.Sprintf("%s\n", r.Result)
	s += utils.Colors.Lowkey(fmt.Sprintf("(%s / total cost USD: %v)", r.Subtype, r.TotalCostUsd))
	return s
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
	if ops.Continue {
		args = append(args, "--continue")
	}

	args = append(args, "--output-format", "json")
	// args = append(args, "--verbose")

	cmd := exec.Command("claude", args...)

	// execute the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute claude command: %w", err)
	}

	var response ClaudeJSONResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal claude response: %w", err)
	}

	TotalCostUsd += response.TotalCostUsd

	return response.String(), nil
}
