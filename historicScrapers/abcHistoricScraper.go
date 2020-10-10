package historicScrapers

import (
	"strconv"
	"strings"
	"time"

	"news-scrapers-workers-go/models"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	uuid "github.com/nu7hatch/gouuid"
)

type AbcHistoricScraper struct {
	Config models.ScrapingConfig
}

func (scraper AbcHistoricScraper) ScrapDate(date time.Time) []models.NewScraped {
	newsScraped := []models.NewScraped{}
	for i := 1; i < 80; i++ {
		time.Sleep(3 * time.Second)
		newsScraped = append(newsScraped, scraper.ScrapDatePage(date, i))
	}
	return newsScraped
}

// http://hemeroteca.abc.es/nav/Navigate.exe/hemeroteca/madrid/abc/2019/04/10/003.html
func (scraper AbcHistoricScraper) ScrapDatePage(date time.Time, page int) models.NewScraped {
	formatedDate := date.Format("20060102")
	formatedPage := strconv.Itoa(page) //pad.Left(strconv.Itoa(page), 3, "0")
	//https://www.abc.es/archivo/periodicos/abc-madrid-20201003-1.html
	url := "https://www.abc.es/archivo/periodicos/abc-madrid-" + formatedDate + "-" + formatedPage + ".html"
	result := models.NewScraped{}
	// Instantiate default collector

	c := colly.NewCollector(
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})
	content := ""

	// On every a element which has href attribute call callback
	c.OnHTML("p", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "ABC") && !strings.Contains(e.Text, "{") {
			content = content + "\n" + e.Text
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

	result.ScraperID = scraper.Config.ScraperId
	result.FullPage = true
	result.NewsPaper = "abc"
	result.Url = url
	result.Date = date
	result.Content = content

	u, _ := uuid.NewV4()
	result.ID = u.String()
	return result
}
