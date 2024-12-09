package pages

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HERE IS THE CONSTANT commands info
var commandsInfoList = `
Navigation:
  ↑ ↓ : Navigate up and down the treeview
  Enter : Open the selected directory
  
  1 : To navigate to home page
  2 : To add new dotfile
  3 : To navigate to settings page

Commands: 
  CTRL + r : Pull all the changes from all the dotfiles locally stored.
  CTRL + p : Pushes the current dotfiles directory to github
`

// BELOW IS ALL THE TREEVIEW LOGIC
var highlightIndex int

// BuildTree recursively builds a tree structure from the root directory.
func BuildTree(rootPath string) *tview.TreeNode {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		rootNode := tview.NewTreeNode("Error: " + err.Error())
		return rootNode
	}

	rootNode := CreateTreeNode(rootPath, true)
	rootNode.SetExpanded(false) // Root node is collapsed by default

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == rootPath {
			return nil
		}

		parentPath := filepath.Dir(path)
		parentNode := findParentNode(rootNode, parentPath)
		if parentNode != nil {
			childNode := CreateTreeNode(path, info.IsDir())
			childNode.SetExpanded(false) // All nodes are collapsed by default
			parentNode.AddChild(childNode)
		}
		return nil
	})

	if err != nil {
		rootNode.SetText("Error: " + err.Error())
	}

	return rootNode
}

// findParentNode finds the parent node for the given directory based on its absolute path.
func findParentNode(node *tview.TreeNode, targetAbsPath string) *tview.TreeNode {
	nodePath := node.GetReference().(string)
	if nodePath == targetAbsPath {
		return node
	}

	for _, child := range node.GetChildren() {
		if result := findParentNode(child, targetAbsPath); result != nil {
			return result
		}
	}
	return nil
}

// CreateTreeNode creates a tree node for the given directory or file path.
func CreateTreeNode(path string, isDir bool) *tview.TreeNode {
	node := tview.NewTreeNode(filepath.Base(path)).
		SetReference(path)

	if isDir {
		node.SetColor(tcell.ColorGreen)
	} else {
		node.SetColor(tcell.ColorWhite)
	}
	return node
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

func CreateDotfilesPage(app *tview.Application) *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexColumn)
	page.SetTitle("Dotfiles").SetBorder(true)

	rootDir := getDestinationDir()
	tree := ShowDirectoryTree(rootDir)
	page.AddItem(tree, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlP:
			PUSH()
		case tcell.KeyCtrlR:
			PULL(tree, rootDir)
		}
		return event
	})

	commandsInfoListView := tview.NewTextView().
		SetText(commandsInfoList)
	page.AddItem(commandsInfoListView, 0, 1, false)

	return page
}

func PUSH() {
	panic("push")
}

// This function re copies all the dotfiles (whose location is in file_paths.json) to the destination directory. So basically it pulls all the changes from the configs. I could do something like GIT but ill do it later.

func PULL(tree *tview.TreeView, rootDir string) {
	destDir := getDestinationDir()
	paths, err := loadFilePaths()
	if err != nil {
		log.Fatalf("Error loading file paths: %v", err)
	}
	for path := range paths {
		err = ValidateAndSync(path, destDir)
		if err != nil {
			log.Fatalf("Error syncing files: %v", err)
		}
	}

	rootNode := BuildTree(rootDir)
	tree.SetRoot(rootNode).SetCurrentNode(rootNode)
	highlightIndex = 0
}
