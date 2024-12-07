package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SetupRouting configures the input capture for page navigation.
// I know there could be a better way to do this but I am not that good at golang.
func SetupRouting(pages *tview.Pages, page1, page2, page3 *tview.Flex) {
	page1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '2':
				pages.SwitchToPage("add")
			case '3':
				pages.SwitchToPage("settings")
			}
		}
		return event
	})

	page2.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '1':
				pages.SwitchToPage("dotfiles")
			case '3':
				pages.SwitchToPage("settings")
			}
		}
		return event
	})

	page3.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '1':
				pages.SwitchToPage("dotfiles")
			case '2':
				pages.SwitchToPage("add")
			}
		}
		return event
	})
}

