package pipes

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/KonishchevDmitry/rsspipes/util"
)

var log = util.MustGetLogger("pipes")

type UrlBuilder struct {
	baseUrl string
}

func (b *UrlBuilder) getUrl(url string) string {
	if strings.HasPrefix(url, "/") {
		url = b.baseUrl + url
	}
	return url
}

func getDescriptionFromSelection(selection *goquery.Selection, urlBuilder UrlBuilder) (string, error) {
	selection.Find("script").Remove()

	selection.Find("a").Each(func(i int, link *goquery.Selection) {
		if url, exists := link.Attr("href"); exists {
			link.SetAttr("href", urlBuilder.getUrl(url))
		}
	})

	selection.Find("img").Each(func(i int, image *goquery.Selection) {
		if url, exists := image.Attr("src"); exists {
			image.SetAttr("src", urlBuilder.getUrl(url))
		}
	})

	return selection.Html()
}

func getSelectionHtml(selection *goquery.Selection) string {
	html, err := selection.Html()
	if err != nil {
		html = fmt.Sprintf("[Failed to render the HTML: %s]", err)
	}

	return html
}
