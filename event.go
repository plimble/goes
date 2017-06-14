package goes

// Event entity
//go:generate msgp
type Event struct {
	AggregateID   string `json:"aggregateID" dynamodbav:"aggregateID" db:"aggregateID" msg:"a"`
	AggregateType string `json:"aggregateType" dynamodbav:"aggregateType" db:"aggregateType" msg:"b"`
	EventID       string `json:"eventID" dynamodbav:"eventID" db:"eventID" msg:"c"`
	EventType     string `json:"eventType" dynamodbav:"eventType" db:"eventType" msg:"d"`
	Revision      int    `json:"revision" dynamodbav:"revision" db:"revision" msg:"e"`
	Time          int64  `json:"time" dynamodbav:"time" db:"time" msg:"f"`
	Data          []byte `json:"data" dynamodbav:"data" db:"data" msg:"g"`
}
