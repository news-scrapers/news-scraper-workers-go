package dailyScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"

	"github.com/joho/godotenv"
)

func TestRecursiveScraperPublico(t *testing.T) {
	t.Skip()
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{ScraperId: "testScraperPublico", DeviceID: "testDevicePublico"}
	index := models.ScrapingIndex{ScraperID: "test", PageIndex: 25}
	scraper := ElPublicoFullIndexManager{Config: config}
	baseUrl := "https://www.publico.es/sociedad/"

	scraper.ScrapNewsInItems(baseUrl, &index)

}
