// I know the name explains it all but well.
// This file is responsible for creating a new configuration for dotfiles.

// NOTE
// I know i have mentioned 'dotfile' but it can be both a file or a directory.
// For example you can copy both 'nvim' directory and a simple .bashrc file.
package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// All the functions used here are written in utils.go file!
func CreateAddDotfilesPage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)
	page.SetTitle("Add Dotfile").SetBorder(true)

	dotFilePath := tview.NewInputField().
		SetFieldBackgroundColor(tcell.NewHexColor(0x041727)).
		SetLabel("Dotfile Path")

	// The notification split is used to display the result of the add dotfile operation.
	// I couldnt find a better name, it is what it is.
	notificationSplit := tview.NewTextView()
	form := tview.NewForm().
		AddFormItem(dotFilePath).
		AddButton("Add Dotfile", func() {
			filePath := dotFilePath.GetText()
			// This path is temporary
			err := ValidateAndSync(filePath, "./temp/")
			if err != nil {
				notificationSplit.SetText(err.Error())
			} else {
				notificationSplit.SetText("Added Dotfile")
			}
		}).
		SetButtonBackgroundColor(tcell.NewHexColor(0x041727)).
		SetFieldBackgroundColor(tcell.NewHexColor(0x041727))

	page.AddItem(form, 0, 1, true)
	page.AddItem(notificationSplit, 1, 1, false)
	return page
}
