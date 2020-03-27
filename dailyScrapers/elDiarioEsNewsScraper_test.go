package dailyScrapers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestScraperDiario(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperdiario", DeviceID: "testDevicediario"}
	scraper := ElDiaroEsNewsScraper{Config: config}
	baseUrl := "https://www.eldiario.es/sociedad/Balance_0_1008949150.html"
	newUrl := models.UrlNew{baseUrl, time.Now()}

	result := scraper.ScrapNewUrl(newUrl)
	assert.NotEqual(t, result.Content, "")
	assert.NotEmpty(t, result.Tags)
	assert.NotEqual(t, result.Date.IsZero(),true)
	assert.NotEqual(t, result.Headline, "")
	fmt.Println(result)

}
