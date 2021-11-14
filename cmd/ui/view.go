package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
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
