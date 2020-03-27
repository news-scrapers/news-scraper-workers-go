package historicScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestScraperlavanguadia(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{UrlBase: "http://localhost:8000", NewsPaper: []string{"lavanguardia"}, ScraperId: "testScraperLavanguardia", DeviceID: "testDevicelavanguardia"}

	scraper := LaVanguardiaScraper{Config: config}
	date := time.Date(2018, 05, 04, 0, 0, 0, 0, time.UTC)
	urlNews := scraper.GetNewsUrls(date, 1)

	result := scraper.ScrapPage(urlNews[1])
	assert.NotEqual(t, result.Content, "", "OK response is expected")
}
