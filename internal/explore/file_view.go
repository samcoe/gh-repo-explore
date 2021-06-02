package explore

import "github.com/rivo/tview"

func BuildFileView() *tview.TextView {
	fileView := tview.NewTextView()
	fileView.SetDynamicColors(true)
	fileView.SetBorder(true)
	return fileView
}
