/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	claude_code "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/claude-code"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/codegen"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/logging"
	"github.com/spf13/cobra"
)

// devTestCmd represents the devTest command
var devTestCmd = &cobra.Command{
	Use:   "devTest",
	Short: "test code used for debugging; no official purpose",
	Long:  `test code used for debugging; no official purpose`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := claude_code.Prompt("explain the process of how python code compiles to machine code", claude_code.PromptOps{})
		if err != nil {
			logging.Fatal(err)
		}

		fmt.Println("claude response:", output.Result)
		fmt.Println("will open in claude interactive in 10 seconds...")
		time.Sleep(10 * time.Second)
		codegen.OpenInClaudeInteractive(output.SessionID)
	},
}

func init() {
	rootCmd.AddCommand(devTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
