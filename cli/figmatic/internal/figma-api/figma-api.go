package figma_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type FigmaNode struct {
	URL       string
	LayerName string
	ID        string
}

type GetFileNodeDataResponse struct {
	Name  string `json:"name"`
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

func getFileNodeDataAPI(fileKey, nodeID, figmaAccessToken string) (GetFileNodeDataResponse, error) {
	// Build API URL
	apiURL := fmt.Sprintf("https://api.figma.com/v1/files/%s/nodes?ids=%s", fileKey, nodeID)

	// Make the HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return GetFileNodeDataResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Figma-Token", figmaAccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetFileNodeDataResponse{}, fmt.Errorf("failed to call Figma API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return GetFileNodeDataResponse{}, fmt.Errorf("figma API error: %s", string(body))
	}

	// Parse the JSON response
	var response GetFileNodeDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GetFileNodeDataResponse{}, fmt.Errorf("failed to decode Figma API response: %w", err)
	}

	return response, nil
}

// GetFileNodeData gets the data for an individual node in a figma design. Makes an API call, so avoid using in a loop.
func GetFileNodeData(figmaDesignURL string, figmaAccessToken string) (FigmaNode, error) {
	vals, err := parseValuesFromURL(figmaDesignURL)
	if err != nil {
		return FigmaNode{}, err
	}

	resp, err := getFileNodeDataAPI(vals.FileKey, vals.NodeIDURL, figmaAccessToken)
	if err != nil {
		return FigmaNode{}, err
	}

	return FigmaNode{
		URL:       figmaDesignURL,
		LayerName: resp.Name,
		ID:        vals.NodeIDURL,
	}, nil
}

type figmaNodeValues struct {
	FileKey    string
	FileName   string
	NodeIDURL  string
	NodeIDJSON string
}

func parseValuesFromURL(figmaDesignURL string) (figmaNodeValues, error) {
	// Parse the URL
	u, err := url.Parse(figmaDesignURL)
	if err != nil {
		return figmaNodeValues{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Extract file key and file name from the path
	// Path is /design/<FILE_KEY>/<FILE_NAME>
	parts := strings.Split(strings.TrimPrefix(u.Path, "/design/"), "/")
	if len(parts) < 2 {
		return figmaNodeValues{}, fmt.Errorf("unexpected path format in URL: %s", u.Path)
	}
	fileKey := parts[0]
	fileName := parts[1]

	// Extract node-id from query
	nodeID := u.Query().Get("node-id")
	if nodeID == "" {
		return figmaNodeValues{}, fmt.Errorf("node-id not found in URL query")
	}

	// For Figma API, node ID uses dashes, but in the JSON response, colons are used
	nodeIDJQ := strings.ReplaceAll(nodeID, "-", ":")

	return figmaNodeValues{
		FileKey:    fileKey,
		FileName:   fileName,
		NodeIDURL:  nodeID,
		NodeIDJSON: nodeIDJQ,
	}, nil
}

// GetNodeURLs takes a Figma design URL and access token, and returns URLs for each child node under the specified node.
func GetNodeURLs(figmaDesignURL string, figmaAccessToken string) ([]FigmaNode, error) {
	vals, err := parseValuesFromURL(figmaDesignURL)
	if err != nil {
		return []FigmaNode{}, err
	}

	apiResp, err := getFileNodeDataAPI(vals.FileKey, vals.NodeIDJSON, figmaAccessToken)
	if err != nil {
		return []FigmaNode{}, err
	}

	// Find the correct node using nodeIDJQ
	node, ok := apiResp.Nodes[vals.NodeIDJSON]
	if !ok {
		return []FigmaNode{}, fmt.Errorf("node %s not found in API response", vals.NodeIDJSON)
	}

	// Build URLs for each child
	var nodes []FigmaNode
	for _, child := range node.Document.Children {
		childID := strings.ReplaceAll(child.ID, ":", "-")
		childURL := fmt.Sprintf("https://www.figma.com/design/%s/%s?node-id=%s", vals.FileKey, vals.FileName, childID)
		nodes = append(nodes, FigmaNode{
			URL:       childURL,
			LayerName: child.Name,
			ID:        child.ID,
		})
	}

	return nodes, nil
}
