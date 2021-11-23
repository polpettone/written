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
