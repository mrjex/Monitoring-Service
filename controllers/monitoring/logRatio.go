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

// Calculates the percentage of responses relative to requests for a specific service
func CalculatePercentage(service string) (float64, error) {
	requests := database.GetCollection("RequestLogs")
	responses := database.GetCollection("ResponseLogs")

	// Count requests with the specified service
	countRequests, err := requests.CountDocuments(context.TODO(), bson.D{{"service", service}})
	if err != nil {
		return 0, err
	}

	// Count all responses
	countResponses, err := responses.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	}

	if countRequests == 0 {
		return 0, errors.New("cannot calculate percentage with zero requests for the specified service")
	}

	percentage := (float64(countResponses) / float64(countRequests)) * 100
	return percentage, nil
}
