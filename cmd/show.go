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
	simpleGrid()
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

func simpleGrid() {

	config.Log.InfoLog.Printf("Simple Grid")
	WrittenDirectory := viper.GetString(WrittenDirectory)
	app := tview.NewApplication()

	documents, _ := readDocuments()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	documentContentView := tview.NewTextView()

	documentList := tview.NewList()
	for _, document := range documents {
		documentList.AddItem(document.Name, "", ' ', nil)
	}

	documentList.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			document := documents[index]
			bytes, err := ioutil.ReadFile(WrittenDirectory + "/" + document.Name)
			if err != nil {
				config.Log.ErrorLog.Printf("{}", err)
			}
			documentContentView.SetText(string(bytes))
		})

	grid := tview.NewGrid().
		SetRows(2, 0, 2).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 2, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(documentList, 0, 0, 0, 0, 0, 0, false).
		AddItem(documentContentView, 1, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(documentList, 1, 0, 1, 1, 0, 100, false).
		AddItem(documentContentView, 1, 1, 1, 1, 0, 100, false)

	if err := app.
				SetRoot(grid, true).
				SetFocus(documentList).
				SetInputCapture(
					func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyCtrlC {
							app.Stop()
						}
						return event
					}).
				Run(); err != nil {
		panic(err)
	}
}

