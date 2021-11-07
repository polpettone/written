package ui

import (
	"fmt"
	"github.com/polpettone/written/cmd/models"
	"strings"
	"time"
)

const outputPattern = `
%s
%s
%s
`

func documentMetaView(document models.Document) string {
		lastModified := document.Info.ModTime().Format(time.RFC822)
		name := document.Info.Name()
		tags := strings.Join(document.Tags, SPACE)
		return fmt.Sprintf(outputPattern, name, lastModified, tags)
}

