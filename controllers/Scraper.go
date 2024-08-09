package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"github.com/pedro-hos/guess-who-web/database"
	"github.com/pedro-hos/guess-who-web/models"
)

var wikipediaUrl = "https://pt.wikipedia.org"
var bornInBrazilByUF = fmt.Sprintf("%s/wiki/Categoria:Naturais_do_Brasil_por_unidade_federativa", wikipediaUrl)

func readData(url string, rvsection string) models.Response {
	name := strings.ReplaceAll(url, wikipediaUrl+"/wiki/", "")
	wikiEndpoint := fmt.Sprintf("https://pt.wikipedia.org/w/api.php?action=query&prop=revisions&titles=%s&rvslots=*&rvprop=content&formatversion=2&format=json&rvsection=%s", name, rvsection)

	response, err := http.Get(wikiEndpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject models.Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getInfoFromWikiApi(wg *sync.WaitGroup, url string) {
	defer wg.Done()

	responseObject := readData(url, "1")
	content := responseObject.Query.Pages[0].Revisions[0].Slots.Main.Content

	if content == "" {
		responseObject := readData(url, "0")
		content = responseObject.Query.Pages[0].Revisions[0].Slots.Main.Content
	}

	database.DB.Model(&models.Card{}).Where("wikipedia_url = ?", url).Update("content_page", content)
}

func RunScraper(stateName string, cityName string) {
	defer fmt.Println("Scrap Finished!")
	fmt.Println("Starting Wikipedia Scrap ...")

	ufLinks := make(map[string]string)
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	allStates := stateName == ""

	c.OnRequest(func(r *colly.Request) {})
	c.OnError(func(_ *colly.Response, err error) {
		log.Panic("Something went wrong: ", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if isStateCategoryLink(link) {
			uf := clearCityAndStateNames(e.Text)
			state := models.State{}

			if allStates {
				database.DB.FirstOrCreate(&state, models.State{Name: uf})
				ufLinks[uf] = wikipediaUrl + link
			} else {
				if strings.EqualFold(stateName, uf) {
					database.DB.FirstOrCreate(&state, models.State{Name: uf})
					ufLinks[uf] = wikipediaUrl + link
				}
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {

		if len(ufLinks) == 0 {
			log.Fatalf("Can't find the state %s", stateName)
		} else {
			citiesScrap(ufLinks, cityName)
		}
	})

	c.Visit(bornInBrazilByUF)
}

func citiesScrap(states map[string]string, cityName string) {
	var wg sync.WaitGroup

	for k, v := range states {

		uf := k
		link := v
		allCities := cityName == ""
		var cityFound bool = false

		state := models.State{}
		database.DB.Where(&models.State{Name: uf}).First(&state)

		fmt.Println("#### Scrappring State [ " + strconv.FormatUint(uint64(state.ID), 10) + " ] " + state.Name)

		if state.ID == 0 {
			fmt.Println("We can't find the State")
			continue
		}

		c := colly.NewCollector()
		c.OnRequest(func(r *colly.Request) {})

		c.OnError(func(_ *colly.Response, err error) {
			log.Panic("Something went wrong: ", err)
		})

		c.OnHTML("#mw-subcategories", func(e *colly.HTMLElement) {

			e.ForEach(".mw-category-group", func(_ int, elem *colly.HTMLElement) {
				h3Title := elem.ChildText("h3") //From A to Z, however, some pages have ~ and empty strings that we don't care;
				elem.ForEach("a[href]", func(_ int, elem2 *colly.HTMLElement) {
					cityNameNormalized := clearCityAndStateNames(elem2.Text)

					//This is because the São Paulo and Rio de Janeiro has cities with the same name
					isSpOrRJCityNatural := cityNameNormalized == "Naturais da cidade de São Paulo" || cityNameNormalized == "Naturais da cidade do Rio de Janeiro"

					if regexp.MustCompile(`[A-Z]`).MatchString(h3Title) || isSpOrRJCityNatural {
						link := elem2.Attr("href")
						city := models.City{}

						if allCities {
							database.DB.FirstOrCreate(&city, models.City{Name: cityNameNormalized, StateId: state.ID})
							e.Request.Visit(wikipediaUrl + link)
						} else {
							if strings.EqualFold(cityName, cityNameNormalized) {
								database.DB.FirstOrCreate(&city, models.City{Name: cityNameNormalized, StateId: state.ID})
								cityFound = true
								e.Request.Visit(wikipediaUrl + link)
							}
						}
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

		c.OnHTML("main#content", func(e *colly.HTMLElement) {

			DOM := e.DOM
			title := DOM.Find(".mw-page-title-main").Text()
			isCity := strings.Contains(title, "Naturais")

			if isCity {
				cityNameNormalized := clearCityAndStateNames(title)
				city := models.City{}
				database.DB.Where(&models.City{Name: cityNameNormalized, StateId: state.ID}).First(&city)

				if city.ID == 0 {
					//some pages have people there are not related to the city, but to the state. So far, I don't want to save them. Need to think if I will save on State Capital, or just don't save
					fmt.Printf("Cant find the city %s, on the database", cityNameNormalized)
				} else {
					categoryGroup := DOM.Find(".mw-category-generated").Find("div.mw-category-group")
					treeSection := categoryGroup.Find(".CategoryTreeSection")

					categoryGroup.Each(func(i int, s *goquery.Selection) {
						if s.HasNodes(treeSection.Nodes...).Length() == 0 {
							s.Find("a[href]").Each(func(count int, s2 *goquery.Selection) {
								link, _ := s2.Attr("href")
								name := s2.Text()
								pageUrl := wikipediaUrl + link

								database.DB.Create(&models.Card{Answer: name, CityId: city.ID, WikipediaURL: pageUrl})

								wg.Add(1)
								go getInfoFromWikiApi(&wg, pageUrl)
								e.Request.Visit(link)
							})
						}
					})

				}
			} else {
				img, hasImg := DOM.Find("table.infobox").Find("a[href].mw-file-description > img.mw-file-element").Attr("src")

				if hasImg {
					imageURL := fmt.Sprintf("https:%s", img)
					database.DB.Model(&models.Card{}).Where("answer = ?", title).Update("ImageURL", imageURL)
				}
			}
		})

		c.OnScraped(func(r *colly.Response) {
			if !cityFound {
				log.Fatalf("Can't find the city %s on the wikipedia page. Please, review", cityName)
			}
		})

		c.Visit(link)

	}

	wg.Wait()
}

func clearCityAndStateNames(text string) string {
	phrases := []string{
		"Naturais do estado de ",
		"Naturais do estado do ",
		"Naturais do ",
		"Naturais de ",
		"Naturais da ",
		"(estado)",
		"(Brasil)",
	}

	uf := text
	for _, phrase := range phrases {
		uf = strings.ReplaceAll(uf, phrase, "")
	}

	return uf
}

// Verify if the link is related to State Naturals born link or not
func isStateCategoryLink(link string) bool {
	return strings.Contains(link, "/wiki/Categoria:Naturais_") && !strings.Contains(link, "_Brasil")
}
