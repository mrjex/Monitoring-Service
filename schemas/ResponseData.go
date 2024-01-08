package schemas

type ResponseData struct {
	Status         int              `json:"status,omitempty"`
	RequestID      string           `json:"requestID,omitempty"`
	Message        string           `json:"message,omitempty"`
	AvailableTimes *[]AvailableTime `json:"availabletimes,omitempty"`
}
