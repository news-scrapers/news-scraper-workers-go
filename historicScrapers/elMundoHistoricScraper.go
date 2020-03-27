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

type ElMundoHistoricScraper struct {
	Config models.ScrapingConfig
}

type UrlNew struct {
	url  string
	date time.Time
}

func (scraper ElMundoHistoricScraper) ScrapDate(date time.Time) []models.NewScraped {
	newsScraped := []models.NewScraped{}
	urlNews := scraper.GetNewsUrls(date)
	for _, url := range urlNews {
		result := scraper.ScrapPage(url)
		newsScraped = append(newsScraped, result)
		// _ = result.Save(scraper.Config)
	}
	return newsScraped
}

func (scraper *ElMundoHistoricScraper) ScrapPage(urlNew UrlNew) models.NewScraped {
	result := models.NewScraped{}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.elmundo.es"),
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

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Info("Something went wrong:", err)
	})
	c.OnHTML("a", func(e *colly.HTMLElement) {
				if strings.Contains(e.Attr("href"), "/www.elmundo.es/t/") {
					result.Tags = append(result.Tags, e.Text)
				}
	})

	c.Visit(urlNew.url)
	result.ScraperID = scraper.Config.ScraperId
	result.NewsPaper = "elmundo"
	result.Url = urlNew.url
	result.Date = urlNew.date
	result.Headline = headline
	result.Content = headline + "\n" + content

	u, _ := uuid.NewV4()
	result.ID = u.String()
	return result

}

func (scraper *ElMundoHistoricScraper) GetNewsUrls(date time.Time) []UrlNew {
	formatedDate := date.Format("2006/01/02")
	newsUrls := []UrlNew{}
	url := "https://www.elmundo.es/elmundo/hemeroteca/" + formatedDate + "/m/espana.html"
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.elmundo.es"),
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

func (scraper *ElMundoHistoricScraper) IsANew(url string) bool {
	blacklist := []string{"/vivienda.htm", "?intcmp=MENUDES22301", "/expansion-empleo.html", "/codigo-etico.html", "/avisolegal.html", "/contacto.html", "/menu.html", "/archivo.html", "/espana.html", "/noticias-mas-leidas.html", "/index.html"}
	notInBlacklist := true
	for _, element := range blacklist {
		containsElement := strings.Contains(url, element)
		if containsElement {
			notInBlacklist = false
			break
		}
	}
	return strings.Contains(url, ".html") && notInBlacklist
}

func removeDuplicatesFromSlice(s []UrlNew) []UrlNew {
	m := make(map[UrlNew]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
			//fmt.Println(item, "is a duplicate")
		} else {
			m[item] = true
		}
	}

	var result []UrlNew
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}
