package historicScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	uuid "github.com/nu7hatch/gouuid"
)

type ElPaisScraper struct {
	Config models.ScrapingConfig
}

func (scraper ElPaisScraper) ScrapDate(date time.Time) []models.NewScraped {
	newsScraped := []models.NewScraped{}
	for i := 0; i < 12; i++ {
		newsScraped = append(newsScraped, scraper.ScrapDatePage(date, i)...)
	}
	return newsScraped
}

func (scraper ElPaisScraper) ScrapDatePage(date time.Time, page int) []models.NewScraped {
	newsScraped := []models.NewScraped{}
	urlNews := scraper.GetNewsUrls(date, page)
	for _, url := range urlNews {
		result := scraper.ScrapPage(url)
		newsScraped = append(newsScraped, result)
		// _ = result.Save(scraper.Config)
	}
	return newsScraped
}

func (scraper *ElPaisScraper) ScrapPage(urlNew UrlNew) models.NewScraped {
	result := models.NewScraped{}
	result.Tags = []string{}
	tags :=  []string{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("elpais.com"),
	)
	content := ""
	headline := ""
	// On every a element which has href attribute call callback
	c.OnHTML("p", func(e *colly.HTMLElement) {
		content = content + "\n" + e.Text
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		headline = e.Text
	})

	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if e.Attr("name")=="news_keywords" {
			tags = strings.Split(content, ",")
		}
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
	result.NewsPaper = "elpais"
	result.Url = urlNew.url
	result.Date = urlNew.date
	result.Headline = headline
	result.Content = cleanUpPublicity( headline + "\n" + content)
	result.Tags = tags

	u, _ := uuid.NewV4()
	result.ID = u.String()
	return result

}

func cleanUpPublicity (content string) string{
	if (strings.Contains(content, "NEWSLETTER")){
		return strings.Split(content, "NEWSLETTER")[0]
	}
	return content
}

func (scraper *ElPaisScraper) GetNewsUrls(date time.Time, page int) []UrlNew {
	formatedDate := date.Format("20060102")
	newsUrls := []UrlNew{}
	url := fmt.Sprintf("https://elpais.com/tag/fecha/%s/%v", formatedDate, page)
	//url := "https://elpais.com/tag/fecha/" + formatedDate + "/"

	c := colly.NewCollector(
		colly.AllowedDomains("elpais.com"),
	)

	c.OnHTML("article", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		if !strings.Contains(link, "/hemeroteca/"){
			url := e.Request.AbsoluteURL(link)

			msg := fmt.Sprintf("Link found: %q -> %s\n", e.Text, link)
			log.Info(msg)

			urlNew := UrlNew{url: url, date: date}

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
