package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/service"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

const SPACE = " "
const commandOverview = "CTRL+Q: Query, CTRL+F: Filter, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open, CTRL+R: Refresh"

func FlexView() {

	writtenDirectory := viper.GetString(config.WrittenDirectory)
	metaField := tview.NewTextView()
	contentField := tview.NewTextView()
	historyField := tview.NewTextView()

	tree := TreeView(contentField, metaField, historyField, writtenDirectory)

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(metaField, 0, 1, false).
			AddItem(contentField, 0, 3, false).
			AddItem(historyField, 5, 1, false), 0, 2, false)
	if err := app.SetRoot(flex, true).
		SetFocus(tree).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				config.Log.DebugLog.Printf("Key: %s", event.Key())

				if event.Key() == tcell.KeyCtrlC {
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
				}

				if event.Key() == tcell.KeyCtrlQ {
				}

				if event.Key() == tcell.KeyCtrlF {
				}

				if event.Key() == tcell.KeyCtrlR {
				}

				if event.Key() == tcell.KeyCtrlO {
				}
				return event
			}).
		Run(); err != nil {


		panic(err)
	}
}

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
				if err != nil {
					panic(err)
				}
				historyField.SetText(fmt.Sprintf("%s", history))
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	return tree
}
