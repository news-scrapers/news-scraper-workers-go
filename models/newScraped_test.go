package models

import (
	"fmt"
	"testing"
)

func TestPostNew(t *testing.T) {
	config := ScrapingConfig{}
	config.CreateFromJson()

	newScraped := NewScraped{}
	newScraped.NewsPaper = "elpais"
	newScraped.Url = "http://test"
	newScraped.Content = "test"
	newScraped.Headline = "test headline"
	newScraped.ScraperID = config.ScraperId
	err := newScraped.Create()
	if err != nil {
		fmt.Println(err)
	}
	//assert.NotEqual(t, nil, err, "no error")
	fmt.Println(newScraped)

}
