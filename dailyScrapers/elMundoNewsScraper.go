package dailyScrapers

import (
	uuid "github.com/nu7hatch/gouuid"
	"news-scrapers-workers-go/models"
	"news-scrapers-workers-go/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type ElMundoNewsScraper struct {
	Config models.ScrapingConfig
}

func (scraper *ElMundoNewsScraper) ScrapNewUrl(urlNew models.UrlNew) models.NewScraped {
	result := models.NewScraped{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.elmundo.es"),
	)
	headline := ""
	content := ""
	date := ""
	tags :=  []string{}

	c.OnHTML("p", func(e *colly.HTMLElement) {
		content = content + "\n" + e.Text
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		headline = e.Text
	})
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "/www.elmundo.es/t/") {
			tags = append(tags, e.Text)
		}
	})
	c.Visit(urlNew.Url)

	result.Url = urlNew.Url
	result.Headline = headline
	result.ScraperID = scraper.Config.ScraperId

	result.NewsPaper = "elmundo"
	result.Content = strings.TrimSpace(content)
	result.Content = cleanUpPublicityElMundo(result.Content)

	date = extractDateFromUrlElMundo(urlNew.Url)

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

func cleanUpPublicityElMundo(content string) string{
	result :=content
	breakingStrings := []string{"Conforme a los criterios deThe Trust Project", "Para seguir leyendo, hazte Premium", "El director de El Mundo selecciona las noticias de mayor interés para ti."}

	for _, str := range(breakingStrings){
		if utils.StringInSlice(str, breakingStrings){
			result = strings.Split(result, str)[0]
		}
	}

	return result
}

func extractDateFromUrlElMundo(url string) string {
	date := ""
	splittedUrl := strings.Split(url, "/")

	if strings.Contains(url, "/") && len(splittedUrl)>4 {
		if checkIfInt(splittedUrl[2]) {
			year := splittedUrl[2]
			month := splittedUrl[3]
			day := splittedUrl[4]
			date = year + "-" + month + "-" + day
			return date
		}
	}
	if strings.Contains(url, "/") && len(splittedUrl)>5 {
		if checkIfInt(splittedUrl[3]) {
			year := splittedUrl[3]
			month := splittedUrl[4]
			day := splittedUrl[5]
			date = year + "-" + month + "-" + day
			return date
		}
	}
	if strings.Contains(url, "/") && len(splittedUrl)>6 {
		if checkIfInt(splittedUrl[4]) {
			year := splittedUrl[4]
			month := splittedUrl[5]
			day := splittedUrl[6]
			date = year + "-" + month + "-" + day
		}
		return date
	}
	return date
}

func checkIfInt (str string) bool{
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}