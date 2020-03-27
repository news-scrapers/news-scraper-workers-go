package managers

import "news-scrapers-workers-go/models"

type ScraperManager interface {
	StartScraping(config models.ScrapingConfig)
}
