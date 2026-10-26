package extract

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Title    string
	Body     string
	Headings []string
	Links    []string
}

func Parse(reader io.Reader) (Page, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return Page{}, err
	}

	title := strings.TrimSpace(doc.Find("title").Text())

	content := doc.Find("main")
	if content.Length() == 0 {
		content = doc.Find("body")
	}

	content.Find("script, style, noscript, nav, footer").Remove()
	cleanContent := content.Text()
	cleanContent = clean(cleanContent)

	var headings []string
	doc.Find("h1, h2, h3").Each(func(_ int, s *goquery.Selection) {
		heading := clean(s.Text())
		if heading != "" {
			headings = append(headings, heading)
		}
	})

	var links []string
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if ok {
			links = append(links, link)
		}
	})

	return Page{
		Title:    title,
		Body:     cleanContent,
		Headings: headings,
		Links:    links,
	}, nil
}

func clean(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
