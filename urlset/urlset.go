package urlset

import (
	"bufio"
	"encoding/xml"
	"os"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func PrintXML(pages []string) {

	urls := make([]Url, 0, len(pages))

	for _, url := range pages {

		urls = append(urls, Url{Loc: url})
	}

	urlset := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	out, _ := xml.MarshalIndent(urlset, " ", "  ")

	createFileXML(out)
}

func createFileXML(outputXml []byte) {

	fo, err := os.Create("sitemap.xml")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(fo)

	if _, err := w.Write(outputXml); err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if err = w.Flush(); err != nil {
		panic(err)
	}
}
