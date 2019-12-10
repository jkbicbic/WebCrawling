package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const outputDir = "./pages/"

func main() {
	DownloadPages()
	CrawlPages()
}

// DownloadPages ...
func DownloadPages() {
	c := colly.NewCollector()

	c.OnHTML("a.tNpZ-r8HSFPRZ6NJvAkbQ", func(e *colly.HTMLElement) {
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	log.Println("saving output to", outputDir)
	os.MkdirAll(outputDir, os.ModePerm)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		err := r.Save(outputDir + strings.Replace(r.FileName(), "unknown", "html", 1))
		if err != nil {
			log.Print(err)
		}
	})

	c.Visit("https://www.jobstreet.com.ph/en/companies")
}

// CrawlPages ...
func CrawlPages() {
	t := &http.Transport{}
	companyJSON := []string{}

	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := colly.NewCollector()
	c.WithTransport(t)
	c.IgnoreRobotsTxt = true

	files, err := ioutil.ReadDir("./pages")
	if err != nil {
		log.Fatal(err)
	}

	// On every a element which has href attribute call callback
	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		log.Println(e.Text)
		companyJSON = append(companyJSON, strings.TrimSpace(e.Text))
	})

	for _, file := range files {
		fmt.Println(fmt.Sprintf("file:///Users/johnkennedybicbic/go/src/github.com/personal/WebCrawling/pages/%s", file.Name()))
		c.Visit(fmt.Sprintf("file:///Users/johnkennedybicbic/go/src/github.com/personal/WebCrawling/pages/%s", file.Name()))
	}

	// log.Println(companyJSON)

	WriteInFile(companyJSON)

}

// WriteInFile ...
func WriteInFile(c []string) {

	marshalledData, err := json.MarshalIndent(c, "", "")
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile("companies.json", marshalledData, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

}
