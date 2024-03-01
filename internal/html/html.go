package html

import (
	"strings"

	"golang.org/x/net/html"
)

func GetTextFromHtml(htmlStr string) string {
	plain := ""
	domDocTest := html.NewTokenizer(strings.NewReader(htmlStr))
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" {
				continue
			}
			TxtContent := html.UnescapeString(string(domDocTest.Text()))
			if len(TxtContent) > 0 {
				plain += TxtContent
			}
		}
	}
	return plain
}
