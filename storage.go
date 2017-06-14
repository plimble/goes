package goes

// Snapshot entity
type Snapshot struct {
	AggregateID      string `json:"aggregateID" dynamodbav:"aggregateID" db:"aggregateID"`
	AggregateType    string `json:"aggregateType" dynamodbav:"aggregateType" db:"aggregateType"`
	AggregateVersion int    `json:"aggregateVersion" dynamodbav:"aggregateVersion" db:"aggregateVersion"`
	Revision         int    `json:"revision" dynamodbav:"revision" db:"revision"`
	Time             int64  `json:"time" dynamodbav:"time" db:"time"`
	Data             []byte `json:"data" dynamodbav:"data" db:"data"`
}

// Storage interface
//go:generate mockery -name Storage
type Storage interface {
	Save(es []Event) error
	GetLastEvent(id string) ([]Event, error)
	SaveSnapshot(snap *Snapshot) error
	GetSnapshot(id string, version int) (*Snapshot, error)
	GetFromRevision(id string, from int) ([]Event, error)
	GetUndispatchedEvent() ([]Event, error)
	MarkDispatchedEvent(es []Event) error
}
