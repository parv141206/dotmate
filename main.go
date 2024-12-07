// main.go
// Hello, its me Parv
// For some reason I have decided to document this project completely.
// The comments may or may not be useful but I tried ha!

package main

import (
	// "github.com/gdamore/tcell/v2"
	"dotmate/pages"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pagesObj := tview.NewPages()
	// Declaring 3 pages.
	// The first page shows all the dot files and the second page is a form which allows you to add a new dot file.
	// The third page contains some settings for the application.
	// I have splitted the pages into their own files. Have a look at the pages folder.
	page1 := pages.CreateDotfilesPage()
	page2 := pages.CreateAddDotfilesPage()
	page3 := pages.CreateSettingsPage()

	pagesObj.AddPage("dotfiles", page1, true, true)
	pagesObj.AddPage("add", page2, true, false)
	pagesObj.AddPage("settings", page3, true, false)
	// I have splitted the routing configuration into a separate file. Have a look at the routing.go file.
	SetupRouting(pagesObj, page1, page2, page3)

	app.SetRoot(pagesObj, true).SetFocus(pagesObj)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
