package pages

import (
	"github.com/rivo/tview"
)

func CreateSettingsPage() *tview.Flex {
	page := tview.NewFlex()
	page.SetTitle("Settings").SetBorder(true)
	return page
}
