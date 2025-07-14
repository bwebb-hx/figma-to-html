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
	"github.com/spf13/cobra"
)

var figmaURL string
var figmaAccessToken string
var genSubNodes bool

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			urls, err := figma_api.GetNodeURLs(figmaURL, figmaAccessToken)
			if err != nil {
				log.Fatal(err)
			}
			// if no sub-node URLs are found, try using the original node URL
			if len(urls) == 0 {
				log.Println("no sub-node URLs found")
				urls = append(urls, figmaURL)
			}

			i := 0
			for _, url := range urls {
				i++
				fmt.Printf("[%v/%v] Generating HTML for %s\n", i, len(urls), url)
				start := time.Now()
				output, err := codegen.GenerateHTML(url)
				if err != nil {
					fmt.Fprint(os.Stderr, output+"\n")
					log.Fatal(err)
				}
				fmt.Printf("done! %v seconds elapsed\n", time.Since(start).Seconds())
			}

			if len(urls) == 1 {
				log.Println("only one URL generated; skipping combine step.")
				os.Exit(0)
			}

			output, err := codegen.CombineHTML(figmaURL)
			if err != nil {
				fmt.Fprint(os.Stderr, output+"\n")
				log.Fatal(err)
			}
			fmt.Println(output)
		} else {
			// generate HTML for the original node
			output, err := codegen.GenerateHTML(figmaURL)
			if err != nil {
				log.Fatal(err)
			}
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
