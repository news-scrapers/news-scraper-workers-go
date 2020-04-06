package models

import (
	"context"
	"errors"
	"time"

	log "log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScrapingIndex struct {
	DateLastNew     time.Time `json:"date_last_new" bson:"date_last_new"`
	DateScraping    time.Time `json:"date_scraping" bson:"date_scraping"`
	LastHistoricUrl string    `json:"last_historic_url" bson:"last_historic_url"`
	Page            int64     `json:"page" bson:"page"`
	PageIndex     	int       `json:"page_index" bson:"page_index"`
	UrlIndex        int     `json:"url_index" bson:"url_index"`
	NewsPaper       string    `json:"newspaper" bson:"newspaper"`
	ScraperType     string    `json:"scraper_type" bson:"scraper_type"`
	ScraperID       string    `json:"scraper_id" bson:"scraper_id"`
	DeviceID        string    `json:"device_id" bson:"device_id"`
	StartingUrls  []string  `json:"startingUrls" bson:"startingUrls"`
}

func (scrapingIndex *ScrapingIndex) Save() error {
	scrapingIndex.DateScraping = time.Now()
	db := GetDB()
	collection := db.Collection("ScrapingIndex")
	options := options.FindOneAndReplaceOptions{}
	upsert := true
	options.Upsert = &upsert

	log.Println("saving scraping index for id" + scrapingIndex.ScraperID + " and newspaper ")
	log.Println(scrapingIndex.NewsPaper)
	err := collection.FindOneAndReplace(context.Background(), bson.M{"scraper_id": scrapingIndex.ScraperID, "scraper_type":scrapingIndex.ScraperType, "newspaper": scrapingIndex.NewsPaper}, scrapingIndex, &options)

	if err != nil {
		return errors.New("error finding")
	} else {
		return nil
	}

}

func CreateScrapingIndex(config ScrapingConfig, newspaper string) *ScrapingIndex {
	scrapingIndexNew := ScrapingIndex{}
	scrapingIndexNew.DateLastNew = time.Now()
	scrapingIndexNew.NewsPaper = newspaper
	scrapingIndexNew.ScraperType = config.ScraperType
	scrapingIndexNew.ScraperID = config.ScraperId
	scrapingIndexNew.DateLastNew = time.Now()
	scrapingIndexNew.StartingUrls = config.InitialUrls[newspaper]
	scrapingIndexNew.UrlIndex = 0
	scrapingIndexNew.PageIndex = 1
	scrapingIndexNew.DeviceID = config.DeviceID
	return &scrapingIndexNew
}

func GetCurrentIndex(scraperID string, source string, scraperType string) (scrapingIndex *ScrapingIndex, err error) {
	db := GetDB()
	collection := db.Collection("ScrapingIndex")

	options := options.FindOneOptions{}
	// Sort by `_id` field descending
	options.Sort = bson.D{{"date_last_new", int32(1)}}

	results := &ScrapingIndex{}
	err = collection.FindOne(context.Background(), bson.M{"scraper_id": scraperID, "scraper_type": scraperType, "newspaper": source}, &options).Decode(&results)
	if err != nil {
		return scrapingIndex, err
	}
	scrapingIndex = results

	return scrapingIndex, nil

}


func (scrapingIndex *ScrapingIndex) UpdateUrls(config ScrapingConfig, source string) error {
	log.Println("Updating initial urls for scraping in " + source)
	scrapingIndex.StartingUrls = config.InitialUrls[source]
	scrapingIndex.Save()

	return nil
}


func GetCurrentIndexNewsPaper(scraperID string, newsPaper string) (scrapingIndex *ScrapingIndex, err error) {
	db := GetDB()
	collection := db.Collection("ScrapingIndex")

	options := options.FindOneOptions{}
	// Sort by `_id` field descending
	options.Sort = bson.D{{"date_last_new", int32(1)}}

	results := ScrapingIndex{}
	err = collection.FindOne(context.Background(), bson.M{"scraper_id": scraperID, "newspaper": newsPaper}, &options).Decode(&results)
	if err != nil {
		return nil, err
	}
	scrapingIndex = &results
	return scrapingIndex, nil

}
