package codegen

import (
	"fmt"
	"log"

	claude_code "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/claude-code"
)

func GenerateHTML(url string, nodeName string) (claude_code.ClaudeJSONResponse, error) {
	basePrompt := fmt.Sprintf("Generate HTML/CSS for the following figma design, without using SPA frameworks like React. Put the generated code in a directory called 'generated/%s'.", nodeName)
	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, url)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return claude_code.ClaudeJSONResponse{}, fmt.Errorf("failed to generate HTML: %w", err)
	}
	return output, nil
}

// ImproveHTML attempts to improve the HTML generated from a Figma design. It tells Claude to look for the specific file, and then consider the given Figma design.
// Works without continuing a previous conversation, but also seems to be a little less effective than ContinueImproveHTML.
func ImproveHTML(url string, nodeName string) (claude_code.ClaudeJSONResponse, error) {
	basePrompt := fmt.Sprintf("There is HTML/CSS code in the directory called 'generated/%s', which was generated from the following figma design. Compare the code to the figma design and fix any missing elements, inaccuracies, etc.", nodeName)
	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, url)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return claude_code.ClaudeJSONResponse{}, fmt.Errorf("failed to improve HTML: %w", err)
	}
	return output, nil
}

// ContinueImproveHTML attempts to improve the HTML for the Figma design used in the last conversation with Claude (if no session ID provided), or the conversation from the given session ID.
// It seems to work better than ImproveHTML, since Claude seems to have more context of what work was being done previously.
func ContinueImproveHTML(sessionID string) (claude_code.ClaudeJSONResponse, error) {
	prompt := "Compare the HTML/CSS code to the figma design again, and add any missing elements, fix mistakes or inaccuracies, etc."
	ops := claude_code.DefaultPromptOps

	if sessionID == "" {
		ops.Continue = true
	} else {
		ops.Resume = sessionID
	}

	output, err := claude_code.Prompt(prompt, ops)
	if err != nil {
		return claude_code.ClaudeJSONResponse{}, fmt.Errorf("failed to improve HTML: %w", err)
	}
	return output, nil
}

func CombineHTML(topLevelURL string) (claude_code.ClaudeJSONResponse, error) {
	basePrompt := `In the directory called 'generated', there are several directories that each have the HTML/CSS for a layer in a figma design.
Put them all together into a single HTML file (and single CSS file).
Take a look at this figma design, and try to put all of these HTML files together.`

	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, topLevelURL)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return claude_code.ClaudeJSONResponse{}, fmt.Errorf("failed to combine HTML: %w", err)
	}
	return output, nil
}

func OpenInClaudeInteractive(sessionID string) {
	ops := claude_code.DefaultPromptOps
	if sessionID != "" {
		ops.Resume = sessionID
	} else {
		ops.Continue = true
	}

	err := claude_code.RunClaudeInTerminal(ops)
	if err != nil {
		log.Fatal(err)
	}
}
