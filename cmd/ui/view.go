package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

const SPACE = " "
const commandOverview = "CTRL+Q: Query, CTRL+F: Filter, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open, CTRL+R: Refresh"


func FlexView() {

	wide := true

	writtenDirectory := viper.GetString(config.WrittenDirectory)
	documentMetaField := tview.NewTextView()
	documentField := tview.NewTextView()
	documentHistoryField := tview.NewTextView()

	tree := TreeView(documentField, documentMetaField, documentHistoryField, writtenDirectory)

	app := tview.NewApplication()

	rightPanelFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(documentMetaField, 0, 1, false).
		AddItem(documentHistoryField, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(tree, 0, 2, false).
		AddItem(documentField, 0, 5, false).
		AddItem(rightPanelFlex, 0, 2, false)

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

				if event.Key() == tcell.KeyCtrlW {
					if !wide  {
						flex.AddItem(rightPanelFlex, 0, 2, false)
						wide = true
					}
				}

				if event.Key() == tcell.KeyCtrlN {
					flex.RemoveItem(rightPanelFlex)
					wide = false
				}

				if event.Key() == tcell.KeyCtrlO {
					node := tree.GetCurrentNode()
					if node.GetLevel() != 0 {
						err := pkg.OpenFileInTerminator(node.GetReference().(string))
						if err != nil {
							config.Log.ErrorLog.Printf("", err)
						}
					}
				}
				return event
			}).
		Run(); err != nil {

		panic(err)
	}
}
