package dailyScrapers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestScraperLaVanguardia(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperdiario", DeviceID: "testDevicediario"}
	scraper := LaVanguardiaNewsScraper{Config: config}
	baseUrl := "https://www.lavanguardia.com/internacional/20200327/48113436277/coronavirus-francia-adolescente-covid-19-paris.html"
	newUrl := models.UrlNew{baseUrl, time.Now()}

	result := scraper.ScrapNewUrl(newUrl)
	assert.NotEqual(t, result.Content, "")
	assert.NotEmpty(t, result.Tags)
	assert.NotEqual(t, result.Date.IsZero(),true)
	assert.NotEqual(t, result.Headline, "")
	fmt.Println(result)

}
