package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/rivo/tview"

	"github.com/samcoe/gh-repo-explore/internal/explore"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	var err error
	var hostname string
	var repo string
	var branch string
	flag.StringVar(&branch, "branch", "", "Explore a specific branch of the repository")
	flag.StringVar(&hostname, "hostname", "", "The GitHub hostname for the request (default \"github.com\")")
	flag.Parse()
	repo = flag.Arg(0)

	if repo == "" {
		return errors.New("repository argument required")
	}
	if branch == "" {
		branch, err = explore.RetrieveDefaultBranch(hostname, repo)
		if err != nil {
			return err
		}
	}

	gitTree, err := explore.RetrieveGitTree(hostname, repo, branch)
	if err != nil {
		return err
	}

	fileView := explore.BuildFileView()
	treeView := explore.BuildTreeView(repo, gitTree)
	treeView.SetSelectedFunc(explore.SelectTreeNode(hostname, repo, branch, fileView))
	searchView := explore.BuildSearchView(repo)
	searchView.SetChangedFunc(explore.SearchTreeView(repo, gitTree, treeView))

	app := buildApplication(treeView, fileView, searchView)
	return app.Run()
}

func buildApplication(treeView *tview.TreeView, fileView *tview.TextView, searchView *tview.InputField) *tview.Application {
	app := tview.NewApplication()
	topRow := tview.NewFlex().
		AddItem(treeView, 0, 1, false).
		AddItem(fileView, 0, 4, false)
	bottomRow := tview.NewFlex().
		AddItem(searchView, 0, 1, false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 17, false).
		AddItem(bottomRow, 0, 1, false)
	app.SetRoot(flex, true).EnableMouse(true).SetFocus(searchView)
	return app
}
