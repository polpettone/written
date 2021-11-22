package ui

import (
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

func MainView() {
	app := tview.NewApplication()

	queryInputField := tview.NewInputField().
		SetLabel("Query: ").
		SetFieldBackgroundColor(tcell.Color240)

	filterInputField := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldBackgroundColor(tcell.Color240)

	commandOverviewView := tview.NewTextView()
	commandOverviewView.SetText(commandOverview)

	documentView := NewDocumentView()

	queryInputField.
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {
				config.Log.DebugLog.Printf("Key: %s", event.Key())
				query := queryInputField.GetText()
				config.Log.InfoLog.Printf("Filter Documents by: s%", query)
				documentView.update(queryInputField.GetText())
				return event
			},
		)

	documentView.update(queryInputField.GetText())

	documentGrid := tview.NewGrid().
		SetRows(10, 0).
		SetBorders(true).
		AddItem(documentView.MetaInfoView, 0, 0, 1, 2, 10, 0, false).
		AddItem(documentView.ContentView, 1, 0, 1, 2, 0, 0, false)

	grid := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetColumns(70, 0).
		SetBorders(true)

	grid.
		AddItem(queryInputField, 0, 0, 1, 1, 0, 100, false).
		AddItem(filterInputField, 0, 1, 1, 1, 0, 100, false).
		AddItem(documentView.Table, 1, 0, 1, 1, 0, 100, false).
		AddItem(documentGrid, 1, 1, 1, 1, 0, 100, false).
		AddItem(commandOverviewView, 2, 0, 1, 2, 0, 100, false)

	if err := app.
		SetRoot(grid, true).
		SetFocus(documentView.Table).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				if event.Key() == tcell.KeyCtrlC {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(documentView.Table)
				}

				if event.Key() == tcell.KeyCtrlQ {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(queryInputField)
				}

				if event.Key() == tcell.KeyCtrlF {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(filterInputField)
				}

				if event.Key() == tcell.KeyCtrlR {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					documentView.update(queryInputField.GetText())
				}

				if event.Key() == tcell.KeyCtrlO {

					selectedRow, _ := documentView.Table.GetSelection()

					config.Log.DebugLog.Printf("selected Row %d", selectedRow)
					config.Log.DebugLog.Printf("len documents %d", len(documentView.Documents))

					if selectedRow < len(documentView.Documents) {
						document := documentView.Documents[selectedRow]
						err := pkg.OpenFileInTerminator(document.Path)
						if err != nil {
							config.Log.ErrorLog.Printf("%s", err)
						}
					}
				}

				return event
			}).
		Run(); err != nil {
		panic(err)
	}
}

func FlexView() {

	writtenDirectory := viper.GetString(config.WrittenDirectory)
	metaField := tview.NewTextView()
	contentField := tview.NewTextView()
	tree := TreeView(contentField, metaField, writtenDirectory)

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(metaField, 0, 1, false).
			AddItem(contentField, 0, 3, false).
			AddItem(metaField, 5, 1, false), 0, 2, false)
	if err := app.SetRoot(flex, true).SetFocus(tree).Run(); err != nil {
		panic(err)
	}
}

func TreeView(contentField, metaField *tview.TextView, rootDir string) *tview.TreeView {
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
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tree
}
