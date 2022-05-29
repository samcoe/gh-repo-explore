package explore

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			PaddingTop(1)
	}()

	listStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.NormalBorder())
	}()
)

type treeView struct {
	active   bool
	delegate list.DefaultDelegate
	list     list.Model
	tree     *GitTree
}

func NewTreeView() *treeView {
	delegate := list.NewDefaultDelegate()
	list := list.New([]list.Item{}, delegate, 0, 0)
	list.Title = "Filetree"
	return &treeView{
		delegate: delegate,
		list:     list,
	}
}

func (t *treeView) SetSize(width, height int) {
	// horizontal, vertical := listStyle.GetFrameSize()
	// t.list.Styles.StatusBar.Width(width - horizontal)
	// t.list.SetSize(width-horizontal, height-vertical)
	t.list.SetSize(width, height)
}

func (t *treeView) Init() tea.Cmd {
	return nil
}

func (t *treeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case gitTreeDownloadMsg:
		t.tree = msg.tree
		names := make([]list.Item, len(*msg.tree))
		for i, s := range *msg.tree {
			names[i] = &s
		}
		cmd = t.list.SetItems(names)
		cmds = append(cmds, cmd)
	}

	return t, tea.Batch(cmds...)
}

func (t *treeView) View() string {
	return listStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			t.list.View(),
		),
	)
}

// func SearchTreeView(repo string, gt GitTree, treeView *tview.TreeView) func(string) {
// 	return func(query string) {
// 		root := treeView.GetRoot()
// 		root.ClearChildren()
// 		dirs := map[string]*tview.TreeNode{".": root}
// 		r, err := regexp.Compile("(?i)" + query)
// 		if err != nil {
// 			return
// 		}
//
// 		for _, n := range gt {
// 			if n.IsDir() || !r.MatchString(n.Path) {
// 				continue
// 			}
// 			node := tview.NewTreeNode(n.Name())
// 			node.SetReference(n)
// 			parentNode := makeParentNodes(dirs, n.Dir(), len(query) != 0)
// 			parentNode.AddChild(node)
// 		}
// 	}
// }
//
// func SelectTreeNode(client api.RESTClient, repo, branch string, fileView *tview.TextView) func(*tview.TreeNode) {
// 	return func(node *tview.TreeNode) {
// 		reference := node.GetReference()
// 		if reference == nil {
// 			return
// 		}
// 		rtn := reference.(GitTreeNode)
// 		if rtn.IsDir() {
// 			node.SetExpanded(!node.IsExpanded())
// 			return
// 		}
// 		fileBytes, err := RetrieveFileContent(client, repo, branch, rtn.Path)
// 		if err != nil {
// 			return
// 		}
// 		file := string(fileBytes)
// 		fileView.Clear()
// 		coloredFileView := tview.ANSIWriter(fileView)
// 		_ = quick.Highlight(coloredFileView, file, rtn.Ext(), "terminal256", "solarized-dark")
// 		fileView.ScrollToBeginning()
// 	}
// }
//
// func makeParentNodes(dirs map[string]*tview.TreeNode, dir string, expanded bool) *tview.TreeNode {
// 	parentNode := dirs[dir]
// 	if parentNode != nil {
// 		return parentNode
// 	}
// 	parentNode = makeParentNodes(dirs, filepath.Dir(dir), expanded)
// 	node := tview.NewTreeNode(dir)
// 	node.SetReference(GitTreeNode{Path: dir, Type: "tree"})
// 	node.SetColor(tcell.ColorGreen)
// 	node.SetExpanded(expanded)
// 	parentNode.AddChild(node)
// 	dirs[dir] = node
// 	return node
// }
