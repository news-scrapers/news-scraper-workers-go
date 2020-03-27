package dailyScrapers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestScraper20Minutos(t *testing.T) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperdiario", DeviceID: "testDevicediario"}
	scraper := VeinteMinutosNewsScraper{Config: config}
	baseUrl := "https://www.20minutos.es/noticia/4206112/0/coronavirus-autonomos-madrid-requisitos-prestacion-cese-actividad/"
	newUrl := models.UrlNew{baseUrl, time.Now()}

	result := scraper.ScrapNewUrl(newUrl)
	assert.NotEqual(t, result.Content, "")
	assert.NotEqual(t, result.Date.IsZero(),true)
	assert.NotEmpty(t, result.Tags)
	assert.NotEqual(t, result.Headline, "")
	fmt.Println(result)

}
