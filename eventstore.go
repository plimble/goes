package goes

import (
	"encoding/json"

	"github.com/plimble/errors"
)

// EventStore interface
type EventStore interface {
	NewAggregateBase(id string, aggtype string, version int, applyFunc ApplyFunc) *AggregateBase
	GetStream(a Aggregate) error
	GetStreamWithSnapshot(a Aggregate) error
	Commit(a Aggregate) error
	SaveSnapshot(a Aggregate) error
}

type eventStore struct {
	storage Storage
	encoder Encoder
	idgen   IDGen
}

// NewEventStore create new event store
func NewEventStore(storage Storage, encoder Encoder, idgen IDGen) EventStore {
	return &eventStore{storage, encoder, idgen}
}

func (e *eventStore) NewAggregateBase(id string, aggtype string, version int, applyFunc ApplyFunc) *AggregateBase {
	return &AggregateBase{
		ID:         id,
		Type:       aggtype,
		Version:    version,
		encoder:    e.encoder,
		idgen:      e.idgen,
		applyEvent: applyFunc,
	}
}

func (e *eventStore) GetStream(a Aggregate) error {
	b := a.GetAggregateBase()
	events, err := e.storage.GetLastEvent(b.ID)
	if err != nil {
		return err
	}

	last := events[len(events)-1]
	b.Revision = last.Revision
	b.Time = last.Time

	for _, e := range events {
		if err = b.applyEvent(e); err != nil {
			return err
		}
	}

	return nil
}

func (e *eventStore) GetStreamWithSnapshot(a Aggregate) error {
	b := a.GetAggregateBase()
	snap, err := e.storage.GetSnapshot(b.ID, b.Version)
	if err != nil {
		if errors.IsNotFound(err) {
			return e.GetStream(a)
		}
		return err
	}

	b.Revision = snap.Revision
	b.Time = snap.Time

	if err := e.encoder.Decode(snap.Data, a); err != nil {
		return err
	}

	events, err := e.storage.GetFromRevision(b.ID, snap.Revision)
	if err != nil {
		return err
	}

	last := events[len(events)-1]
	b.Revision = last.Revision
	b.Time = last.Time

	for _, e := range events {
		if err = b.applyEvent(e); err != nil {
			return err
		}
	}

	return nil
}

// SaveSnapshot save snapshot
func (e *eventStore) SaveSnapshot(a Aggregate) error {
	b := a.GetAggregateBase()
	data, _ := json.Marshal(a)
	snap := &Snapshot{
		AggregateID:      b.ID,
		AggregateType:    b.Type,
		AggregateVersion: b.Version,
		Revision:         b.Revision,
		Time:             b.Time,
		Data:             data,
	}
	return e.storage.SaveSnapshot(snap)
}

// Commit save the event into store
func (e *eventStore) Commit(a Aggregate) error {
	b := a.GetAggregateBase()
	err := e.storage.Save(b.events)
	if err != nil {
		return err
	}

	go func() {
		b := a.GetAggregateBase()
		if b.Revision%20 == 0 {
			data, _ := json.Marshal(a)
			snap := &Snapshot{
				AggregateID:      b.ID,
				AggregateType:    b.Type,
				AggregateVersion: b.Version,
				Revision:         b.Revision,
				Time:             b.Time,
				Data:             data,
			}
			e.storage.SaveSnapshot(snap)
		}
	}()

	return nil
}
