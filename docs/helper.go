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

func InsertImage(link string, location *docs.Location, size float64) *docs.Request {
	heightweight := docs.Dimension{
		Magnitude:       size,
		Unit:            "PT",
		ForceSendFields: nil,
		NullFields:      nil,
	}

	if size == 0 {
		heightweight.Magnitude = 128
	}

	sizeD := &docs.Size{
		Height: &heightweight,
		Width:  &heightweight,
	}

	return &docs.Request{
		InsertInlineImage: &docs.InsertInlineImageRequest{Location: location, Uri: link, ObjectSize: sizeD},
	}
}
