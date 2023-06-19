package docs

import (
	"errors"
	"google.golang.org/api/docs/v1"
	"log"
	"strings"
)

type GoogleDocs struct {
	srv *docs.Service
}

func NewGoogleDocs(srvDocs *docs.Service) (res GoogleDocs) {
	res = GoogleDocs{srv: srvDocs}
	return
}

func (gdocs *GoogleDocs) Service() *docs.Service {
	return gdocs.srv
}

func (gdocs *GoogleDocs) FindAndReplace(docId string, request ...*docs.Request) (err error) {
	_, err = gdocs.srv.Documents.BatchUpdate(docId, &docs.BatchUpdateDocumentRequest{Requests: request}).Do()
	if err != nil {
		return
	}

	return
}

func (gdocs *GoogleDocs) FindTextLocation(docId string, pattern string) (loc docs.Location, err error) {
	file, err := gdocs.srv.Documents.Get(docId).Do()

	for _, element := range file.Body.Content {
		if element.Paragraph == nil {
			continue
		}
		for _, run := range element.Paragraph.Elements {
			if run.TextRun == nil {
				continue
			}
			if !strings.ContainsAny(run.TextRun.Content, pattern) {
				continue
			}
			loc.Index = run.StartIndex
			break
		}
	}

	return
}

func (gdocs *GoogleDocs) GetTableLocation(docID string) (data *docs.Location, err error) {

	doc, err := gdocs.srv.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve document: %v", err)
	}

	for _, content := range doc.Body.Content {
		if content.Table != nil {
			return &docs.Location{
				Index:           content.StartIndex,
				ForceSendFields: content.ForceSendFields,
				NullFields:      content.NullFields,
			}, nil
		}
	}

	err = errors.New("table not found")
	return
}
