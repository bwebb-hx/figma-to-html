/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwebb-hx/figma-to-html/cli/figmatic/cmd"
	claude_code "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/claude-code"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/utils"
)

func main() {
	mainStart := time.Now()
	defer func() {
		// print session summary
		utils.Colors.LowkeyPrint("\n--------------------------------")
		utils.Colors.LowkeyPrint(fmt.Sprintf("time elapsed: %s", time.Since(mainStart).Truncate(time.Second)))
		utils.Colors.LowkeyPrint(fmt.Sprintf("total cost USD: $%.2f", claude_code.TotalCostUsd))
	}()

	if err := claude_code.ConfirmClaudeIsInstalled(); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
