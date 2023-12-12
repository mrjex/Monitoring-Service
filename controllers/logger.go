package controllers

import (
	"Monitoring-service/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeLogger(client mqtt.Client) {

}

func GetAppointmentCollection() *mongo.Collection {
	col := database.Database.Collection("Logs")
	return col
}
