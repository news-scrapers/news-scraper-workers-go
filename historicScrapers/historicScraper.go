package historicScrapers

import (
	"news-scrapers-workers-go/models"
	"time"
)

type NewsScraper interface {
	ScrapDate(date time.Time) []models.NewScraped
}
