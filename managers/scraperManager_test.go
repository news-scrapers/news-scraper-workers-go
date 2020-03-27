package managers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"

	"github.com/joho/godotenv"
)



func TestMainAll(t *testing.T) {
	t.Skip()

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	//testConfig := models.ScrapingConfig{UrlBase: "http://localhost:8000", NewsPaper: "elmundo", ScraperId: "testScraper", DeviceID: "testDevice"}
	//testConfig := models.ScrapingConfig{UrlBase: "http://localhost:8000", NewsPaper: "elpais", ScraperId: "testScraperElpais", DeviceID: "testDeviceElpais"}
	testConfig := models.ScrapingConfig{UrlBase: "http://localhost:8000", NewsPaper: []string{"elmundo", "elpais"}, ScraperId: "testScraperAbc", DeviceID: "testDeviceAbc"}

	mainScraper := ScraperManagerAllHistoricNewspapers{}
	mainScraper.StartScraping(testConfig)

}
