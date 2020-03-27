package dailyScrapers
import (
	"news-scrapers-workers-go/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type VeinteMinutosFullIndexManager struct {
	Config models.ScrapingConfig
	Index  models.ScrapingIndex
	MaxPages int

}

var baseUrl20minutos = "https://www.20minutos.es"

func (scraper VeinteMinutosFullIndexManager) ScrapNewsInItems(baseUrl string, scrapingIndex *models.ScrapingIndex) {
	if scraper.MaxPages == 0 {
		scraper.MaxPages = 100
	}
	urlsPending := []models.UrlNew{}
	newsScraper := VeinteMinutosNewsScraper{scraper.Config}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("https://elpais.com/"),
	)
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, elem *colly.HTMLElement) {
			if strings.Contains(elem.Attr("href"), baseUrl20minutos + "/noticia/") {
				url := elem.Attr("href")
				date := time.Now()
				urlScrap := models.UrlNew{url,date}
				fmt.Println(url)
				urlsPending = append(urlsPending, urlScrap)
			}
		})

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("title"), "Último")  {
			currentPage := getCurrentPage20Minutos(baseUrl)
			currentPageStr := strconv.Itoa(currentPage)

			nextPage := currentPage +1
			nextPageStr :=  strconv.Itoa(nextPage)
			nextPageUrl := strings.ReplaceAll(baseUrl, "/"+currentPageStr+ "/", "/") +nextPageStr+ "/"
			baseUrl = nextPageUrl
			c.Visit(nextPageUrl)


		}
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

func getCurrentPage20Minutos(url string) int {
	splitted := strings.Split(url, "/")
	page := splitted[len(splitted) -2]
	pageInt, _ := strconv.Atoi(page)
	return pageInt
}

func (scraper VeinteMinutosFullIndexManager) scrapAllReviewsInUrl(urlbase models.UrlNew, reviewsScraper *VeinteMinutosNewsScraper, out chan models.NewScraped) models.NewScraped {
	result := reviewsScraper.ScrapNewUrl(urlbase)

	out <- result
	return result
}
