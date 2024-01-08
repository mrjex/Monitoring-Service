package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type DownTime struct{
    TimeDown primitive.DateTime `bson:"time_down" json:"time_down"`
    TimeUp   primitive.DateTime `bson:"time_up" json:"time_up,omitempty"`
    Service  string             `bson:"service" json:"service,omitempty"`
}
