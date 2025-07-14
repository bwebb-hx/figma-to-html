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
			urls, err := figma_api.GetNodeURLs(figmaURL, figmaAccessToken)
			if err != nil {
				log.Fatal(err)
			}
			// if no sub-node URLs are found, try using the original node URL
			if len(urls) == 0 {
				fmt.Println("No sub-node URLs found. Using the original node URL.")
				urls = append(urls, figmaURL)
			} else {
				fmt.Printf("Found %v sub-node URLs.\n", len(urls))
			}

			var totalElapsed time.Duration

			for i, url := range urls {
				fmt.Printf("\n[%v/%v] Generating HTML for %s\n", i+1, len(urls), url)
				if totalElapsed > 0 {
					timeEstimate := totalElapsed / time.Duration(i) * time.Duration(len(urls)-i)
					utils.Colors.LowkeyPrint(fmt.Sprintf("time remaining: %s", timeEstimate.Truncate(time.Second)))
				}

				start := time.Now()
				output, err := codegen.GenerateHTML(url)
				if err != nil {
					fmt.Fprint(os.Stderr, output+"\n")
					log.Fatal(err)
				}
				utils.Colors.Lowkey(fmt.Sprintf("done! %v seconds elapsed\n", time.Since(start).Seconds()))
				totalElapsed += time.Since(start)
			}

			if len(urls) == 1 {
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
			output, err := codegen.GenerateHTML(figmaURL)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("\nClaude:")
			fmt.Println(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&figmaURL, "url", "u", "", "Figma design URL")
	genCmd.Flags().StringVarP(&figmaAccessToken, "figma-access-token", "t", "", "Figma access token")
	genCmd.Flags().BoolVarP(&genSubNodes, "sub-nodes", "s", false, "Generate HTML for sub-nodes and then combine them. Increases accuracy, especially for complex designs, but increases time and uses more claude code API calls.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
