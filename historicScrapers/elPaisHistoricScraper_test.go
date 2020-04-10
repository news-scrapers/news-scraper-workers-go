package historicScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestScraperElpais(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{UrlBase: "http://localhost:8000", NewsPaper: []string{"elpais"}, ScraperId: "testScraperElpais", DeviceID: "testDeviceElpais"}
	scraper := ElPaisScraper{Config: config}
	date := time.Date(2017, 10, 02, 0, 0, 0, 0, time.UTC)
	urlNewsElPais := scraper.GetNewsUrls(date, 1)

	resultElPais := scraper.ScrapPage(urlNewsElPais[1])
	resultElPais2 := scraper.ScrapPage(urlNewsElPais[2])
	resultElPais3 := scraper.ScrapPage(urlNewsElPais[5])


	assert.NotEqual(t, resultElPais.Headline, "", "OK response is expected")
	assert.NotEmpty(t, resultElPais.Content,"Should fill content")
	assert.NotEmpty(t, resultElPais.Url,"Should fill url")
	assert.NotEmpty(t, resultElPais.Tags,"Should fill tags")

	assert.NotEqual(t, resultElPais2.Headline, "", "OK response is expected")
	assert.NotEmpty(t, resultElPais2.Content,"Should fill content")
	assert.NotEmpty(t, resultElPais2.Url,"Should fill url")
	assert.NotEmpty(t, resultElPais2.Tags,"Should fill tags")

	//assert.NotEqual(t, resultElPais3.Headline, "", "OK response is expected")
	assert.NotEmpty(t, resultElPais3.Content,"Should fill content")
	assert.NotEmpty(t, resultElPais3.Url,"Should fill url")
	//assert.NotEmpty(t, resultElPais3.Tags,"Should fill tags")

}
