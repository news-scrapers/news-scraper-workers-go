package dailyScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type LaVanguardiaFullIndexManager struct {
	Config             models.ScrapingConfig
	Index              models.ScrapingIndex
	MaxPagesVanguardia int

}

var baseUrlLaVanguardia = "https://www.lavanguardia.com"

func (scraper LaVanguardiaFullIndexManager) ScrapNewsInItems(baseUrl string, scrapingIndex *models.ScrapingIndex) {
	if scraper.MaxPagesVanguardia == 0 {
		scraper.MaxPagesVanguardia = 100
	}
	urlsPending := []models.UrlNew{}
	newsScraper := LaVanguardiaNewsScraper{scraper.Config}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("https://elpais.com/"),
	)
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, elem *colly.HTMLElement) {
			if strings.Contains(elem.Attr("class"), "story-header-title-link") {
				url := elem.Attr("href")
				date := time.Now()
				urlScrap := models.UrlNew{url,date}
				fmt.Println(url)
				urlsPending = append(urlsPending, urlScrap)
			}
		})

	})


	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting  ", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Visit(baseUrl)

	log.Println("Collected pages")
	log.Println(urlsPending)

	for scrapingIndex.PageIndex < len(urlsPending)-3 {
		log.Println("---------------")
		log.Printf("Scraping page %d", scrapingIndex.PageIndex)
		log.Println("---------------")

		out1 := make(chan models.NewScraped)
		out2 := make(chan models.NewScraped)
		out3 := make(chan models.NewScraped)

		url1 := urlsPending[scrapingIndex.PageIndex]
		url2 := urlsPending[scrapingIndex.PageIndex+1]
		url3 := urlsPending[scrapingIndex.PageIndex+2]

		go scraper.scrapAllReviewsInUrl(url1, &newsScraper, out1)
		go scraper.scrapAllReviewsInUrl(url2, &newsScraper, out2)
		go scraper.scrapAllReviewsInUrl(url3, &newsScraper, out3)

		resultsInPage1, resultsInPage2, resultsInPage3 := <-out1, <-out2, <-out3
		resultsInPage1.SaveOrUpdate()
		resultsInPage2.SaveOrUpdate()
		resultsInPage3.SaveOrUpdate()

		//results = append(results, resultsInPage1...)
		//results = append(results, resultsInPage2...)

		scrapingIndex.PageIndex = scrapingIndex.PageIndex + 3
		scrapingIndex.Save()

	}
	scrapingIndex.PageIndex = 0
	scrapingIndex.Save()

	//return results

}

func (scraper LaVanguardiaFullIndexManager) scrapAllReviewsInUrl(urlbase models.UrlNew, reviewsScraper *LaVanguardiaNewsScraper, out chan models.NewScraped) models.NewScraped {
	result := reviewsScraper.ScrapNewUrl(urlbase)

	out <- result
	return result
}
