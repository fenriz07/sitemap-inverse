package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	link "github.com/fenriz07/link/students/fenriz"
	"github.com/fenriz07/sitemap-inverse/urlset"
	"github.com/fenriz07/sitemap/helpers"
)

func main() {

	start := time.Now()

	urlFlag := flag.String("url", "none", "Domain to scan Ex. url=https://google.com")
	depthFlag := flag.Int("depth", 2, "Specify the scanner depth number. By default is 2")

	flag.Parse()

	if *urlFlag == "none" {
		fmt.Println("I'm sorry you must specify a domain,use the flag url.")
		os.Exit(2)
	}

	//Obtenemos las paginas
	pages := bfs(*urlFlag, *depthFlag)

	externalurls := make(map[string][]string)

	for _, page := range pages {
		v := get(page, true)

		externalurls[page] = v
	}

	urlset.PrintXML(externalurls)

	elapsed := time.Since(start)

	fmt.Printf("Tiempo en ejecución: %v \n", elapsed)

}

//Algoritmo BFS
func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})

	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i < maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		linkChanel := make(chan []string)
		lenchanel := 0

		for url, _ := range q {
			/* con la linea # podemos comprobar si la llave existe en el mapa.
			Si dicha llave existe quiere decir que la url fue analizada */

			if _, ok := seen[url]; ok {
				continue
			}

			lenchanel++

			/* Se le asigna el link que se va a analizar, para que no pueda ser analizado
			en un futuro */
			seen[url] = struct{}{}

			go func() {
				links := get(url, false)

				linkChanel <- links
			}()

			/*Se prepara nq con los valores obtenidos que posteriormente se analizaran*/
			/*for _, link := range links {
				nq[link] = struct{}{}
			}*/
		}

		for i := 0; i < lenchanel; i++ {
			links := <-linkChanel

			for _, link := range links {
				nq[link] = struct{}{}
			}
		}

		close(linkChanel)
	}

	//Se obtiene el resultado

	ret := make([]string, 0, len(seen))

	for url, _ := range seen {
		ret = append(ret, url)
	}

	return ret
}

func get(urlStr string, inverse bool) []string {

	resp, err := http.Get(urlStr)

	if err != nil {
		helpers.DD(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	base := baseUrl.String()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		helpers.DD(err)
	}

	allLinks := link.ParseHtml(string(body))

	pages := filter(base, createPages(*allLinks, base), inverse)

	return pages

}

func createPages(links []link.Link, base string) []string {

	var href string
	var allLinks []string

	for _, l := range links {

		href = l.Href

		switch {
		case strings.HasPrefix(href, "/"):
			allLinks = append(allLinks, base+href)
		case strings.HasPrefix(href, "http"):
			allLinks = append(allLinks, href)
		}
	}

	return allLinks
}

func filter(base string, links []string, inverse bool) []string {
	var ret []string

	for _, link := range links {

		if inverse == false {
			if strings.HasPrefix(link, base) {
				ret = append(ret, link)
			}
		} else {
			if !strings.HasPrefix(link, base) && !strings.HasPrefix(link, "/") {
				ret = append(ret, link)
			}
		}

	}

	if inverse == true {
		return ret
	}

	return unique(ret)
}

func unique(elements []string) []string {

	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
