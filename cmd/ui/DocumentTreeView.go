package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/service"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path/filepath"
)

func TreeView(contentField, metaField, historyField *tview.TextView, rootDir string) *tview.TreeView {
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}

		config.Log.InfoLog.Printf(node.GetText())

		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			file, _ := os.Open(path)
			fileInfo, _ := file.Stat()
			if fileInfo.IsDir() {
				add(node, path)
			} else {
				document, _ := service.ReadDocument(path)
				contentField.SetText(document.Content)
				metaField.SetText(documentMetaView(*document))
				history, err := pkg.GetHistory(rootDir, document.Info.Name())
				if err == nil {
					historyField.SetText(fmt.Sprintf("%s", history))
				} else {
					historyField.SetText("no .git in root dir. no history available")
				}
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	return tree
}
