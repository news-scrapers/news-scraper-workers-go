package historicScrapers

import (
	"fmt"
	"strings"
	"time"

	"news-scrapers-workers-go/models"

	"github.com/gocolly/colly"
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type LaVanguardiaScraper struct {
	Config models.ScrapingConfig
}

func (scraper LaVanguardiaScraper) ScrapDate(date time.Time) []models.NewScraped {
	newsScraped := []models.NewScraped{}
	urlNews := []UrlNew{}
	for i := 1; i <= 6; i++ {
		urlNewsPage := scraper.GetNewsUrls(date, i)
		urlNews = append(urlNews, urlNewsPage...)
	}

	for _, url := range urlNews {
		result := scraper.ScrapPage(url)
		newsScraped = append(newsScraped, result)
		// _ = result.Save(scraper.Config)
	}
	return newsScraped
}

func (scraper *LaVanguardiaScraper) ScrapPage(urlNew UrlNew) models.NewScraped {
	result := models.NewScraped{}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("hemeroteca.lavanguardia.com"),
	)
	content := ""
	headline := ""
	// On every a element which has href attribute call callback
	c.OnHTML(".text", func(e *colly.HTMLElement) {
		content = content + "\n" + e.Text
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Visit(urlNew.url)
	result.ScraperID = scraper.Config.ScraperId
	result.NewsPaper = "lavanguardia"
	result.Url = urlNew.url
	result.Date = urlNew.date
	result.Headline = headline
	result.Content = headline + "\n" + content

	u, _ := uuid.NewV4()
	result.ID = u.String()
	return result

}

func (scraper *LaVanguardiaScraper) GetNewsUrls(date time.Time, page int) []UrlNew {
	year, month, day := date.Date()

	newsUrls := []UrlNew{}
	url := fmt.Sprintf("http://hemeroteca.lavanguardia.com/edition.html?bd=%d&bm=%d&by=%d&page=%d", day, month, year, page)
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("hemeroteca.lavanguardia.com"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		msg := fmt.Sprintf("Link found: %q -> %s\n", e.Text, link)
		log.Info(msg)
		url := e.Request.AbsoluteURL(link)
		urlNew := UrlNew{url: url, date: date}
		if scraper.IsANew(url) {
			newsUrls = append(newsUrls, urlNew)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Visit(url)
	urlsNoRepetition := removeDuplicatesFromSlice(newsUrls)

	return urlsNoRepetition
}

func (scraper *LaVanguardiaScraper) IsANew(url string) bool {
	urlAllowed := "hemeroteca.lavanguardia.com/preview/"
	return strings.Contains(url, urlAllowed)

}
