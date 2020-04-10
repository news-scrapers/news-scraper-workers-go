package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UrlNew struct {
	Url  string
	Date time.Time
}

type NewScraped struct {
	Page       int       `json:"page" bson:"page"`
	FullPage   bool      `json:"full_page" bson:"full_page"`
	Headline   string    `json:"headline" bson:"headline"`
	Date       time.Time `json:"date" bson:"date"`
	DateString string    `json:"date_string" bson:"date_string"`
	Content    string    `json:"content" bson:"content"`
	Url        string    `json:"url" bson:"url"`
	NewsPaper  string    `json:"newspaper" bson:"newspaper"`
	ScraperID  string    `json:"scraper_id" bson:"scraper_id"`
	ID         string    `json:"id" bson:"id"`
	Tags       []string  `json:"tags" bson:"tags"`
	TagsConcatForTextSearch       string  `json:"tags_concat_for_text_search" bson:"tags_concat_for_text_search"`
}

const collectionNameNewsScraped = "NewsContentScraped"

func init(){
	db := GetDB()
	PopulateTextIndex(db, collectionNameNewsScraped)
}

func PopulateTextIndex(db *mongo.Database, collectionName string) {
	log.Println("Creating text index for collection  " + collectionName)

	coll := db.Collection(collectionName)
	index := []mongo.IndexModel{{Keys: bsonx.Doc{{Key: "tags_concat_for_text_search", Value: bsonx.String("text")}}}}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, errIndex := coll.Indexes().CreateMany(context.Background(), index, opts)
	if errIndex != nil {
		panic(errIndex)
	}
}

func (newScraped *NewScraped) Create() error {

	db := GetDB()
	collection := db.Collection(collectionNameNewsScraped)

	// options := options.FindOneAndReplaceOptions{}
	// upsert := true
	// options.Upsert = &upsert
	// err := collection.FindOneAndReplace(context.Background(), bson.M{"url": newScraped.Url}, newScraped, &options)
	result := &NewScraped{}
	err := collection.FindOne(context.Background(), bson.M{"url": newScraped.Url}).Decode(result)

	return err

	//_, err := collection.InsertOne(context.Background(), newScraped)

}


func (user *NewScraped) SaveOrUpdate() error {
	db := GetDB()
	collection := db.Collection(collectionNameNewsScraped)
	filter := bson.M{"headline": &user.Headline, "newspaper": &user.NewsPaper, "url": &user.Url}
	update := bson.M{"$set": user}
	user.ConcatenateTagsForSearches()
	_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	if err != nil {
		log.Info(err)
		return err

	}
	return nil
}
func (newScraped *NewScraped) ConcatenateTagsForSearches() {
	newScraped.TagsConcatForTextSearch = strings.Join(newScraped.Tags[:], " ")
}

func CreateOrUpdateManyNewsScraped(newsScraped *[]NewScraped) error {
	for _, result := range *newsScraped {
		result.SaveOrUpdate()
	}

	return nil
}
