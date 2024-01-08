package monitoring

import (
	"Monitoring-service/database"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CalculatePercentage Calculates the percentage of responses relative to requests for a specific service
func CalculatePercentage(service string) (float64, error) {
	requests := database.GetCollection("RequestLogs")
	responses := database.GetCollection("ResponseLogs")

	var filter bson.D
	if service != "AllServices" {
		filter = bson.D{{"service", service}}
	} else {
		filter = bson.D{{}}
	}

	// Count requests based on the filter
	countRequests, err := requests.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	// Count all responses based on the filter
	countResponses, err := responses.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	if countRequests == 0 {
		return 0, errors.New("cannot calculate percentage with zero requests for the specified service")
	}

	percentage := (float64(countResponses) / float64(countRequests)) * 100
	return percentage, nil
}
