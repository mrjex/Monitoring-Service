package controllers

import (
	"Monitoring-service/database"
	"Monitoring-service/schemas"
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	AppointmentFlag     = make(chan bool, 1)
	appointmentResponse = make(chan struct{}, 1)
	appointmentDown     bool

	UserFlag     = make(chan bool, 1)
	userResponse = make(chan struct{}, 1)
	userDown     bool

	AvailableTimesChan = make(chan schemas.ResponseData, 1)
)

func InitialiseAvailability(client mqtt.Client) {
	checkServiceStatus()
	go CheckUserService(client)
	go CheckAppointmentService(client)

	tokenUserService := client.Subscribe("grp20/res/patients/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
		userResponse <- struct{}{}
	})
	if tokenUserService.Error() != nil {
		panic(tokenUserService.Error())
	}

	tokenAppointmentService := client.Subscribe("grp20/res/availabletimes/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
		var payload schemas.ResponseData

		err1 := json.Unmarshal(m.Payload(), &payload)
		if err1 != nil {
			go countAvailableTimes(payload)
			return
		}
		AvailableTimesChan <- payload
		appointmentResponse <- struct{}{}
	})
	if tokenAppointmentService.Error() != nil {
		panic(tokenAppointmentService.Error())
	}

}

func CheckUserService(client mqtt.Client) {
	for {
		requestID := primitive.NewObjectID().Hex()
		message := `{"requestID": "` + requestID + `"}`
		token := client.Publish("grp20/req/patients/get", 0, false, message)
		token.Wait()

		timeout := time.After(5 * time.Second)
		select {
		case <-userResponse:
			UserFlag <- true
			closeDownTime(primitive.NewDateTimeFromTime(time.Now()), "User")
			userDown = false
			time.Sleep(5 * time.Second)
		case <-timeout:
			UserFlag <- false
			if !userDown {
				downTime := schemas.DownTime{
					TimeDown: primitive.NewDateTimeFromTime(time.Now()),
					Service:  "User",
				}
				insertDownTimeDB(downTime)
				userDown = true
			}

		}
	}
}

func CheckAppointmentService(client mqtt.Client) {
	for {
		requestID := primitive.NewObjectID().Hex()
		message := `{"requestID": "` + requestID + `"}`
		token := client.Publish("grp20/req/availabletimes/get", 0, false, message)
		token.Wait()

		timeout := time.After(5 * time.Second)
		select {
		case <-appointmentResponse:
			AppointmentFlag <- true
			closeDownTime(primitive.NewDateTimeFromTime(time.Now()), "Appointment")
			appointmentDown = false
			time.Sleep(5 * time.Second)
		case <-timeout:
			AppointmentFlag <- false
			if !appointmentDown {
				downTime := schemas.DownTime{
					TimeDown: primitive.NewDateTimeFromTime(time.Now()),
					Service:  "Appointment",
				}
				insertDownTimeDB(downTime)
				appointmentDown = true
			}

		}
	}
}

func insertDownTimeDB(downTime schemas.DownTime) {
	col := database.Database.Collection("DownTime")

	col.InsertOne(context.TODO(), downTime)

}

func checkServiceStatus() {
	var userDownTimes []bson.M
	col := database.Database.Collection("DownTime")
	zeroTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	zeroTimePrim := primitive.NewDateTimeFromTime(zeroTime)

	// User service
	filter := bson.M{"time_up": zeroTimePrim, "service": "User"}
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error getting from database")
		return
	}
	if err := res.All(context.Background(), &userDownTimes); err != nil {
		fmt.Println("Error decoding")
		return
	}

	if len(userDownTimes) == 1 {
		userDown = true
		fmt.Println("Service down")
	} else if len(userDownTimes) == 0 {
		userDown = false
		fmt.Println("Service not down")
	}

	// Appointment service
	var appointDownTimes []bson.M
	filter = bson.M{"time_up": zeroTimePrim, "service": "Appointment"}
	res, err = col.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error getting from database")
		return
	}
	if err = res.All(context.Background(), &appointDownTimes); err != nil {
		fmt.Println("Error decoding")
		return
	}

	if len(appointDownTimes) == 1 {
		appointmentDown = true
		fmt.Println("Service down")
	} else if len(appointDownTimes) == 0 {
		appointmentDown = false
		fmt.Println("Service not down")
	}
}

func closeDownTime(upTime primitive.DateTime, service string) {
	col := database.Database.Collection("DownTime")
	zeroTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	zeroTimePrim := primitive.NewDateTimeFromTime(zeroTime)

	filter := bson.M{"service": service, "time_up": zeroTimePrim}
	update := bson.M{"$set": bson.M{"time_up": upTime}}

	res, err := col.UpdateOne(context.TODO(), filter, update)
	_ = res
	if err != nil {
		fmt.Println("Error updating db")
	}
}

func countAvailableTimes(payload schemas.ResponseData) int {
	if payload.AvailableTimes != nil {
		availableTimes := *payload.AvailableTimes
		return len(availableTimes)
	}
	return 0
}
