package explore

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/cli/go-gh/pkg/api"
)

func RetrieveDefaultBranch(client api.RESTClient, repo string) (string, error) {
	var response struct {
		DefaultBranch string `json:"default_branch"`
	}

	path := fmt.Sprintf("repos/%s", repo)

	err := client.Get(path, &response)
	if err != nil {
		return "", err
	}

	return response.DefaultBranch, nil
}

func RetrieveGitTree(client api.RESTClient, repo, branch string) (GitTree, error) {
	var response struct {
		SHA       string  `json:"sha"`
		URL       string  `json:"url"`
		Tree      GitTree `json:"tree"`
		Truncated bool    `json:"truncated"`
	}

	path := fmt.Sprintf("repos/%s/git/trees/%s?recursive=1", repo, branch)

	err := client.Get(path, &response)
	if err != nil {
		return nil, err
	}

	return response.Tree, nil
}

func RetrieveFileContent(client api.RESTClient, repo, branch, filePath string) ([]byte, error) {
	var response struct {
		Content string `json:"content"`
	}

	path := fmt.Sprintf("repos/%s/contents/%s", repo, filePath)
	if branch != "" {
		path += fmt.Sprintf("?ref=%s", url.QueryEscape(branch))
	}

	err := client.Get(path, &response)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	return decoded, nil
}
