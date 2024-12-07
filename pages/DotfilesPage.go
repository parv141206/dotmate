// This is the main home page!
// It displays the main dotfiles folder

package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
)

var highlightIndex int

// CreateTreeNode creates a tree node for the given directory or file path.
func CreateTreeNode(path string, isDir bool) *tview.TreeNode {
	node := tview.NewTreeNode(filepath.Base(path))
	if isDir {
		node.SetColor(tcell.ColorGreen)
	} else {
		node.SetColor(tcell.ColorWhite)
	}
	return node
}

// BuildTree recursively builds a tree structure from the root directory.
func BuildTree(rootPath string) *tview.TreeNode {
	rootNode := CreateTreeNode(rootPath, true)
	rootNode.SetExpanded(true)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}
		relativePath, _ := filepath.Rel(rootPath, path)
		parentNode := findParentNode(rootNode, filepath.Dir(relativePath))
		if parentNode != nil {
			childNode := CreateTreeNode(path, info.IsDir())
			parentNode.AddChild(childNode)
		}
		return nil
	})

	if err != nil {
		rootNode.SetText("Error: " + err.Error())
	}

	return rootNode
}

// findParentNode finds the parent node for the given directory.
func findParentNode(node *tview.TreeNode, target string) *tview.TreeNode {
	if node.GetText() == target {
		return node
	}
	for _, child := range node.GetChildren() {
		if result := findParentNode(child, target); result != nil {
			return result
		}
	}
	return nil
}

// ShowDirectoryTree displays the directory tree for a given root directory.

func ShowDirectoryTree(rootDir string) *tview.TreeView {
	highlightIndex = 0
	tree := tview.NewTreeView()

	rootNode := BuildTree(rootDir)
	tree.SetRoot(rootNode).SetCurrentNode(rootNode)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			return nil
		case tcell.KeyUp:
			if highlightIndex > 0 {
				highlightIndex--
				tree.SetCurrentNode(getVisibleNodes(rootNode)[highlightIndex])
			}
			return nil
		case tcell.KeyDown:
			if highlightIndex < len(getVisibleNodes(rootNode))-1 {
				highlightIndex++
				tree.SetCurrentNode(getVisibleNodes(rootNode)[highlightIndex])
			}
			return nil
		case tcell.KeyEnter:
			currentNode := getVisibleNodes(rootNode)[highlightIndex]
			if len(currentNode.GetChildren()) > 0 {
				currentNode.SetExpanded(!currentNode.IsExpanded())
				tree.SetRoot(rootNode).SetCurrentNode(currentNode)
			}
			return nil
		}
		return event
	})

	return tree
}

// getVisibleNodes returns a slice of visible nodes in the tree.
func getVisibleNodes(node *tview.TreeNode) []*tview.TreeNode {
	var visibleNodes []*tview.TreeNode
	var walk func(n *tview.TreeNode)
	walk = func(n *tview.TreeNode) {
		visibleNodes = append(visibleNodes, n)
		if n.IsExpanded() {
			for _, child := range n.GetChildren() {
				walk(child)
			}
		}
	}
	walk(node)
	return visibleNodes
}

func PUSH() {
	panic("push")
}

func CreateDotfilesPage(app *tview.Application) *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexColumn)
	page.SetTitle("Dotfiles").SetBorder(true)
	page.AddItem(ShowDirectoryTree("."), 0, 1, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlP {
			PUSH()
		}
		return event
	})
	// Alright so the following is basically all the 'control' buttons:
	commandsInfoList := tview.NewTextView().
		SetText("Commands:\n\nCTRL + p : Pushes the current dotfiles directory to github")
	page.AddItem(commandsInfoList, 0, 1, false)
	return page

}
