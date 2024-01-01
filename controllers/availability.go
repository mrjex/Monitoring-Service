package controllers

import (
	"Monitoring-service/database"
	"Monitoring-service/schemas"
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var(
    UserFlag = make(chan bool, 1)
    userResponse = make(chan struct{}, 1)
    recievedMessage string
    userDown bool
    appointmentDown bool
)

func InitialiseAvailability(client mqtt.Client) {
    checkServiceStatus()
    go CheckUserService(client)

    tokenUserService := client.Subscribe("grp20/res/patients/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
        userResponse <- struct{}{}
    })
    if tokenUserService.Error() != nil{
        panic(tokenUserService.Error())
    }
}

func CheckUserService(client mqtt.Client) {
    for{
        requestID := primitive.NewObjectID().Hex()
        message := `{"requestID": "` + requestID + `"}`
        token := client.Publish("grp20/req/patients/get", 0, false, message)
        token.Wait()

        timeout := time.After(5*time.Second)
        select{
        case <- userResponse:
            UserFlag <- true
            closeDownTime(primitive.NewDateTimeFromTime(time.Now()), "User")
            userDown = false
            time.Sleep(5*time.Second)
        case <- timeout:
            UserFlag <- false
            if !userDown{
                downTime := schemas.DownTime{
                    TimeDown: primitive.NewDateTimeFromTime(time.Now()),
                    Service: "User",
                }
                insertDownTimeDB(downTime)
                userDown = true
            }

        }
    }
}

func insertDownTimeDB(downTime schemas.DownTime) {
    col := database.Database.Collection("DownTime")

    col.InsertOne(context.TODO(), downTime)

}

func checkServiceStatus() {
    var downTimes []bson.M
    col := database.Database.Collection("DownTime")
    zeroTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
    zeroTimePrim := primitive.NewDateTimeFromTime(zeroTime)

    filter := bson.M{"time_up": zeroTimePrim, "service": "User"}
    res, err := col.Find(context.TODO(), filter)
    if err != nil{
        fmt.Println("Error getting from database")
        return
    }
    if err := res.All(context.Background(), &downTimes); err != nil{
        fmt.Println("Error decoding")
        return
    }

    if len(downTimes) == 1{
        userDown = true
        fmt.Println("Service down")
    } else if len(downTimes) == 0 {
        userDown = false
        fmt.Println("Service not down")
    }
}

func closeDownTime(upTime primitive.DateTime, service string) {
    col := database.Database.Collection("DownTime")
    zeroTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
    zeroTimePrim := primitive.NewDateTimeFromTime(zeroTime)

    filter := bson.M{"service": service,"time_up": zeroTimePrim}
    update := bson.M{"$set": bson.M{"time_up": upTime}}

    res, err := col.UpdateOne(context.TODO(), filter, update)
    _ = res
    if err != nil{
        fmt.Println("Error updating db")
    }
}
