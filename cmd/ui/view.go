package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
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
				documentMetaInfoView.SetText(documentMetaView(*document))
				tagInputField.SetText(fmt.Sprintf("%s", strings.Join(document.Tags, SPACE)))
			},
		)
	return documentTable
}

const commandOverview = "CTRL+T: Tabs, CTRL+D: Document Table, CTRL+C: Quit, CTRL+O: Open"

func MainView(documents []*models.Document) {
	app := tview.NewApplication()

	tagInputField := tview.NewInputField().
		SetLabel("Tags: ").
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
		SetRows(3, 0, 2).
		SetColumns(70, 0).
		SetBorders(true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(documentTable, 0, 0, 0, 0, 0, 0, false).
		AddItem(documentGrid, 1, 0, 1, 2, 0, 0, false).
		AddItem(commandOverviewView, 2, 0, 0, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(documentTable, 1, 0, 1, 1, 0, 100, false).
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
					openFileInTerminator(document.Path)
				}

				return event
			}).
		Run(); err != nil {
		panic(err)
	}
}

func openFileInTerminator(path string) {
	exectuable := "terminator"
	command := exec.Command(exectuable, "-e", "vim "+path)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		config.Log.ErrorLog.Printf("%s", err)
	}
	config.Log.InfoLog.Printf("open %s in editor", path)
}
