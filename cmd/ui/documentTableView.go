package ui

import (
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/rivo/tview"
	"io/ioutil"
	"strings"
)

func fillDocumentTable(documents []*models.Document, documentTable tview.Table) {
	for row, document := range documents {
		documentTable.SetCell(row, 0, tview.NewTableCell(document.Info.Name()))
		documentTable.SetCell(row, 1, tview.NewTableCell(strings.Join(document.Tags, SPACE)))
	}
}

func buildDocumentTable(documents []*models.Document,
	documentContentView *tview.TextView,
	documentMetaInfoView *tview.TextView) *tview.Table {

	documentTable := tview.NewTable().
		SetBorders(false).
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
			},
		)
	return documentTable
}