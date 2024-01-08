package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	Topic   string             `bson:"topic,omitempty" json:"topic,omitempty"`
	Service string             `bson:"service,omitempty" json:"service,omitempty"`
	Request string             `bson:"request,omitempty" json:"request,omitempty"`
	Time    primitive.DateTime `bson:"time,omitempty" json:"time,omitempty"`
}
