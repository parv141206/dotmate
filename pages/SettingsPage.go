package pages

import (
	"github.com/rivo/tview"
)

func CreateSettingsPage() *tview.Box {
	page := tview.NewBox().SetTitle("Settings").SetBorder(true)
	return page
}
