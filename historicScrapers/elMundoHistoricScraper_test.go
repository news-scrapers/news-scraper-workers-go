package historicScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestScraper(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{}
	config.CreateFromJson()
	scraper := ElMundoHistoricScraper{Config: config}
	date := time.Date(2018, 05, 04, 0, 0, 0, 0, time.UTC)
	urlNews := scraper.GetNewsUrls(date)

	result := scraper.ScrapPage(urlNews[1])
	fmt.Println("----")

	fmt.Println(result.Tags)

	result2 := scraper.ScrapPage(urlNews[2])
	fmt.Println("----")

	fmt.Println(result2.Tags)

	result3 := scraper.ScrapPage(UrlNew{"https://www.elmundo.es/comunidad-valenciana/castellon/2019/12/19/5dfa7cc6fdddff6c6e8b45e9.html", date})
	fmt.Println("----")

	fmt.Println(result3.Tags)
	assert.NotEqual(t, result.Headline, "", "OK response is expected")
	assert.NotEmpty(t, result.Content,"Should fill content")
	assert.NotEmpty(t, result.Url,"Should fill url")

	assert.NotEmpty(t, result3.Tags,"Should fill tags")
}
