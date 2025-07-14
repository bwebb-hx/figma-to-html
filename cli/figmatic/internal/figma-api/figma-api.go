package figma_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetNodeURLs takes a Figma design URL and access token, and returns URLs for each child node under the specified node.
func GetNodeURLs(figmaDesignURL string, figmaAccessToken string) ([]string, error) {
	// Parse the URL
	u, err := url.Parse(figmaDesignURL)
	if err != nil {
		return []string{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Extract file key and file name from the path
	// Path is /design/<FILE_KEY>/<FILE_NAME>
	parts := strings.Split(strings.TrimPrefix(u.Path, "/design/"), "/")
	if len(parts) < 2 {
		return []string{}, fmt.Errorf("unexpected path format in URL: %s", u.Path)
	}
	fileKey := parts[0]
	fileName := parts[1]

	// Extract node-id from query
	nodeID := u.Query().Get("node-id")
	if nodeID == "" {
		return []string{}, fmt.Errorf("node-id not found in URL query")
	}

	// For Figma API, node ID uses dashes, but in the JSON response, colons are used
	nodeIDJQ := strings.ReplaceAll(nodeID, "-", ":")

	// Build API URL
	apiURL := fmt.Sprintf("https://api.figma.com/v1/files/%s/nodes?ids=%s", fileKey, nodeID)

	// Make the HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return []string{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Figma-Token", figmaAccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("failed to call Figma API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return []string{}, fmt.Errorf("figma API error: %s", string(body))
	}

	// Parse the JSON response
	var apiResp struct {
		Nodes map[string]struct {
			Document struct {
				Children []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"children"`
			} `json:"document"`
		} `json:"nodes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return []string{}, fmt.Errorf("failed to decode Figma API response: %w", err)
	}

	// Find the correct node using nodeIDJQ
	node, ok := apiResp.Nodes[nodeIDJQ]
	if !ok {
		return []string{}, fmt.Errorf("node %s not found in API response", nodeIDJQ)
	}

	// Build URLs for each child
	var urls []string
	for _, child := range node.Document.Children {
		childID := strings.ReplaceAll(child.ID, ":", "-")
		childURL := fmt.Sprintf("https://www.figma.com/design/%s/%s?node-id=%s", fileKey, fileName, childID)
		urls = append(urls, childURL)
	}

	return urls, nil
}
