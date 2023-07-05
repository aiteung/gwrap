package docs

import (
	"strings"

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

func searchTextElement(pattern string, elements ...*docs.ParagraphElement) (res *docs.ParagraphElement) {
	if len(elements) < 1 {
		return
	}
	for _, v := range elements {

		if v.TextRun == nil {
			continue
		}

		curString1 := strings.ReplaceAll(v.TextRun.Content, "\n", "")
		curString2 := strings.ReplaceAll(curString1, "\t", "")
		curString3 := strings.ReplaceAll(curString2, " ", "")
		if strings.Contains(curString3, pattern) {
			res = v
			break
		}
	}

	return
}
