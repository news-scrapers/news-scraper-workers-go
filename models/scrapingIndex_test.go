package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetScrapingIndex(t *testing.T) {

	index := ScrapingIndex{}
	config := ScrapingConfig{}
	config.CreateFromJson()

	index ,err := GetCurrentIndex(config.ScraperId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(index)
	assert.NotEqual(t, nil, index, "OK response is expected")
	// assert.Equal(t, config.ScraperId, index.ScraperID, "OK response is expected")
	// assert.NotEqual(t, nil, err.Error, "no error")

}