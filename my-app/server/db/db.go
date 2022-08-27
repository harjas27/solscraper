package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ramdeuter.org/solscraper/model"
	"ramdeuter.org/solscraper/query"
)

type DB interface {
	GetTechnologies() ([]*model.Technology, error)
	GetMetadata() ([]*model.Meta, error)
	GetData() ([]*model.Data, error)
	GetDataForQuery(name string) ([]*model.Data, error)
	SaveData(qName string, data []string)
	SaveQuery(q query.Query)
}

type MongoDB struct {
	collection *mongo.Collection
	meta       *mongo.Collection
	data       *mongo.Collection
}

func NewMongo(client *mongo.Client) DB {
	tech := client.Database("tech").Collection("tech")
	meta := client.Database("tech").Collection("meta")
	data := client.Database("tech").Collection("data")
	return MongoDB{
		collection: tech,
		meta:       meta,
		data:       data,
	}
}

func (m MongoDB) GetTechnologies() ([]*model.Technology, error) {
	res, err := m.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error while fetching technologies:", err.Error())
		return nil, err
	}
	var tech []*model.Technology
	err = res.All(context.TODO(), &tech)
	if err != nil {
		log.Println("Error while decoding technologies:", err.Error())
		return nil, err
	}
	return tech, nil
}

func (m MongoDB) GetMetadata() ([]*model.Meta, error) {
	res, err := m.meta.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error while fetching technologies:", err.Error())
		return nil, err
	}
	var tech []*model.Meta
	err = res.All(context.TODO(), &tech)
	if err != nil {
		log.Println("Error while decoding technologies:", err.Error())
		return nil, err
	}
	return tech, nil
}

func (m MongoDB) GetData() ([]*model.Data, error) {
	res, err := m.data.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error while fetching technologies:", err.Error())
		return nil, err
	}
	var tech []*model.Data
	err = res.All(context.TODO(), &tech)
	if err != nil {
		log.Println("Error while decoding technologies:", err.Error())
		return nil, err
	}
	return tech, nil
}

func (m MongoDB) GetDataForQuery(name string) ([]*model.Data, error) {
	res, err := m.data.Find(context.TODO(), bson.M{"name": name})
	if err != nil {
		log.Println("Error while fetching technologies:", err.Error())
		return nil, err
	}
	var tech []*model.Data
	err = res.All(context.TODO(), &tech)
	if err != nil {
		log.Println("Error while decoding technologies:", err.Error())
		return nil, err
	}
	return tech, nil
}

func (m MongoDB) SaveQuery(q query.Query) {
	data := model.Meta{
		Name:  "test",
		Query: q,
	}
	//b, _ := bson.Marshal(data)
	one, err := m.meta.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Println(one.InsertedID)
}

func (m MongoDB) SaveData(qName string, data []string) {
	for _, row := range data {
		formattedRow := model.Data{
			Name: qName,
			Row:  row,
		}
		b, _ := bson.Marshal(formattedRow)
		one, err := m.data.InsertOne(context.TODO(), b)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		fmt.Println(one.InsertedID)
	}
}
