package ui

import (
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"io/ioutil"
	"sort"
	"strings"
)

type DocumentView struct {
	Documents    []*models.Document
	Table        *tview.Table
	ContentView  *tview.TextView
	MetaInfoView *tview.TextView
}


func (view DocumentView) update(query string) {
	WrittenDirectory := viper.GetString(config.WrittenDirectory)
	documents, err := service.Read(WrittenDirectory)
	if err != nil {
		config.Log.ErrorLog.Printf("%s", err)
	}
	view.Documents = documents

	view.Table.Clear()

	filteredDocuments := service.FilterDocuments(documents, query)

	sort.Slice(filteredDocuments, func(i, j int) bool {
		return filteredDocuments[i].Info.ModTime().After(filteredDocuments[j].Info.ModTime())
	})

	for row, document := range filteredDocuments {
		view.Table.SetCell(row, 0, tview.NewTableCell(document.Info.Name()))
		view.Table.SetCell(row, 1, tview.NewTableCell(strings.Join(document.Tags, SPACE)))
	}

	view.Table.SetSelectionChangedFunc(
		func(row int, column int) {
			if len(filteredDocuments) > row {
				document := filteredDocuments[row]
				content, err := ioutil.ReadFile(document.Path)
				tags := service.ExtractTags(string(content))
				document.Tags = tags
				if err != nil {
					config.Log.ErrorLog.Printf("{}", err)
				}
				view.ContentView.SetText(string(content))
				view.MetaInfoView.SetText(documentMetaView(*document))
			} else {
				view.ContentView.Clear()
				view.MetaInfoView.Clear()
			}
		},
	)
	view.Table.Select(0, 0)
}

