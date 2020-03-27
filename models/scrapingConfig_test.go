package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapingConfig(t *testing.T) {

	config := ScrapingConfig{}
	config.CreateFromJson()
	fmt.Println(config)
	assert.NotEqual(t, nil, config.UrlBase, "OK response is expected")

}
