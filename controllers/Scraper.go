package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

var wikipediaUrl = "https://pt.wikipedia.org"

func RunScraper() {
	fmt.Println("Starting Wikipedia Scrapping ...")
	federatedUnitBrazilScrap()
}

func federatedUnitBrazilScrap() {

	ufLink := make(map[string]string)
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Panic("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("span.mw-page-title-main", func(e *colly.HTMLElement) {
		fmt.Println("Title " + e.Text)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isStateCategoryLink(link) {

			uf := strings.ReplaceAll(e.Text, "Naturais ", "")
			uf = strings.ReplaceAll(uf, "estado ", "")
			uf = strings.ReplaceAll(uf, "(estado)", "")
			uf = strings.ReplaceAll(uf, "(Brasil)", "")
			uf = strings.ReplaceAll(uf, "de ", "")
			uf = strings.ReplaceAll(uf, "do ", "")
			uf = strings.ReplaceAll(uf, "da ", "")

			ufLink[uf] = link
			//e.Request.Visit(wikipediaUrl + link)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit(wikipediaUrl + "/wiki/Categoria:Naturais_do_Brasil_por_unidade_federativa")
}

// Verify if the link is related to State Naturals born link or not
func isStateCategoryLink(link string) bool {
	return strings.Contains(link, "/wiki/Categoria:Naturais_") && !strings.Contains(link, "_Brasil")
}
