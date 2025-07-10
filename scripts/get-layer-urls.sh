#!/bin/bash

# this script is used to get the URLs for each layer below a given node in a figma design
# the URLs are printed to stdout, one per line

NODE_URL=$1
FIGMA_ACCESS_TOKEN=$2

# from a URL, parse out the different parts
# this is for the following style of URL:
# https://www.figma.com/design/<FILE_KEY>/<FILE_NAME>?node-id=<NODE_ID_URL>&t=cXDUkWIRVTwE31dg-4

# extract the section of the URL after 'design/', up to but not including the ?
PATH_PART=$(echo "$NODE_URL" | sed -E 's|https://www.figma.com/design/([^?]+).*|\1|')

# remove everything after the first / in PATH_PART
FILE_KEY=${PATH_PART%%/*}
# remove everything before the first / in PATH_PART
FILE_NAME=${PATH_PART#*/}
# extract the node-id from the URL
NODE_ID_URL=$(echo "$NODE_URL" | grep -o 'node-id=[^&]*' | cut -d '=' -f 2)

# replace the - with : in the NODE_ID_URL for use in jq
NODE_ID_JQ="${NODE_ID_URL//-/:}"

# echo "FILE_KEY: $FILE_KEY"
# echo "FILE_NAME: $FILE_NAME"
# echo "NODE_ID_URL: $NODE_ID_URL"
# echo "NODE_ID_JQ: $NODE_ID_JQ"

# get the layers for the node from the given URL
response=$(curl -s \
  -H "X-Figma-Token: $FIGMA_ACCESS_TOKEN" \
  "https://api.figma.com/v1/files/$FILE_KEY/nodes?ids=$NODE_ID_URL" \
  | jq -c ".nodes[\"$NODE_ID_JQ\"].document.children[] | {id, name, type}")

# print out each URL to stdout
echo "$response" | while read -r line; do
    id=$(echo "$line" | jq -r '.id')
    # replace : with -
    id=${id//:/-}
    echo "https://www.figma.com/design/$FILE_KEY/$FILE_NAME?node-id=$id"
done