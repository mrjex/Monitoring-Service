package controllers

import (
	"Monitoring-service/database"
	"Monitoring-service/schemas"
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitialiseLogger(client mqtt.Client) {
	tokenLog := client.Subscribe("grp20/#", byte(0), func(c mqtt.Client, m mqtt.Message) {
		go Log(m)
	})
	if tokenLog.Error() != nil {
		panic(tokenLog.Error())
	}
}

func Log(message mqtt.Message) bool {
	if message.Retained() {
		return false
	}
	var requestLog schemas.Log
	topic := message.Topic()

	requestLog.Topic = topic
	requestLog.Request = string(message.Payload())
	requestLog.Service = GetService(topic)

	currentTime := time.Now()
	requestLog.Time = primitive.NewDateTimeFromTime(currentTime)

	reqRes := GetReqRes(topic)
	if reqRes == "req" {
		col := GetRequestCollection()
		result, err := col.InsertOne(context.TODO(), requestLog)
		_ = result
		return err == nil
	} else if reqRes == "res" {
		col := GetResponseCollection()
		result, err := col.InsertOne(context.TODO(), requestLog)
		_ = result
		return err == nil
	}
	return false
}

func GetRequestCollection() *mongo.Collection {
	col := database.Database.Collection("RequestLogs")
	return col
}

func GetResponseCollection() *mongo.Collection {
	col := database.Database.Collection("ResponseLogs")
	return col
}
func GetCollection() *mongo.Collection {
	col := database.Database.Collection("Logs")
	return col
}
