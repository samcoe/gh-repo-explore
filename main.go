package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
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
	flag.Usage = usageFunc
	flag.StringVar(&branch, "branch", "", "Explore a specific branch of the repository")
	flag.StringVar(&hostname, "hostname", "", "The GitHub hostname for the request (default \"github.com\")")
	flag.Parse()
	repo = flag.Arg(0)

	if repo == "" || repo == "help" {
		usageFunc()
		return nil
	}

	opts := api.ClientOptions{}
	if hostname != "" {
		opts.Host = hostname
	}
	client, err := gh.RESTClient(&opts)
	if err != nil {
		return err
	}

	if branch == "" {
		branch, err = explore.RetrieveDefaultBranch(client, repo)
		if err != nil {
			return err
		}
	}

	gitTree, err := explore.RetrieveGitTree(client, repo, branch)
	if err != nil {
		return err
	}

	fileView := explore.BuildFileView()
	treeView := explore.BuildTreeView(repo, gitTree)
	treeView.SetSelectedFunc(explore.SelectTreeNode(client, repo, branch, fileView))
	searchView := explore.BuildSearchView(repo)
	searchView.SetChangedFunc(explore.SearchTreeView(repo, gitTree, treeView))

	app := buildApplication(treeView, fileView, searchView)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		return event
	})
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

func usageFunc() {
	fmt.Fprintf(os.Stdout, "Interactively explore a repository.\n\n")
	fmt.Fprintf(os.Stdout, "USAGE\n")
	fmt.Fprintf(os.Stdout, "  gh repo-explore <owner>/<repository>\n\n")
	fmt.Fprintf(os.Stdout, "FLAGS\n")
	fmt.Fprintf(os.Stdout, "  --branch\tExplore a specific branch of the repository\n")
	fmt.Fprintf(os.Stdout, "  --hostname\tThe GitHub host for the request (default \"github.com\")\n")
}
