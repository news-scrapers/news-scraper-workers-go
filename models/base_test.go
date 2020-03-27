package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {

}

func TestConnection(t *testing.T) {
	err := loadDb()
	assert.Nil(t,err, "Connection obtained")

}
