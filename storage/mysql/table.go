package mysql

import (
	"fmt"
)

// FlushAll tables
func (s *Storage) FlushAll() {
	s.db.Exec(fmt.Sprintf("DELETE FROM %s", s.eventTable))
	s.db.Exec(fmt.Sprintf("DELETE FROM %s", s.undispatchCommitedTable))
	s.db.Exec(fmt.Sprintf("DELETE FROM %s", s.snapshotTable))
}

// CreateTables create table
func (s *Storage) CreateTables() {
	schema := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	aggregateID varchar(64) NOT NULL DEFAULT '',
	revision int(11) NOT NULL,
	aggregateType varchar(64) NOT NULL DEFAULT '',
	eventID varchar(64) NOT NULL DEFAULT '',
	eventType varchar(64) NOT NULL DEFAULT '',
	time bigint(20) NOT NULL,
	data blob NOT NULL,
	PRIMARY KEY (aggregateID,revision)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, s.eventTable)
	s.db.MustExec(schema)

	schema = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	aggregateID varchar(64) NOT NULL DEFAULT '',
	revision int(11) NOT NULL,
	aggregateType varchar(64) NOT NULL DEFAULT '',
	eventID varchar(64) NOT NULL DEFAULT '',
	eventType varchar(64) NOT NULL,
	time bigint(20) NOT NULL,
	data blob NOT NULL,
	PRIMARY KEY (aggregateID,revision)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, s.undispatchCommitedTable)
	s.db.MustExec(schema)

	schema = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	aggregateID varchar(64) NOT NULL DEFAULT '',
	aggregateVersion int(11) NOT NULL,
	revision int(11) NOT NULL,
	aggregateType varchar(64) NOT NULL DEFAULT '',
	time bigint(20) NOT NULL,
	data blob NOT NULL,
	PRIMARY KEY (aggregateID,aggregateVersion)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, s.snapshotTable)
	s.db.MustExec(schema)
}

// DropTables drop table
func (s *Storage) DropTables() {
	schema := fmt.Sprintf(`
	DROP TABLE IF EXISTS %s, %s, %s;
	`, s.eventTable, s.undispatchCommitedTable, s.snapshotTable)
	s.db.MustExec(schema)
}
