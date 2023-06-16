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

func AddTable(column, rows int64) *docs.Request {
	temp := new(docs.Request)
	table := new(docs.InsertTableRequest)
	table.Columns = column
	table.Rows = rows
	return temp
}
