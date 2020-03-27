package managers

import (
	"news-scrapers-workers-go/dailyScrapers"
	"news-scrapers-workers-go/models"
	"news-scrapers-workers-go/utils"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ScraperManagerAllDailyNewspapers struct {
}

func (mainScraper ScraperManagerAllDailyNewspapers) StartScraping(config models.ScrapingConfig) {

	scraperDiarioEs := dailyScrapers.ElDiarioEsFullIndexManager{Config: config}
	scraperPublico := dailyScrapers.ElPublicoFullIndexManager{Config: config}
	scraperElpais := dailyScrapers.ElPaisEsFullIndexManager{Config: config}
	scraperElMundo := dailyScrapers.ElMundoFullIndexManager{Config: config}
	scraper20Minutos := dailyScrapers.VeinteMinutosFullIndexManager{Config: config}
	scraperLaVanguardia := dailyScrapers.LaVanguardiaFullIndexManager{Config: config}

	log.Info("using daily Scrapers:")
	for {
		var wg sync.WaitGroup

		scrapAll := utils.StringInSlice("all", config.NewsPaper)

		if utils.StringInSlice("eldiario.es", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperDiarioEs, "eldiario.es", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("publico", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperPublico, "publico", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("elpais", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperElpais, "elpais", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("elmundo", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperElMundo, "elmundo", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("20minutos", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraper20Minutos, "20minutos", config, &wg)
			wg.Add(1)
		}
		if utils.StringInSlice("lavanguardia", config.NewsPaper) || scrapAll {
			go mainScraper.ScrapOneIteration(scraperLaVanguardia, "lavanguardia", config, &wg)
			wg.Add(1)
		}

		wg.Wait()
		log.Info("-------------------------------------------------------------------------------------------------")
		log.Info("-------------------Finished one iteration, all news from page scraped----------------------------")
		log.Info("-------------------------------------------------------------------------------------------------")
	}

}

func (mainScraper *ScraperManagerAllDailyNewspapers) ScrapOneIteration(scraper dailyScrapers.DailyScraper, source string, config models.ScrapingConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("starting scraping using " + source)
	scrapingIndex, err := models.GetCurrentIndex(config.ScraperId, source)

	if scrapingIndex.ScraperType == "" || err != nil {
		scrapingIndex = *models.CreateScrapingIndex(config, source)
	}

	scrapingIndex.UpdateUrls(config, source)

	index := scrapingIndex.UrlIndex
	if index >= len(scrapingIndex.StartingUrls) {
		index = 0
	}

	log.Printf("starting with url number %d", index)

	nextUrl := scrapingIndex.StartingUrls[index]
	scraper.ScrapNewsInItems(nextUrl, &scrapingIndex)

	//models.SaveMany(results, config)

	index = index + 1

	scrapingIndex.UrlIndex = index

	scrapingIndex.Save()
}
