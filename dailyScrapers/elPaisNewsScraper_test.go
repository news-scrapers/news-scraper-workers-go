package dailyScrapers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestScraperElPais(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperdiario", DeviceID: "testDevicediario"}
	scraper := ElPaisEsNewsScraper{Config: config}
	baseUrl := "https://elpais.com/sociedad/2020-03-26/china-asegura-que-espana-compro-los-test-rapidos-a-una-empresa-sin-licencia.html"
	newUrl := models.UrlNew{baseUrl, time.Now()}

	result := scraper.ScrapNewUrl(newUrl)
	assert.NotEqual(t, result.Content, "")
	assert.NotEmpty(t, result.Tags)
	assert.NotEqual(t, result.Date.IsZero(),true)
	assert.NotEqual(t, result.Headline, "")
	fmt.Println(result)

}
