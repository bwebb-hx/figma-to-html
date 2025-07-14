/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/codegen"
	figma_api "github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/figma-api"
	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/utils"
	"github.com/spf13/cobra"
)

var figmaURL string
var figmaAccessToken string
var genSubNodes bool
var iterations int

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate HTML/CSS code from a Figma design",
	Long: `
		Generate HTML/CSS code from a Figma design.
		You can pass a Figma design URL, and the tool will generate HTML/CSS code for the design.
		
		For larger designs (or for increased accuracy), you can pass the --sub-nodes flag to generate the sub-nodes individually and then combine them.
		Note that this will take longer, and will use more Claude Code API calls
		
		# Example: generating a figma design, using the "sub-nodes" technique.
		figmatic gen --url "https://www.figma.com/design/FigmaDesignURL" -t "$FIGMA_ACCESS_TOKEN" --sub-nodes
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if figmaURL == "" {
			log.Fatal("figma URL is required")
		}
		if figmaAccessToken == "" {
			// try to get from env
			figmaAccessToken = os.Getenv("FIGMA_ACCESS_TOKEN")
			if figmaAccessToken == "" {
				log.Fatal("Figma access token is required. Either set it in the FIGMA_ACCESS_TOKEN environment variable, or pass it with the --figma-access-token flag.")
			}
		}

		if genSubNodes {
			fmt.Println("Getting URLs for sub-nodes of the provided Figma design node...")
			nodes, err := figma_api.GetNodeURLs(figmaURL, figmaAccessToken)
			if err != nil {
				log.Fatal(err)
			}
			// if no sub-node URLs are found, try using the original node URL
			if len(nodes) == 0 {
				fmt.Println("No sub-node URLs found. Using the original node URL.")
				nodes = append(nodes, figma_api.FigmaNode{URL: figmaURL})
			} else {
				fmt.Printf("Found %v sub-node URLs.\n", len(nodes))
			}

			var totalElapsed time.Duration

			for i, node := range nodes {
				fmt.Printf("\n[%v/%v] Generating HTML for %s\n", i+1, len(nodes), node.URL)
				if totalElapsed > 0 {
					timeEstimate := totalElapsed / time.Duration(i) * time.Duration(len(nodes)-i)
					utils.Colors.LowkeyPrint(fmt.Sprintf("time remaining: %s", timeEstimate.Truncate(time.Second)))
				}

				start := time.Now()
				output, err := codegen.GenerateHTML(node.URL, node.LayerName)
				if err != nil {
					fmt.Fprint(os.Stderr, output+"\n")
					log.Fatal(err)
				}
				if iterations > 0 {
					for i := range iterations {
						fmt.Printf("..[%v/%v] Iterating on the HTML for %s...\n", i+1, iterations, figmaURL)
						output, err = codegen.ImproveHTML(node.URL, node.LayerName)
						if err != nil {
							log.Println(err)
							break // for now, break out since we've already generated code
						}
						utils.Colors.LowkeyPrint("..Iteration output from Claude:")
						utils.Colors.LowkeyPrint(output)
					}
				}
				utils.Colors.LowkeyPrint(fmt.Sprintf("done! %v seconds elapsed\n", time.Since(start).Seconds()))
				totalElapsed += time.Since(start)
			}

			if len(nodes) == 1 {
				log.Println("only one URL generated; skipping combine step.")
				os.Exit(0)
			}

			fmt.Println("\nCombining HTML files. This could take a few minutes...")
			start := time.Now()
			output, err := codegen.CombineHTML(figmaURL)
			if err != nil {
				fmt.Fprint(os.Stderr, output+"\n")
				log.Fatal(err)
			}
			utils.Colors.Lowkey(fmt.Sprintf("done! %v seconds elapsed\n", time.Since(start).Seconds()))
			fmt.Println("\nClaude:")
			fmt.Println(output)
		} else {
			// generate HTML for the original node
			fmt.Printf("Generating HTML for %s...\n", figmaURL)

			output, err := codegen.GenerateHTML(figmaURL, "root")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("\nClaude:")
			fmt.Println(output)

			if iterations > 0 {
				for i := range iterations {
					fmt.Printf("[%v/%v] Iterating on the HTML for %s...\n", i+1, iterations, figmaURL)
					output, err = codegen.ImproveHTML(figmaURL, "root")
					if err != nil {
						log.Println(err)
						break // for now, break out since we've already generated code
					}
					utils.Colors.LowkeyPrint("Iteration output from Claude:")
					utils.Colors.LowkeyPrint(output)
				}
				fmt.Println("Iterations complete!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&figmaURL, "url", "u", "", "Figma design URL")
	genCmd.Flags().StringVarP(&figmaAccessToken, "figma-access-token", "t", "", "Figma access token")
	genCmd.Flags().BoolVarP(&genSubNodes, "sub-nodes", "s", false, "Generate HTML for sub-nodes and then combine them. Increases accuracy, especially for complex designs, but increases time and cost.")
	genCmd.Flags().IntVarP(&iterations, "iterations", "i", 0, "Number of extra iterations for Claude to refine the generated code. Defaults to 0. An iteration involves asking Claude to review how accurate the code is and try to improve it.")
}
