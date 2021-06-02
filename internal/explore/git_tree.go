package explore

import "path/filepath"

type GitTree []GitTreeNode

type GitTreeNode struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

func (gtn *GitTreeNode) IsDir() bool {
	return gtn.Type == "tree"
}

func (gtn *GitTreeNode) Name() string {
	return filepath.Base(gtn.Path)
}

func (gtn *GitTreeNode) Dir() string {
	return filepath.Dir(gtn.Path)
}

func (gtn *GitTreeNode) Ext() string {
	return filepath.Ext(gtn.Path)
}
