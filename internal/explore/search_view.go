package explore

import "github.com/rivo/tview"

func BuildSearchView(repo string) *tview.InputField {
	searchView := tview.NewInputField()
	searchView.SetLabel(repo)
	searchView.SetFieldBackgroundColor(0)
	searchView.SetLabelWidth(len(repo) + 1)
	return searchView
}
