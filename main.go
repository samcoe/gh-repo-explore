package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

	app, err := explore.NewApplication(hostname, repo, branch)
	if err != nil {
		return err
	}

	var opts []tea.ProgramOption
	opts = append(opts, tea.WithAltScreen(), tea.WithMouseCellMotion())

	p := tea.NewProgram(app, opts...)
	if err := p.Start(); err != nil {
		return err
	}

	return nil
}

func usageFunc() {
	fmt.Fprintf(os.Stdout, "Interactively explore a repository.\n\n")
	fmt.Fprintf(os.Stdout, "USAGE\n")
	fmt.Fprintf(os.Stdout, "  gh repo-explore <owner>/<repository>\n\n")
	fmt.Fprintf(os.Stdout, "FLAGS\n")
	fmt.Fprintf(os.Stdout, "  --branch\tExplore a specific branch of the repository\n")
	fmt.Fprintf(os.Stdout, "  --hostname\tThe GitHub host for the request (default \"github.com\")\n")
}
