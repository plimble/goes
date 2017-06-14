package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/plimble/errors"
	"github.com/plimble/goes"
)

// Storage mysql
type Storage struct {
	eventTable                string
	undispatchCommitedTable   string
	snapshotTable             string
	db                        *sqlx.DB
	getLastEventQuery         string
	saveSnapshotQuery         string
	getSnapshotQuery          string
	getFromRevisionQuery      string
	getUndispatchedEventQuery string
	markDispatchedEventQuery  string
}

// New storage
func New(db *sqlx.DB, eventTable, undispatchCommitedTable, snapshotTable string) *Storage {
	return &Storage{
		eventTable,
		undispatchCommitedTable,
		snapshotTable,
		db,
		fmt.Sprintf("SELECT * FROM %s WHERE aggregateID = ?", eventTable),
		fmt.Sprintf("INSERT INTO %s (aggregateID, aggregateType, aggregateVersion, revision, time, data) VALUES (?,?,?,?,?,?)", snapshotTable),
		fmt.Sprintf("SELECT * FROM %s WHERE aggregateID = ? AND aggregateVersion = ?", snapshotTable),
		fmt.Sprintf("SELECT * FROM %s WHERE aggregateID = ? AND revision > ?", eventTable),
		fmt.Sprintf("SELECT * FROM %s", undispatchCommitedTable),
		fmt.Sprintf("DELETE FROM %s WHERE aggregateID = ? AND revision = ?", undispatchCommitedTable),
	}
}

// Save events
func (s *Storage) Save(es []goes.Event) error {
	total := len(es)
	if total == 0 {
		return nil
	}

	dataSqls := make([]string, total)
	valSqls := make([]interface{}, total*7)
	for i, e := range es {
		dataSqls[i] = "(?, ?, ?, ?, ?, ?, ?)"
		valSqls[(i * 7)] = e.AggregateID
		valSqls[(i*7)+1] = e.AggregateType
		valSqls[(i*7)+2] = e.EventID
		valSqls[(i*7)+3] = e.EventType
		valSqls[(i*7)+4] = e.Revision
		valSqls[(i*7)+5] = e.Time
		valSqls[(i*7)+6] = e.Data
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (aggregateID, aggregateType, eventID, eventType, revision, time, data) VALUES %s", s.eventTable, strings.Join(dataSqls, ","))
	if _, err := tx.Exec(query, valSqls...); err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (aggregateID, aggregateType, eventID, eventType, revision, time, data) VALUES %s", s.undispatchCommitedTable, strings.Join(dataSqls, ","))
	if _, err := tx.Exec(query, valSqls...); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetLastEvent get all event
func (s *Storage) GetLastEvent(id string) ([]goes.Event, error) {
	es := []goes.Event{}

	if err := s.db.Select(&es, s.getLastEventQuery, id); err != nil {
		return nil, err
	}

	if len(es) == 0 {
		return es, errors.NotFound("events not found")
	}

	return es, nil
}

// SaveSnapshot save snapshot
func (s *Storage) SaveSnapshot(snap *goes.Snapshot) error {
	if snap == nil {
		return nil
	}

	_, err := s.db.Exec(s.saveSnapshotQuery, snap.AggregateID, snap.AggregateType, snap.AggregateVersion, snap.Revision, snap.Time, snap.Data)

	return err
}

// GetSnapshot get snapshot with version
func (s *Storage) GetSnapshot(id string, version int) (*goes.Snapshot, error) {
	snap := goes.Snapshot{}
	if err := s.db.Get(&snap, s.getSnapshotQuery, id, version); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("snapshot not found")
		}

		return nil, err
	}

	return &snap, nil
}

// GetFromRevision get events from revision
func (s *Storage) GetFromRevision(id string, from int) ([]goes.Event, error) {
	es := []goes.Event{}

	if err := s.db.Select(&es, s.getFromRevisionQuery, id, from); err != nil {
		return nil, err
	}

	if len(es) == 0 {
		return es, errors.NotFound("events not found")
	}

	return es, nil
}

// GetUndispatchedEvent get undispatched event
func (s *Storage) GetUndispatchedEvent() ([]goes.Event, error) {
	es := []goes.Event{}

	if err := s.db.Select(&es, s.getUndispatchedEventQuery); err != nil {
		return nil, err
	}

	if len(es) == 0 {
		return es, errors.NotFound("events not found")
	}

	return es, nil
}

// MarkDispatchedEvent makr event dispatched
func (s *Storage) MarkDispatchedEvent(es []goes.Event) error {
	total := len(es)
	if total == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, e := range es {
		if _, err := tx.Exec(s.markDispatchedEventQuery, e.AggregateID, e.Revision); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
