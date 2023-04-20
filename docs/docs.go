package docs

import (
	"google.golang.org/api/docs/v1"
)

func ReplaceTextDocs(char, text string) *docs.Request {
	tempData := new(docs.Request)
	tempData.ReplaceAllText = &docs.ReplaceAllTextRequest{
		ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: char},
		ReplaceText:  text,
	}

	return tempData
}

func FindAndReplaceOne(srvDocs *docs.Service, docId string, request ...*docs.Request) (err error) {
	_, err = srvDocs.Documents.BatchUpdate(docId, &docs.BatchUpdateDocumentRequest{Requests: request}).Do()
	if err != nil {
		return
	}

	return
}
