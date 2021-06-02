package explore

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
)

func RetrieveDefaultBranch(hostname, repo string) (string, error) {
	var response struct {
		DefaultBranch string `json:"default_branch"`
	}

	args := []string{}
	path := fmt.Sprintf("/repos/%s", repo)
	args = append(args, path)
	if hostname != "" {
		host := fmt.Sprintf("--hostname=%s", hostname)
		args = append(args, host)
	}

	err := apiGet(args, &response)
	if err != nil {
		return "", err
	}

	return response.DefaultBranch, nil
}

// https://docs.github.com/en/rest/reference/git#get-a-tree
func RetrieveGitTree(hostname, repo, branch string) (GitTree, error) {
	var response struct {
		SHA       string  `json:"sha"`
		URL       string  `json:"url"`
		Tree      GitTree `json:"tree"`
		Truncated bool    `json:"truncated"`
	}

	args := []string{}
	path := fmt.Sprintf("/repos/%s/git/trees/%s?recursive=1", repo, branch)
	args = append(args, path)
	if hostname != "" {
		host := fmt.Sprintf("--hostname=%s", hostname)
		args = append(args, host)
	}

	err := apiGet(args, &response)
	if err != nil {
		return nil, err
	}

	return response.Tree, nil
}

func RetrieveFileContent(hostname, repo, branch, filePath string) ([]byte, error) {
	var response struct {
		Content string `json:"content"`
	}

	args := []string{}
	path := fmt.Sprintf("/repos/%s/contents/%s", repo, filePath)
	if branch != "" {
		path += fmt.Sprintf("?ref=%s", url.QueryEscape(branch))
	}
	args = append(args, path)
	if hostname != "" {
		host := fmt.Sprintf("--hostname=%s", hostname)
		args = append(args, host)
	}

	err := apiGet(args, &response)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	return decoded, nil
}

func apiGet(args []string, response interface{}) error {
	exe, err := exec.LookPath("gh")
	if err != nil {
		return err
	}
	args = append([]string{"api"}, args...)
	cmd := exec.Command(exe, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	err = json.Unmarshal(stdout.Bytes(), &response)
	return err
}
