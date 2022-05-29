package explore

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

type Application struct {
	branch   string
	client   api.RESTClient
	err      error
	fileView *fileView
	height   int
	host     string
	ready    bool
	repo     string
	spinner  spinner.Model
	treeView *treeView
	width    int
}

type gitTreeDownloadMsg struct {
	tree *GitTree
	err  error
}

func NewApplication(host, repo, branch string) (*Application, error) {
	opts := api.ClientOptions{}
	if host != "" {
		opts.Host = host
	}
	client, err := gh.RESTClient(&opts)
	if err != nil {
		return nil, err
	}
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return &Application{
		branch:   branch,
		client:   client,
		fileView: NewFileView(),
		host:     host,
		repo:     repo,
		spinner:  s,
		treeView: NewTreeView(),
	}, nil
}

func (a *Application) Init() tea.Cmd {
	return tea.Batch(
		downloadGitTree(a.client, a.repo, a.branch),
		a.spinner.Tick,
		a.fileView.Init(),
		a.treeView.Init(),
	)
}

func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
		a.treeView.SetSize(msg.Width/2, msg.Height)
		a.fileView.SetSize(msg.Width/2, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return a, tea.Quit
		}
	case gitTreeDownloadMsg:
		a.ready = true
		if msg.err != nil {
			a.err = msg.err
			return a, tea.Quit
		}
	}

	if !a.ready {
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd
	}

	_, cmd = a.fileView.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = a.treeView.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *Application) View() string {
	if !a.ready {
		return a.spinner.View()
	}

	leftBox := a.treeView.View()
	rightBox := a.fileView.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	)
}

func downloadGitTree(client api.RESTClient, repo, branch string) func() tea.Msg {
	return func() tea.Msg {
		if branch == "" {
			var err error
			branch, err = RetrieveDefaultBranch(client, repo)
			if err != nil {
				return gitTreeDownloadMsg{err: err}
			}
		}

		gitTree, err := RetrieveGitTree(client, repo, branch)
		if err != nil {
			return gitTreeDownloadMsg{err: err}
		}

		return gitTreeDownloadMsg{tree: &gitTree}
	}
}
