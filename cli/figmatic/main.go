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
)

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Printf("(time elapsed: %s)\n", time.Since(mainStart))
	}()

	if err := claude_code.ConfirmClaudeIsInstalled(); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
