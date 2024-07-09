package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	c.OnRequest(func(r *colly.Request) {})
	c.OnError(func(_ *colly.Response, err error) {
		log.Panic("Something went wrong: ", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isStateCategoryLink(link) {
			uf := clearCityAndStateNames(e.Text)

			state := models.State{}
			database.DB.Where(&models.State{Name: uf}).First(&state)

			if state.ID == 0 {
				database.DB.Create(&models.State{Name: uf})
			}

			ufLinks[&state] = wikipediaUrl + link
		}
	})

	c.OnScraped(func(r *colly.Response) {
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

		fmt.Println("#### Scrappring State: " + uf.Name + " ####")

		c := colly.NewCollector()
		c.OnRequest(func(r *colly.Request) {})

		c.OnError(func(_ *colly.Response, err error) {
			log.Panic("Something went wrong: ", err)
		})

		c.OnHTML("#mw-subcategories", func(e *colly.HTMLElement) {

			e.ForEach(".mw-category-group", func(_ int, elem *colly.HTMLElement) {
				h3Title := elem.ChildText("h3") //From A to Z, however, some pages have ~ and empty strings that we don't care;
				elem.ForEach("a[href]", func(_ int, elem2 *colly.HTMLElement) {
					cityName := clearCityAndStateNames(elem2.Text)

					//This is because the São Paulo and Rio de Janeiro has cities with the same name
					isSpOrRJCityNatural := cityName == "Naturais da cidade de São Paulo" || cityName == "Naturais da cidade do Rio de Janeiro"

					if regexp.MustCompile(`[A-Z]`).MatchString(h3Title) || isSpOrRJCityNatural {
						link := elem2.Attr("href")

						fmt.Println(cityName + " >>> " + link)

						e.Request.Visit(wikipediaUrl + link)
					}
				})

			})

			alredyVisited := false
			e.DOM.Before(".mw-content-ltr").Find("a[href]").Each(func(i int, s *goquery.Selection) {
				link, _ := s.Attr("href")

				if !alredyVisited && s.Text() == "página seguinte" {
					alredyVisited = true
					e.Request.Visit(wikipediaUrl + link)
				}
			})

		})

		c.OnHTML("#mw-pages", func(e *colly.HTMLElement) {
			//fmt.Println("Page " + e.Text)
		})

		c.OnScraped(func(r *colly.Response) {})
		c.Visit(link)

	}

	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(diff)
}

func clearCityAndStateNames(text string) string {
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
