package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

func ShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "shows all writtens",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleShowCommand()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleShowCommand() (string, error) {
	mainView()
	return "", nil
}

func readDocuments() ([]models.Document, error) {
	WrittenDirectory := viper.GetString(WrittenDirectory)
	documents, err := service.Read(WrittenDirectory)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func init() {
	showCmd := ShowCmd()
	rootCmd.AddCommand(showCmd)
}

func mainView() {
	app := tview.NewApplication()

	documents, _ := readDocuments()

	documentContentView := tview.NewTextView()
	documentMetaInfoView := tview.NewTextView()

	documentTable := tview.
		NewTable().
		SetBorders(true).
		SetSelectable(true, false).
		SetSelectedFunc(
			func(row int, column int) {
				document := documents[row]
				bytes, err := ioutil.ReadFile(document.Path)
				if err != nil {
					config.Log.ErrorLog.Printf("{}", err)
				}
				documentContentView.SetText(string(bytes))
				documentMetaInfoView.SetText(fmt.Sprintf("%s \n %s", document.Path, document.Info.Name()))
			},
		)

	for row, document := range documents {
		documentTable.SetCell(row, 0, tview.NewTableCell(document.Info.Name()))
		documentTable.SetCell(row, 1, tview.NewTableCell(document.Info.ModTime().String()))
		documentTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%s", document.Tags)))
	}

	tagInputField := tview.NewInputField().
		SetLabel("Tag").
		SetFieldBackgroundColor(tcell.Color240).
		SetDoneFunc(func(key tcell.Key) {

			if key == tcell.KeyEnter {

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
		SetColumns(100, 0).
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
					service.Save()
					app.Stop()
				}

				if event.Key() == tcell.KeyCtrlD {
					app.SetFocus(documentTable)
				}
				if event.Key() == tcell.KeyCtrlT {
					app.SetFocus(tagInputField)
				}

				return event
			}).
		Run(); err != nil {
		panic(err)
	}
}
