package goes

import (
	"time"
)

// ApplyFunc event handler func
type ApplyFunc func(e Event) error

// Aggregate interface
type Aggregate interface {
	GetAggregateBase() *AggregateBase
}

// AggregateBase entity
type AggregateBase struct {
	ID         string
	Type       string
	Revision   int
	Version    int
	Time       int64
	events     []Event
	encoder    Encoder
	idgen      IDGen
	applyEvent ApplyFunc
}

// GetAggregateBase return AggregateBase
func (a *AggregateBase) GetAggregateBase() *AggregateBase {
	return a
}

// DecodeEvent decode event data
func (a *AggregateBase) DecodeEvent(data []byte, v interface{}) error {
	return a.encoder.Decode(data, v)
}

// AddEvent add and apply event
func (a *AggregateBase) AddEvent(eventType string, eventData interface{}) {
	data, _ := a.encoder.Encode(eventData)
	a.Revision++
	a.Time = time.Now().Unix()

	e := Event{
		AggregateID:   a.ID,
		AggregateType: a.Type,
		EventID:       a.idgen.Generate(),
		EventType:     eventType,
		Revision:      a.Revision,
		Time:          a.Time,
		Data:          data,
	}

	if err := a.applyEvent(e); err == nil {
		a.events = append(a.events, e)
	}
}
