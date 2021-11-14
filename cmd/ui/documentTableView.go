package ui

import (
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/rivo/tview"
	"io/ioutil"
	"sort"
	"strings"
)

func fillDocumentTable(documents []*models.Document, documentTable tview.Table) {

	sort.Slice(documents, func(i, j int) bool {
		return documents[i].Info.ModTime().After(documents[j].Info.ModTime())
	})

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
				content, err := ioutil.ReadFile(document.Path)
				tags := service.ExtractTags(string(content))
				document.Tags = tags
				if err != nil {
					config.Log.ErrorLog.Printf("{}", err)
				}
				documentContentView.SetText(string(content))
				documentMetaInfoView.SetText(documentMetaView(*document))
			},
		)
	return documentTable
}