package main

import (
	"fmt"
	"io"
	"news-scrapers-workers-go/managers"
	"news-scrapers-workers-go/models"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const logPath = "logs.txt"

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	//log.SetOutput(f)
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	config := models.ScrapingConfig{}
	config.CreateFromEnvVars()
	fmt.Println(config)

	var mainScraper managers.ScraperManager


	if (config.ScraperType == "historic"){
		mainScraper = managers.ScraperManagerAllHistoricNewspapers{}
	} else {
		mainScraper = managers.ScraperManagerAllDailyNewspapers{}
	}

	mainScraper.StartScraping(config)
}
