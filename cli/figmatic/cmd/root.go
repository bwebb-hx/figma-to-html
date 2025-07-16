/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/logging"
	"github.com/spf13/cobra"
)

var setWorkingDir string
var writeLogs bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "figmatic",
	Short: "A CLI for generating HTML/CSS code from Figma designs.",
	Long: `
		figmatic is a CLI for generating HTML/CSS code from Figma designs.
		You can pass a Figma design URL, and the tool will generate HTML/CSS code for the design.
		`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if setWorkingDir != "" {
		setWorkingDir, err := filepath.Abs(setWorkingDir)
		if err != nil {
			logging.Fatal(err)
		}
		if err := os.Chdir(setWorkingDir); err != nil {
			logging.Fatal(err)
		}
		log.Printf("Set working directory to %s", setWorkingDir)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&setWorkingDir, "working-dir", "w", "", "Set the working directory for the application. This affects where tools like claude code will operate.")
	rootCmd.PersistentFlags().BoolVar(&writeLogs, "logs", false, "If set, logs will be written for debugging purposes.")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.figmatic.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
