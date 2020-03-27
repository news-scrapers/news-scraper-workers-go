package dailyScrapers

import (
	uuid "github.com/nu7hatch/gouuid"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type ElPaisEsNewsScraper struct {
	Config models.ScrapingConfig
}

func (scraper *ElPaisEsNewsScraper) ScrapNewUrl(urlNew models.UrlNew) models.NewScraped {
	result := models.NewScraped{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("elpais.com"),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})
	headline := ""
	content := ""
	date := ""
	tags :=  []string{}

	c.OnHTML("h1", func(elem *colly.HTMLElement) {
			headline = strings.TrimSpace(elem.Text)
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		content = content + "\n" + e.Text
	})
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if e.Attr("name")=="news_keywords" {
			content := e.Attr("content")
			tags = strings.Split(content, ",")
		}
	})
	c.Visit(urlNew.Url)


	date = strings.Split(urlNew.Url, "/")[4]

	result.Url = urlNew.Url
	result.Headline = headline
	result.ScraperID = scraper.Config.ScraperId

	result.NewsPaper = "elpais"
	result.Content = strings.TrimSpace(content)
	result.Content = cleanUpPublicityElPais(result.Content)

	t, _ := time.Parse("2006-01-02", date)
	result.Date = t
	result.DateString = date


	result.Tags = tags

	u, _ := uuid.NewV4()
	result.ID = u.String()

	log.Println("obtained new with headline " + headline)

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Wait()

	return result

}

func cleanUpPublicityElPais(content string) string{
	if (strings.Contains(content, "NEWSLETTER")){
		return strings.Split(content, "NEWSLETTER")[0]
	}
	return content
}
