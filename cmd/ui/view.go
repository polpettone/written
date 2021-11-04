package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/rivo/tview"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"strings"
	"time"
)

const SPACE = " "

func fillDocumentTable(documents []*models.Document, documentTable tview.Table) {
	for row, document := range documents {
		documentTable.SetCell(row, 0, tview.NewTableCell(document.Info.Name()))
		documentTable.SetCell(row, 1, tview.NewTableCell(document.Info.ModTime().Format(time.RFC822)))
		documentTable.SetCell(row, 2, tview.NewTableCell(strings.Join(document.Tags, SPACE)))
	}

}

func buildDocumentTable(documents []*models.Document,
	documentContentView *tview.TextView,
	documentMetaInfoView *tview.TextView,
	tagInputField *tview.InputField) *tview.Table {
	documentTable := tview.
		NewTable().
		SetBorders(true).
		SetSelectable(true, false).
		SetSelectionChangedFunc(
			func(row int, column int) {
				document := documents[row]
				bytes, err := ioutil.ReadFile(document.Path)
				if err != nil {
					config.Log.ErrorLog.Printf("{}", err)
				}
				documentContentView.SetText(string(bytes))
				documentMetaInfoView.SetText(fmt.Sprintf("%s \n %s", document.Path, document.Info.Name()))
				tagInputField.SetText(fmt.Sprintf("%s", strings.Join(document.Tags, SPACE)))
			},
		)
	return documentTable
}

func MainView(documents []*models.Document) {
	app := tview.NewApplication()

	tagInputField := tview.NewInputField().
		SetLabel("Tags: ").
		SetFieldBackgroundColor(tcell.Color240)

	documentContentView := tview.NewTextView()
	documentMetaInfoView := tview.NewTextView()
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
		SetRows(2, 0, 2).
		SetColumns(70, 0).
		SetBorders(true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(documentTable, 0, 0, 0, 0, 0, 0, false).
		AddItem(documentGrid, 1, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(documentTable, 1, 0, 1, 1, 0, 100, false).
		AddItem(documentGrid, 1, 1, 1, 1, 0, 100, false)

	if err := app.
		SetRoot(grid, true).
		SetFocus(documentTable).
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {

				if event.Key() == tcell.KeyCtrlC {
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
					app.SetFocus(documentTable)
				}
				if event.Key() == tcell.KeyCtrlT {
					app.SetFocus(tagInputField)
				}

				if event.Key() == tcell.KeyCtrlO {
					selectedRow, _ := documentTable.GetSelection()
					document := documents[selectedRow]
					config.Log.InfoLog.Printf("%s", document.Path)
					err := open.RunWith(document.Path, "vim")
					if err != nil {
						config.Log.ErrorLog.Printf("%v", err)
					}
				}

				return event
			}).
		Run(); err != nil {
		panic(err)
	}
}
