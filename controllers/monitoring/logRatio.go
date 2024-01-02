package monitoring

import (
	"Monitoring-service/database"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// This controller uses no time filter, only an empty Bson.D document. This can and will be changed

// Returns the number of requests with no filter
func GetNumOfRequests() (int32, error) {
	requests := database.GetCollection("RequestLogs")

	// Count all documents
	count, err := requests.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	}

	// Convert int64 to int32
	numOfRequests := int32(count)
	return numOfRequests, nil
}

// Returns the number of requests with no filter
func GetNumOfResponses() (int32, error) {
	requests := database.GetCollection("ResponseLogs")

	// Count all documents
	count, err := requests.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	}

	// Convert int64 to int32
	numOfResponses := int32(count)
	return numOfResponses, nil
}

// Calculates the percentage of responses relative to requests
func CalculatePercentage() (float64, error) {
	numOfRequests, err := GetNumOfRequests()
	if err != nil {
		return 0, err
	}

	numOfResponses, err := GetNumOfResponses()
	if err != nil {
		return 0, err
	}

	if numOfRequests == 0 {
		return 0, errors.New("cannot calculate percentage with zero requests")
	}

	percentage := (float64(numOfResponses) / float64(numOfRequests)) * 100
	return percentage, nil
}
