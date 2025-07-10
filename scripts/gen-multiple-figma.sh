#!/bin/bash

# This script is used to generate multiple figma designs using claude code.
# It takes a list of figma design URLs as input and generates the HTML/CSS for each design.
# It uses the claude code API to generate the HTML/CSS.

# Usage: ./gen-multiple-figma.sh <list-of-figma-urls>

# Example: ./gen-multiple-figma.sh "someURL1" "someURL2" "someURL3"

# Check if the list of figma URLs is provided

if [ -z "$@" ]; then
    echo "Error: No figma URLs provided"
    echo "Usage: ./gen-multiple-figma.sh <list-of-figma-urls>"
    exit 1
fi

# Get the list of figma URLs from the command line
figma_urls="$@"

base_prompt="Generate HTML/CSS for the following figma design, without using SPA frameworks like React. Put the generated code under a directory named "generated", and in a directory named after the figma layer name."

# Loop through each figma URL and generate the HTML/CSS
for url in $figma_urls; do
    echo "Generating HTML/CSS for $url"
    claude -p "$base_prompt. Design URL: $url" \
      --allowedTools mcp__figma-mcp-1 --permission-mode acceptEdits
done