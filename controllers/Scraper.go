package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/pedro-hos/guess-who-web/database"
	"github.com/pedro-hos/guess-who-web/models"
)

var wikipediaUrl = "https://pt.wikipedia.org"
var bornInBrazilByUF = wikipediaUrl + "/wiki/Categoria:Naturais_do_Brasil_por_unidade_federativa"

func RunScraper() {
	fmt.Println("Starting Wikipedia Scrapping ...")
	federatedUnitBrazilScrap()
}

func federatedUnitBrazilScrap() {

	ufLinks := make(map[*models.State]string)
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Panic("Something went wrong: ", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isStateCategoryLink(link) {
			uf := clearStateName(e.Text)

			state := models.State{}
			database.DB.Where(&models.State{Name: uf}).First(&state)

			if state.ID == 0 {
				database.DB.Create(&models.State{Name: uf})
			}

			ufLinks[&state] = wikipediaUrl + link
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
		citiesScrap(ufLinks)
	})

	c.Visit(bornInBrazilByUF)
}

func citiesScrap(states map[*models.State]string) {
	t1 := time.Now()
	for k, v := range states {

		uf := k
		link := v

		//TODO: Need to remove, I just added this to skip all cities and states
		if !strings.Contains(uf.Name, "São Paulo") {
			continue
		}

		fmt.Println(uf.Name)
		c := colly.NewCollector()
		c.OnRequest(func(r *colly.Request) {})

		c.OnError(func(_ *colly.Response, err error) {
			log.Panic("Something went wrong: ", err)
		})

		c.OnHTML("#mw-subcategories", func(e *colly.HTMLElement) {

			e.ForEach(".mw-category-group", func(_ int, elem *colly.HTMLElement) {
				h3Title := elem.ChildText("h3") //From A to Z, however, some pages have ~ and empty strings that we don't care;

				if regexp.MustCompile(`[A-Z]`).MatchString(h3Title) {
					elem.ForEach("a[href]", func(_ int, elem2 *colly.HTMLElement) {
						city := elem2.Text
						link := elem2.Attr("href")
						fmt.Println(city + " >>> " + link)
					})
				} else {
					fmt.Println("Empty or ~")
				}
			})
		})

		c.OnScraped(func(r *colly.Response) {
			//fmt.Println(r.Request.URL, " scraped!")
		})

		c.Visit(link)

	}

	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(diff)
}

func clearStateName(text string) string {
	uf := strings.ReplaceAll(text, "Naturais do estado de ", "")
	uf = strings.ReplaceAll(uf, "Naturais do estado do ", "")
	uf = strings.ReplaceAll(uf, "Naturais do ", "")
	uf = strings.ReplaceAll(uf, "Naturais de ", "")
	uf = strings.ReplaceAll(uf, "Naturais da ", "")
	uf = strings.ReplaceAll(uf, "(estado)", "")
	uf = strings.ReplaceAll(uf, "(Brasil)", "")
	return uf
}

// Verify if the link is related to State Naturals born link or not
func isStateCategoryLink(link string) bool {
	return strings.Contains(link, "/wiki/Categoria:Naturais_") && !strings.Contains(link, "_Brasil")
}
