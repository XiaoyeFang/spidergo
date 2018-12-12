package spider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"io"
)

func ParseAllVersion(resp io.Reader) (err error) {
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		glog.Errorf("NewDocumentFromReader err :%s", err)
	}
	doc.Find("#primary .listWidget").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".table-row").Text()
		glog.Infoln(title)
	})
	return nil
}
