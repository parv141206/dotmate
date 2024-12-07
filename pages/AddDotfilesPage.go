// I know the name explains it all but well.
// This file is responsible for creating a new configuration for dotfiles.

package pages

import (
	"github.com/rivo/tview"
)

func CreateAddDotfilesPage() *tview.Box {
	page := tview.NewBox().SetTitle("Add Dotfile").SetBorder(true)
	return page
}
