package dailyScrapers

import (
	"news-scrapers-workers-go/models"
)

type DailyScraper interface {
	ScrapNewsInItems(baseUrl string, scrapingIndex *models.ScrapingIndex)
}
