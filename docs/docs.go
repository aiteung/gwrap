package docs

import (
	"errors"
	"google.golang.org/api/docs/v1"
	"log"
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

	if len(file.Body.Content) < 1 {
		return
	}

	for _, element := range file.Body.Content {
		switch {
		case element.Paragraph != nil:
			if data := searchTextElement(pattern, element.Paragraph.Elements...); data != nil {
				loc.Index = (data.StartIndex + data.EndIndex) / 2
				return
			}
		case element.Table != nil:
			if len(element.Table.TableRows) < 1 {
				return
			}
			for _, v := range element.Table.TableRows {
				if len(v.TableCells) < 1 {
					continue
				}
				for _, v2 := range v.TableCells {
					if len(v2.Content) < 1 {
						continue
					}
					for _, v3 := range v2.Content {
						if v3.Paragraph == nil {
							continue
						}
						if data := searchTextElement(pattern, v3.Paragraph.Elements...); data != nil {
							loc.Index = data.StartIndex
							return
						}
					}
				}
			}
		case element.TableOfContents != nil:
			if len(element.TableOfContents.Content) < 1 {
				return
			}
			for _, v := range element.TableOfContents.Content {
				if len(v.Paragraph.Elements) < 1 {
					continue
				}
				if data := searchTextElement(pattern, v.Paragraph.Elements...); data != nil {
					loc.Index = data.StartIndex
				}
			}
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
				Index:           (content.StartIndex + content.EndIndex) / 2,
				ForceSendFields: content.ForceSendFields,
				NullFields:      content.NullFields,
			}, nil
		}
	}

	err = errors.New("table not found")
	return
}
