package dailyScrapers

import (
	uuid "github.com/nu7hatch/gouuid"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type PublicoNewsScraper struct {
	Config models.ScrapingConfig
}


func (scraper *PublicoNewsScraper) ScrapNewUrl(urlNew models.UrlNew) models.NewScraped {
	result := models.NewScraped{}

	ajaxUrl := urlNew.Url
	// Instantiate default collector
	c := colly.NewCollector(

	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
			headline := ""
			content := ""
			date := ""
			tags := []string{}

			e.ForEach("div", func(_ int, elem *colly.HTMLElement) {
				if elem.Attr("itemprop")=="headline name"{
					headline = strings.TrimSpace(elem.Text)
				}
			})

			e.ForEach(".article-body", func(_ int, elem *colly.HTMLElement) {
				elem.ForEach("p", func(_ int, elem2 *colly.HTMLElement) {
						content = content + " " + elem2.Text
				})
			})
			e.ForEach("span", func(_ int, elem *colly.HTMLElement) {
				if elem.Attr("class")=="published" && strings.Contains(elem.Text, " ") {
					date = strings.Split(elem.Text," ")[0]
				}
			})

			e.ForEach(".article-tags", func(_ int, elem *colly.HTMLElement) {
				elem.ForEach("a", func(_ int, elem2 *colly.HTMLElement) {
					tags = append(tags, elem2.Text)
				})
			})

			t, _ := time.Parse("02/01/2006", date)
			result.Url=urlNew.Url
			result.Headline=headline
			result.ScraperID = scraper.Config.ScraperId

			result.NewsPaper = "publico"
			result.Content = strings.TrimSpace(content)
			result.Date = t
			result.DateString = date
			u, _ := uuid.NewV4()
			result.ID = u.String()
			result.Tags = tags

			log.Println("obtained new with headline " + headline)


	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		url := "https://www.tripadvisor.es/" + e.Attr("href") //strings.Replace(e.Attr("href"), "/", "", 1)
		if e.Text == "Siguiente" {
			c.Visit(url)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting\n", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})

	c.Visit(ajaxUrl)
	c.Wait()

	return result

}
