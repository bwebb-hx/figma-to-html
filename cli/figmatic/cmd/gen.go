/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
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

			// generate each node concurrently
			generateHTMLWorkGroup(nodes, iterations)

			if len(nodes) == 1 {
				log.Println("only one URL generated; skipping combine step.")
				os.Exit(0)
			}

			fmt.Println("\nCombining HTML files. This could take a few minutes...")
			start := time.Now()
			output, err := codegen.CombineHTML(figmaURL)
			if err != nil {
				fmt.Fprint(os.Stderr, output.String()+"\n")
				log.Fatal(err)
			}
			utils.Colors.Lowkey(fmt.Sprintf("done! %v seconds elapsed\n", time.Since(start).Seconds()))
			fmt.Println("\nClaude:")
			fmt.Println(output)
		} else {
			// generate HTML for the original node
			stats, err := generateHTML(figmaURL, "root", iterations, "")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Finished generating HTML! Turns: %v | Cost: $%.2f\n", stats.TotalTurns, stats.TotalCostUSD)
		}
	},
}

func generateHTMLWorkGroup(nodes []figma_api.FigmaNode, numIterations int) {
	var wg sync.WaitGroup

	fmt.Printf("Generating %v Figma nodes concurrently...\n", len(nodes))
	if numIterations > 0 {
		fmt.Printf("Each node's HTML will be iterated on %v extra time(s)\n", numIterations)
	}

	start := time.Now()
	for _, node := range nodes {
		wg.Add(1)
		go generateHTMLWorker(node.URL, node.LayerName, numIterations, &wg)
	}

	wg.Wait()
	fmt.Printf("All workers done generating HTML! Time elapsed: %s\n", time.Since(start).Truncate(time.Second))
}

func generateHTMLWorker(url string, nodeName string, numIterations int, wg *sync.WaitGroup) {
	defer wg.Done()
	utils.Colors.LowkeyPrint(fmt.Sprintf("[%s] Worker generating HTML for %s...", nodeName, url))

	start := time.Now()
	stats, err := generateHTML(url, nodeName, numIterations, nodeName)
	if err != nil {
		log.Println(err)
	}
	utils.Colors.LowkeyPrint(fmt.Sprintf("[%s] Worker finished! Time: %s | Turns: %v | Cost USD: $%.2f", nodeName, time.Since(start).Truncate(time.Second), stats.TotalTurns, stats.TotalCostUSD))
}

type generateHTMLStats struct {
	TotalTurns   int
	TotalCostUSD float64
}

// generateHTML generates HTML for the given URL. If workerName is set, output will be minimized and shown as though part of a work group that is executing concurrently.
func generateHTML(url string, nodeName string, numIterations int, workerName string) (generateHTMLStats, error) {
	stats := generateHTMLStats{}

	output, err := codegen.GenerateHTML(url, nodeName)
	if err != nil {
		return stats, err
	}
	stats.TotalTurns += output.NumTurns
	stats.TotalCostUSD += output.TotalCostUsd

	if workerName == "" {
		fmt.Println("\nClaude:")
		fmt.Println(output)
	}

	if numIterations > 0 {
		for i := range numIterations {
			if workerName == "" {
				fmt.Printf("[%v/%v] Iterating on the HTML for %s...\n", i+1, numIterations, figmaURL)
			} else {
				utils.Colors.LowkeyPrint(fmt.Sprintf("[%s] Iteration %v/%v...", workerName, i+1, numIterations))
			}
			output, err = codegen.ContinueImproveHTML(output.SessionID)
			if err != nil {
				return stats, err
			}
			stats.TotalTurns += output.NumTurns
			stats.TotalCostUSD += output.TotalCostUsd

			if workerName == "" {
				utils.Colors.LowkeyPrint("Iteration output from Claude:")
				utils.Colors.LowkeyPrint(output.String())
			}
		}
		if workerName == "" {
			fmt.Println("Iterations complete!")
		}
	}

	return stats, nil
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&figmaURL, "url", "u", "", "Figma design URL")
	genCmd.Flags().StringVarP(&figmaAccessToken, "figma-access-token", "t", "", "Figma access token")
	genCmd.Flags().BoolVarP(&genSubNodes, "sub-nodes", "s", false, "Generate HTML for sub-nodes and then combine them. Increases accuracy, especially for complex designs, but increases time and cost.")
	genCmd.Flags().IntVarP(&iterations, "iterations", "i", 0, "Number of extra iterations for Claude to refine the generated code. Defaults to 0. An iteration involves asking Claude to review how accurate the code is and try to improve it.")
}
