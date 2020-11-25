package historicScrapers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestScraperAbc(t *testing.T) {
	t.Skip()
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	config := models.ScrapingConfig{}
	config.CreateFromJson()
	scraper := AbcHistoricScraper{Config: config}
	date := time.Date(2018, 05, 04, 0, 0, 0, 0, time.UTC)
	results := scraper.ScrapDate(date)
	result := results[1]

	assert.NotEqual(t, result.Headline, "", "OK response is expected")

}
