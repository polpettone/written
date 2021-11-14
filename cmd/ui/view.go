package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
)

const SPACE = " "
const commandOverview = "CTRL+Q: Query, CTRL+F: Filter, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open, CTRL+R: Refresh"

func MainView(documents []*models.Document) {
	app := tview.NewApplication()

	queryInputField := tview.NewInputField().
		SetLabel("Query: ").
		SetFieldBackgroundColor(tcell.Color240)

	filterInputField := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldBackgroundColor(tcell.Color240)

	documentContentView := tview.NewTextView()
	documentMetaInfoView := tview.NewTextView()
	commandOverviewView := tview.NewTextView()
	commandOverviewView.SetText(commandOverview)

	documentTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	documentView := DocumentView{
		Documents:    []*models.Document{},
		Table:        documentTable,
		ContentView:  documentContentView,
		MetaInfoView: documentMetaInfoView,
	}

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
		AddItem(documentMetaInfoView, 0, 0, 1, 2, 10, 0, false).
		AddItem(documentContentView, 1, 0, 1, 2, 0, 0, false)

	grid := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetColumns(70, 0).
		SetBorders(true)

	grid.
		AddItem(queryInputField, 0, 0, 1, 1, 0, 100, false).
		AddItem(filterInputField, 0, 1, 1, 1, 0, 100, false).
		AddItem(documentTable, 1, 0, 1, 1, 0, 100, false).
		AddItem(documentGrid, 1, 1, 1, 1, 0, 100, false).
		AddItem(commandOverviewView, 2, 0, 1, 2, 0, 100, false)

	if err := app.
		SetRoot(grid, true).
		SetFocus(documentTable).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				if event.Key() == tcell.KeyCtrlC {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(documentTable)
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
					selectedRow, _ := documentTable.GetSelection()
					document := documents[selectedRow]
					err := pkg.OpenFileInTerminator(document.Path)
					if err != nil {
						config.Log.ErrorLog.Printf("%s", err)
					}
				}

				return event
			}).
		Run(); err != nil {
		panic(err)
	}
}
