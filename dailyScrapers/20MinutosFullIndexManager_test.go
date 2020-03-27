package dailyScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"

	"github.com/joho/godotenv"
)

func TestRecursiveScraper20Minutos(t *testing.T) {
	//t.Skip()
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{UrlBase: "http://localhost:8000", ScraperId: "testScraperElpais", DeviceID: "testDeviceElpais"}
	index := models.ScrapingIndex{ScraperID: "test", PageIndex: 1}
	scraper := VeinteMinutosFullIndexManager{Config: config}

	//baseUrl := "https://www.amazon.es/gp/bestsellers/?ref_=nav_cs_bestsellers"
	baseUrl := "https://www.20minutos.es/internacional/"

	scraper.ScrapNewsInItems(baseUrl, &index)

}
