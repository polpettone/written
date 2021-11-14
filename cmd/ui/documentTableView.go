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

func fillDocumentTable(
	documents []*models.Document,
	documentTable *tview.Table,
	documentContentView *tview.TextView,
	documentMetaInfoView *tview.TextView,
	query string) {

	documentTable.Clear()

	filteredDocuments := service.FilterDocuments(documents, query)

	sort.Slice(filteredDocuments, func(i, j int) bool {
		return filteredDocuments[i].Info.ModTime().After(filteredDocuments[j].Info.ModTime())
	})

	for row, document := range filteredDocuments {
		documentTable.SetCell(row, 0, tview.NewTableCell(document.Info.Name()))
		documentTable.SetCell(row, 1, tview.NewTableCell(strings.Join(document.Tags, SPACE)))
	}

	documentTable.SetSelectionChangedFunc(
		func(row int, column int) {
			document := filteredDocuments[row]
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

}
