package dailyScrapers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestScraperElMundo(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperdiario", DeviceID: "testDevicediario"}
	scraper := ElMundoNewsScraper{Config: config}
	baseUrl := "https://www.elmundo.es/economia/2020/03/26/5e7d0cf2fdddff496a8b465c.html"
	newUrl := models.UrlNew{baseUrl, time.Now()}

	result := scraper.ScrapNewUrl(newUrl)
	assert.NotEqual(t, result.Content, "")
	assert.NotEmpty(t, result.Tags)
	assert.NotEqual(t, result.Date.IsZero(),true)
	assert.NotEqual(t, result.Headline, "")
	fmt.Println(result)

}
