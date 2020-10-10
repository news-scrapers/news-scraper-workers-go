package dailyScrapers

import (
	uuid "github.com/nu7hatch/gouuid"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type ElDiaroEsNewsScraper struct {
	Config models.ScrapingConfig
}


func (scraper *ElDiaroEsNewsScraper) ScrapNewUrl(urlNew models.UrlNew) models.NewScraped {
	result := models.NewScraped{}

	ajaxUrl := urlNew.Url
	// Instantiate default collector
	c := colly.NewCollector(

	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})
	date := ""


	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if (e.Attr("property") == "article:published_time") {
			date = e.Attr("content")
		}
	})

	c.OnHTML(".article-page", func(e *colly.HTMLElement) {
			headline := ""
			content := ""
			tags := []string{}


		e.ForEach(".title", func(_ int, elem *colly.HTMLElement) {
				headline = strings.TrimSpace(elem.Text)
			})

			e.ForEach("p.article-text", func(_ int, elem *colly.HTMLElement) {
				content = content + " " + elem.Text
			})
			e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
				elem.ForEach("a.tag-link", func(_ int, elem2 *colly.HTMLElement) {
					tags = append(tags, strings.TrimSpace(elem2.Text))
				})
			})

			layout := "2006-01-02"
			date = strings.Split(date, "T")[0]
			t, _ := time.Parse(layout, date)
			result.Url=urlNew.Url
			result.Headline=headline
			result.ScraperID = scraper.Config.ScraperId

			result.NewsPaper = "eldiario.es"
			result.Content = strings.TrimSpace(content)
			result.Date = t
			result.DateString = date
			result.Tags = tags
			u, _ := uuid.NewV4()
			result.ID = u.String()

			log.Println("obtained new with headline " + headline)


	})


	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Visit(ajaxUrl)
	c.Wait()

	return result

}
