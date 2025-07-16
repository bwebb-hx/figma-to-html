package claude_code

// package to call the claude code cli from go

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	globalconfig "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/global-config"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/logging"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/utils"
)

var (
	TotalCostUsd float64 = 0
)

type PromptOps struct {
	AllowedTools []string
	AllowEdits   bool
	Continue     bool
	Resume       string
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
	SessionID    string  `json:"session_id"`
	NumTurns     int     `json:"num_turns"`
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

func getArgsFromOptions(ops PromptOps) []string {
	var args []string

	if len(ops.AllowedTools) > 0 {
		args = append(args, "--allowedTools")
		args = append(args, ops.AllowedTools...)
	}
	if ops.AllowEdits {
		args = append(args, "--permission-mode", "acceptEdits")
	}
	if ops.Continue {
		args = append(args, "--continue")
	} else if ops.Resume != "" {
		args = append(args, "--resume", ops.Resume)
	}

	return args
}

func Prompt(prompt string, ops PromptOps) (ClaudeJSONResponse, error) {
	args := []string{"-p", fmt.Sprintf("\"%s\"", prompt)}

	opArgs := getArgsFromOptions(ops)
	if len(opArgs) > 0 {
		args = append(args, opArgs...)
	}

	args = append(args, "--output-format", "json")
	// args = append(args, "--verbose")

	cmd := exec.Command("claude", args...)

	// execute the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ClaudeJSONResponse{}, fmt.Errorf("failed to execute claude command: %w", err)
	}

	var response ClaudeJSONResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return ClaudeJSONResponse{}, fmt.Errorf("failed to unmarshal claude response: %w", err)
	}

	if globalconfig.WRITE_LOG_FILES {
		logging.WriteLog("prompt", prompt+"\n\n"+response.String())
	}

	TotalCostUsd += response.TotalCostUsd

	return response, nil
}

func RunClaudeInTerminal(ops PromptOps) error {
	args := getArgsFromOptions(ops)

	binary, err := exec.LookPath("claude")
	if err != nil {
		return err
	}

	fullArgs := []string{"claude"}
	fullArgs = append(fullArgs, args...)

	env := os.Environ()

	err = syscall.Exec(binary, fullArgs, env)
	// if the above worked, go process should now have exited
	return err
}
