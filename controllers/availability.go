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

	NotificationFlag     = make(chan bool, 1)
	notificationResponse = make(chan struct{}, 1)
	notificationDown     bool

	ClinicFlag         = make(chan bool, 1)
	clinicResponse     = make(chan struct{}, 1)
	clinicDown         bool
	AvailableTimesChan = make(chan schemas.ResponseData, 1)
)

func InitialiseAvailability(client mqtt.Client) {
	checkServiceStatus()
	go CheckUserService(client)
	go CheckAppointmentService(client)
	go CheckNotificationService(client)
	go CheckClinicService(client)

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

	tokenNotificationService := client.Subscribe("grp20/res/subscriber/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
		notificationResponse <- struct{}{}
	})
	if tokenNotificationService.Error() != nil {
		panic(tokenNotificationService.Error())
	}

	tokenClinicService := client.Subscribe("grp20/res/map/nearby", byte(0), func(c mqtt.Client, m mqtt.Message) {
		clinicResponse <- struct{}{}
	})
	if tokenClinicService.Error() != nil {
		panic(tokenClinicService.Error())
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

// Notification service
func CheckNotificationService(client mqtt.Client) {
	for {
		requestID := primitive.NewObjectID().Hex()
		patientID := primitive.NewObjectID().Hex()
		message := `{"requestID": "` + requestID + `", "patient_ID": "` + patientID + `"}`
		token := client.Publish("grp20/req/subscriber/get", 0, false, message)
		token.Wait()

		timeout := time.After(5 * time.Second)
		select {
		case <-notificationResponse:
			NotificationFlag <- true
			closeDownTime(primitive.NewDateTimeFromTime(time.Now()), "Notification")
			notificationDown = false
			time.Sleep(5 * time.Second)
		case <-timeout:
			NotificationFlag <- false
			if !notificationDown {
				downTime := schemas.DownTime{
					TimeDown: primitive.NewDateTimeFromTime(time.Now()),
					Service:  "Notification",
				}
				insertDownTimeDB(downTime)
				notificationDown = true
			}
		}
	}
}

// Clinic service
func CheckClinicService(client mqtt.Client) {
	for {
		requestID := primitive.NewObjectID().Hex()
		message := `{"requestID": "` + requestID + `", "nearby_clinics_number": "4", "reference_position": "50.13,12.10"}`
		token := client.Publish("grp20/req/map/query/nearby/fixed/get", 0, false, message)
		token.Wait()

		timeout := time.After(5 * time.Second)
		select {
		case <-clinicResponse:
			ClinicFlag <- true
			closeDownTime(primitive.NewDateTimeFromTime(time.Now()), "Clinic")
			clinicDown = false
			time.Sleep(5 * time.Second)
		case <-timeout:
			ClinicFlag <- false
			if !clinicDown {
				downTime := schemas.DownTime{
					TimeDown: primitive.NewDateTimeFromTime(time.Now()),
					Service:  "Clinic",
				}
				insertDownTimeDB(downTime)
				clinicDown = true
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
	} else if len(userDownTimes) == 0 {
		userDown = false
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
	} else if len(appointDownTimes) == 0 {
		appointmentDown = false
	}

	//Notification servicee
	var notifDownTimes []bson.M
	filter = bson.M{"time_up": zeroTimePrim, "service": "Notification"}
	res, err = col.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error getting from database")
		return
	}
	if err = res.All(context.Background(), &notifDownTimes); err != nil {
		fmt.Println("Error decoding")
		return
	}

	if len(notifDownTimes) == 1 {
		notificationDown = true
	} else if len(notifDownTimes) == 0 {
		notificationDown = false
	}

	// Clinic service
	var clinicDownTimes []bson.M
	filter = bson.M{"time_up": zeroTimePrim, "service": "Clinic"}
	res, err = col.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error getting from database")
		return
	}
	if err = res.All(context.Background(), &clinicDownTimes); err != nil {
		fmt.Println("Error decoding")
		return
	}

	if len(clinicDownTimes) == 1 {
		clinicDown = true
	} else if len(clinicDownTimes) == 0 {
		clinicDown = false
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
