package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/polpettone/written/pkg"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"strings"
)

const SPACE = " "
const commandOverview = "CTRL+T: Tabs, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open"

func MainView(documents []*models.Document) {
	app := tview.NewApplication()

	tagInputField := tview.NewInputField().
		SetLabel("Tags: ").
		SetFieldBackgroundColor(tcell.Color240)

	queryInputField := tview.NewInputField().
		SetLabel("Query: ").
		SetFieldBackgroundColor(tcell.Color240)

	documentContentView := tview.NewTextView()
	documentMetaInfoView := tview.NewTextView()
	commandOverviewView := tview.NewTextView()
	commandOverviewView.SetText(commandOverview)

	documentTable := buildDocumentTable(documents, documentContentView, documentMetaInfoView, tagInputField)
	fillDocumentTable(documents, *documentTable)

	tagInputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := tagInputField.GetText()
			config.Log.InfoLog.Printf(text)
			selectedRow, _ := documentTable.GetSelection()
			document := documents[selectedRow]
			document.Tags = strings.Split(text, SPACE)
			config.Log.InfoLog.Printf(document.Info.Name())
			config.Log.InfoLog.Printf("%v", documents)
			fillDocumentTable(documents, *documentTable)
		}
	})

	documentGrid := tview.NewGrid().
		SetRows(10, 0, 2).
		SetBorders(true).
		AddItem(documentContentView, 1, 0, 1, 2, 0, 0, false).
		AddItem(documentMetaInfoView, 0, 0, 1, 2, 10, 0, false).
		AddItem(tagInputField, 2, 0, 1, 2, 10, 0, false)

	grid := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetColumns(70, 0).
		SetBorders(true)

	grid.
		AddItem(queryInputField, 0, 0, 1, 2, 0, 100, false).
		AddItem(documentTable, 1, 0, 1, 1, 0, 100, false).
		AddItem(documentGrid, 1, 1, 1, 1, 0, 100, false).
		AddItem(commandOverviewView, 2, 0, 1, 2, 0, 100, false)

	if err := app.
		SetRoot(grid, true).
		SetFocus(documentTable).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				if event.Key() == tcell.KeyCtrlC {
					metaDataPath := viper.GetString(config.MetaDataPath)
					err := service.Save(metaDataPath, documents)
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					if err != nil {
						config.Log.ErrorLog.Printf("%s", err)
					}
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(documentTable)
				}
				if event.Key() == tcell.KeyCtrlT {
					config.Log.DebugLog.Printf("Key: %s", event.Key())
					app.SetFocus(tagInputField)
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

