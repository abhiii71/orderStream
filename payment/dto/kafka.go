package dto

// EventData mirrors product /models.EventData  used by Product  service kafka messages
//
//	we only care about product creation  events here.
type EventData struct {
	ID          *string
	Name        *string
	Description *string
	Price       *float64
	AccountID   *int
}

// EVENT mirrors product/models.Event
type Event struct {
	Type string
	Data EventData
}
