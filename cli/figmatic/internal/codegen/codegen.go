package codegen

import (
	"fmt"

	claude_code "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/claude-code"
)

func GenerateHTML(url string, nodeName string) (string, error) {
	basePrompt := fmt.Sprintf("Generate HTML/CSS for the following figma design, without using SPA frameworks like React. Put the generated code in a directory called 'generated/%s'.", nodeName)
	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, url)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}
	return output, nil
}

func ImproveHTML(url string, nodeName string) (string, error) {
	basePrompt := fmt.Sprintf("There is HTML/CSS code in the directory called 'generated/%s', which was generated from the following figma design. Compare the code to the figma design and fix any missing elements, inaccuracies, etc.", nodeName)
	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, url)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return "", fmt.Errorf("failed to improve HTML: %w", err)
	}
	return output, nil
}

func CombineHTML(topLevelURL string) (string, error) {
	basePrompt := `In the directory called 'generated', there are several directories that each have the HTML/CSS for a layer in a figma design.
Put them all together into a single HTML file (and single CSS file).
Take a look at this figma design, and try to put all of these HTML files together.`

	prompt := fmt.Sprintf("%s Design URL: %s", basePrompt, topLevelURL)
	output, err := claude_code.Prompt(prompt, claude_code.DefaultPromptOps)
	if err != nil {
		return "", fmt.Errorf("failed to combine HTML: %w", err)
	}
	return output, nil
}
