package dailyScrapers

import (
	uuid "github.com/nu7hatch/gouuid"
	"news-scrapers-workers-go/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type VeinteMinutosNewsScraper struct {
	Config models.ScrapingConfig
}

func (scraper *VeinteMinutosNewsScraper) ScrapNewUrl(urlNew models.UrlNew) models.NewScraped {
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
		if strings.Contains(elem.Attr("class"), "article-title"){
			headline = strings.TrimSpace(elem.Text)
		}
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "paragraph"){
			content = content + "\n" + e.Text
		}
	})
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if e.Attr("property")=="article:tag" {
			content := e.Attr("content")
			tags = append(tags, content)
		}
	})
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if e.Attr("class") == "article-date"{
			e.ForEach("a", func(_ int, elem2 *colly.HTMLElement) {
				date = strings.Split(elem2.Text," ")[0]
			})
		}

	})

	c.Visit(urlNew.Url)


	result.Url = urlNew.Url
	result.Headline = headline
	result.ScraperID = scraper.Config.ScraperId

	result.NewsPaper = "20minutos"
	result.Content = strings.TrimSpace(content)
	result.Content = cleanUpPublicityVeinteMinutos(result.Content)

	t, _ := time.Parse("02.01.2006", date)
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

func cleanUpPublicityVeinteMinutos(content string) string{
	if (strings.Contains(content, "NEWSLETTER")){
		return strings.Split(content, "NEWSLETTER")[0]
	}
	return content
}
