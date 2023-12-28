package controllers

import (
	"time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var(
    UserFlag = make(chan bool, 1)
    userResponse = make(chan struct{}, 1)
    recievedMessage string
)

func InitialiseAvailability(client mqtt.Client) {
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
            time.Sleep(5*time.Second)
        case <- timeout:
            UserFlag <- false
        }
    }
}
