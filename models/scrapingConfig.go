package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ScrapingConfig struct {
	UrlBase           string    `json:"url_base"`
	DeviceID          string    `json:"device_id"`
	ScraperId         string    `json:"scraper_id"`
	ScraperType       string    `json:"scraper_type"`
	NewsPaper         []string  `json:"newspaper"`
	AppID             string    `json:"app_id"`
	ScrapingDateLimit time.Time `json:"scraping_date_limit"`
	InitialUrls map[string][]string `json:"initial_urls" bson:"initial_urls"`

}

func (config *ScrapingConfig) CreateFromJson() {
	jsonFile, err := os.Open("scrapingConfig.json")
	if err != nil {
		log.Error(err)
	}
	log.Info("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	config.GetStartingUrls()
}

func (config *ScrapingConfig) CreateFromEnvVars() {
	config.UrlBase = os.Getenv("url_base")
	config.DeviceID = os.Getenv("device_id")
	config.ScraperId = os.Getenv("scraper_id")
	config.ScraperType = os.Getenv("scraper_type")
	config.NewsPaper = strings.Split(os.Getenv("newspaper"), ";")
	config.AppID = os.Getenv("url_base")

	limit := os.Getenv("scraping_date_limit")
	t, _ := time.Parse("2006-01-02", limit)
	config.ScrapingDateLimit = t
	config.GetStartingUrls()
}

func (config *ScrapingConfig) GetStartingUrls(){
	jsonFile, err := os.Open("startingUrls.json")
	if err != nil {
		log.Error(err)
	}
	log.Info("Successfully Opened startingUrls.json")
	defer jsonFile.Close()

	var initialUrls map[string][]string = make(map[string][]string)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &initialUrls)
	config.InitialUrls = initialUrls
}
