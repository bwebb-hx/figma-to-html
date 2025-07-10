#!/bin/bash

# given a figma design URL, this script will build the HTML/CSS for it.
# here's an overview of the process:
# 1. get the URLs for each layer below the given figma design node
# 2. for each layer, generate the HTML/CSS files
# 3. finally, put all the HTML/CSS files together into a single HTML file

# Usage: ./gen-figma-by-layer.sh <figma-design-url> <figma-access-token>

if [ -z "$1" ]; then
    echo "Error: No figma design URL provided"
    echo "Usage: ./gen-figma-by-layer.sh <figma-design-url> <figma-access-token>"
    exit 1
fi

if [ -z "$2" ]; then
    echo "Error: No figma access token provided"
    echo "Usage: ./gen-figma-by-layer.sh <figma-design-url> <figma-access-token>"
    exit 1
fi

# get the URLs for each layer below the given figma design node
urls=$(bash scripts/get-layer-urls.sh "$1" "$2")

if [ -z "$urls" ]; then
    echo "Error: No URLs found for the given figma design"
    exit 1
fi

# put the URLs into an array
urls_array=()
while read -r url; do
    urls_array+=("$url")
done <<< "$urls"

# pass all the URLs to the gen-multiple-figma.sh script
bash scripts/gen-multiple-figma.sh "${urls_array[@]}"