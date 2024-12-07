// I know the name explains it all but well.
// This file is responsible for creating a new configuration for dotfiles.

// NOTE
// I know i have mentioned 'dotfile' but it can be both a file or a directory.
// For example you can copy both 'nvim' directory and a simple .bashrc file.
package pages

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// It was a drag to make 2 functions. I did debug a lot to find a way to combine these but it was way to ugly.
// I think this is much better and readable.
// Both CopyFile and CopyDir are called in the ValidateAndCopy function conditionally.

// As you can tell by the code, go for sure has great error handling.
func CopyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func CopyDir(src string, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateAndCopy is a function which validates the source path and copies it to the destination directory.
// Havent tested with protected files but should work fine without root access.

func ValidateAndCopy(sourcePath string, destDir string) error {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error stating source path: %s", sourcePath)
	}

	if info.IsDir() {
		destPath := filepath.Join(destDir, filepath.Base(sourcePath))
		return CopyDir(sourcePath, destPath)
	} else if info.Mode().IsRegular() {
		destPath := filepath.Join(destDir, filepath.Base(sourcePath))
		return CopyFile(sourcePath, destPath)
	} else {
		return fmt.Errorf("how tf is source path neither a file nor a directory: %s", sourcePath)
	}
}

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
			err := ValidateAndCopy(filePath, "./temp/")
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
