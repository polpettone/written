package cmd

import (
	"fmt"
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

var currentDocumentContent string

func simpleGrid() {

	config.Log.InfoLog.Printf("Simple Grid")
	WrittenDirectory := viper.GetString(WrittenDirectory)
	app := tview.NewApplication()

	documentList := tview.NewList()

	documents, _ := readDocuments()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	textView := tview.NewTextView()

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
			currentDocumentContent = string(bytes)
			textView.SetText(currentDocumentContent)
		})

	grid := tview.NewGrid().
		SetRows(2, 0, 2).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 2, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(documentList, 0, 0, 0, 0, 0, 0, false).
		AddItem(textView, 1, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(documentList, 1, 0, 1, 1, 0, 100, false).
		AddItem(textView, 1, 1, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).SetFocus(documentList).Run(); err != nil {
		panic(err)
	}
}

func simpleListView() {
	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("List item 1", "Some explanatory text", 'a', nil).
		AddItem("List item 2", "Some explanatory text", 'b', nil).
		AddItem("List item 3", "Some explanatory text", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func simpleBox() {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
