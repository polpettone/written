package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

const SPACE = " "
const commandOverview = "CTRL+Q: Query, CTRL+F: Filter, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open, CTRL+R: Refresh"

func FlexView() {

	writtenDirectory := viper.GetString(config.WrittenDirectory)
	documentMetaField := tview.NewTextView()
	documentField := tview.NewTextView()
	documentHistoryField := tview.NewTextView()

	tree := TreeView(documentField, documentMetaField, documentHistoryField, writtenDirectory)

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(documentMetaField, 0, 1, false).
			AddItem(documentField, 0, 3, false).
			AddItem(documentHistoryField, 5, 1, false), 0, 2, false)


	if err := app.SetRoot(flex, true).
		SetFocus(tree).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				config.Log.DebugLog.Printf("Key: %s", event.Key())

				if event.Key() == tcell.KeyCtrlC {
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlT {
					app.SetFocus(tree)
				}

				if event.Key() == tcell.KeyCtrlH {
					app.SetFocus(documentHistoryField)
				}

				if event.Key() == tcell.KeyCtrlD {
					app.SetFocus(documentField)
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
