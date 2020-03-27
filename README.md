# News scraper go

This repository contains the code of the workers that extract data from newspapers **el mundo**, **el pais**, **la vanguardia** and **la razón**. 

## Installation

You will also need **golang** (version 1.12 at least). Follow [this](https://golang.org/doc/install) to install it.

## Configuration
* Clone this repository to a directory:

     git clone https://github.com/news-scrapers/news-scraper-workers-go.git

* Move to the cloned directory and create a file named **.env** . 
Inside this file you will need to specify the url to the mongodb 
database  and the newspapers that you want to scrap. 
Also an id for your scraper. Here is an example that will work if 
you are running the backend locally and if you want to scrap all newspapers:
  
        scraper_id=scraper_test
        device_id=device_test
        app_id=all
        newspaper=all
        scraping_date_limit=2020-03-10
        scraper_type=daily
        database_url =mongodb://localhost:27017
        database_name=news-scraped-with-tags
if you only want to scrap elmundo y el pais you should  replace `newspaper=all` 
with `newspaper=elpais;elmundo` (add them separated by ";"). Right now you can scap the following:
    
    elpais
    elmundo
    abc
    diario.es
    publico
    lavanguardia
    20minutos    
    
if you use `scraper_type=daily` in .env you will scrap the last news in each newspaper, if you want to 
scrap all historic news you should use `scraper_type=historic`, the scraper will stop when it reaches
the limit date for scraping specified in  `scraping_date_limit=2020-03-10`. Historic scraping is only available
for:

    elpais
    elmundo
    abc
    lavanguardia

* Run the golang code. Inside the project folder run:
  
        go run main.go
    you should see the logs with the scraped news of the different newspapers.
    
   Alternatively you can use docker-compose:
        
        docker-compose up
    Or
        
        docker-compose up -d 