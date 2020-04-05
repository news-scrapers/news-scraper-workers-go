package managers

import (
	"fmt"
	"news-scrapers-workers-go/models"
	"news-scrapers-workers-go/historicScrapers"
	"news-scrapers-workers-go/utils"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ScraperManagerAllHistoricNewspapers struct {
}

func (mainScraper ScraperManagerAllHistoricNewspapers) StartScraping(config models.ScrapingConfig) {

	scraperElmundo := historicScrapers.ElMundoHistoricScraper{Config: config}
	scraperEPais := historicScrapers.ElPaisScraper{Config: config}
	scraperAbc := historicScrapers.AbcHistoricScraper{Config: config}
	scraperLavanguardia := historicScrapers.LaVanguardiaScraper{Config: config}

	log.Info("using historicScrapers:")
	log.Info(config.NewsPaper)


	for {
		var wg sync.WaitGroup

		scrapAll := utils.StringInSlice("all", config.NewsPaper)


		if utils.StringInSlice("elmundo", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperElmundo, "elmundo", config, &wg)
			wg.Add(1)
		}

		if utils.StringInSlice("elpais", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperEPais, "elpais", config, &wg)
			wg.Add(1)
		}

		if utils.StringInSlice("abc", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperAbc, "abc", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("lavanguardia", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperLavanguardia, "lavanguardia", config, &wg)
			wg.Add(1)
		}
		wg.Wait()
		log.Info("-------------------------------------------------------------------------------------------------")
		log.Info("-------------------Finished one iteration, all news from page scraped----------------------------")
		log.Info("-------------------------------------------------------------------------------------------------")

	}
}

func (mainScraper *ScraperManagerAllHistoricNewspapers) ScrapOneIteration(scraper historicScrapers.NewsScraper, newspaper string, config models.ScrapingConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	scrapingIndex, _:= models.GetCurrentIndexNewsPaper(config.ScraperId, newspaper)
	if scrapingIndex ==nil {
		scrapingIndex = models.CreateScrapingIndex(config, newspaper)
	}
	date := scrapingIndex.DateLastNew

	if date.Year() < 2 {
		date = time.Now()
		scrapingIndex.DateLastNew = date
		scrapingIndex.Save()
		log.Info("sating date as now")
	}

	diff := date.Sub(config.ScrapingDateLimit)
	if  diff.Hours()<0{
		log.Info("Reached scraping limit, no more news to scrap in this period")
		panic("Reached limit date for scraping")
	}
	log.Info("-------------------------------------------------------------------------------------------------")
	msg := fmt.Sprintf("starting scraping date %v from %s", date, newspaper)
	log.Info(msg)
	log.Info("-------------------------------------------------------------------------------------------------")


	newsScraped := scraper.ScrapDate(date)

	log.Info("-------------------------------------------------------------------------------------------------")
	log.Info("Saving all scrapded results for " +  newspaper)
	log.Info("-------------------------------------------------------------------------------------------------")

	models.CreateOrUpdateManyNewsScraped(&newsScraped)

	scrapingIndex.DateLastNew = scrapingIndex.DateLastNew.Add(-24*time.Hour)
	scrapingIndex.DateScraping = time.Now()
	scrapingIndex.Save()
}
