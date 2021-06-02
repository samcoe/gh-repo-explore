package explore

import (
	"path/filepath"
	"regexp"

	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func BuildTreeView(repo string, gt GitTree) *tview.TreeView {
	root := tview.NewTreeNode(repo).SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true)
	SearchTreeView(repo, gt, tree)("")
	return tree
}

func SearchTreeView(repo string, gt GitTree, treeView *tview.TreeView) func(string) {
	return func(query string) {
		root := treeView.GetRoot()
		root.ClearChildren()
		dirs := map[string]*tview.TreeNode{".": root}
		r := regexp.MustCompile(query)

		for _, n := range gt {
			if n.IsDir() || !r.MatchString(n.Path) {
				continue
			}
			node := tview.NewTreeNode(n.Name())
			node.SetReference(n)
			parentNode := makeParentNodes(dirs, n.Dir(), len(query) != 0)
			parentNode.AddChild(node)
		}
	}
}

func SelectTreeNode(hostname, repo, branch string, fileView *tview.TextView) func(*tview.TreeNode) {
	return func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		rtn := reference.(GitTreeNode)
		if rtn.IsDir() {
			node.SetExpanded(!node.IsExpanded())
			return
		}
		fileBytes, err := RetrieveFileContent(hostname, repo, branch, rtn.Path)
		if err != nil {
			return
		}
		file := string(fileBytes)
		fileView.Clear()
		coloredFileView := tview.ANSIWriter(fileView)
		_ = quick.Highlight(coloredFileView, file, rtn.Ext(), "terminal256", "solarized-dark")
		fileView.ScrollToBeginning()
	}
}

func makeParentNodes(dirs map[string]*tview.TreeNode, dir string, expanded bool) *tview.TreeNode {
	parentNode := dirs[dir]
	if parentNode != nil {
		return parentNode
	}
	parentNode = makeParentNodes(dirs, filepath.Dir(dir), expanded)
	node := tview.NewTreeNode(dir)
	node.SetReference(GitTreeNode{Path: dir, Type: "tree"})
	node.SetColor(tcell.ColorGreen)
	node.SetExpanded(expanded)
	parentNode.AddChild(node)
	dirs[dir] = node
	return node
}
