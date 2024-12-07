package pages

import (
	"github.com/rivo/tview"
)

func CreateDotfilesPage() *tview.Box {
	page := tview.NewBox().SetTitle("Dotfiles").SetBorder(true)
	return page
}
