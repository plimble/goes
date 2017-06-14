CREATE TABLE IF NOT EXISTS eventstore (
  aggregateID varchar(64) NOT NULL DEFAULT '',
  revision int(11) NOT NULL,
  aggregateType varchar(64) NOT NULL DEFAULT '',
  eventID varchar(64) NOT NULL DEFAULT '',
  eventType varchar(64) NOT NULL DEFAULT '',
  time bigint(20) NOT NULL,
  data blob NOT NULL,
  PRIMARY KEY (aggregateID,revision)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS eventstore_log (
	aggregateID varchar(64) NOT NULL DEFAULT '',
	revision int(11) NOT NULL,
	aggregateType varchar(64) NOT NULL DEFAULT '',
	eventID varchar(64) NOT NULL DEFAULT '',
	eventType varchar(64) NOT NULL,
	time bigint(20) NOT NULL,
	data blob NOT NULL,
	PRIMARY KEY (aggregateID,revision)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS eventstore_snapshot (
	aggregateID varchar(64) NOT NULL DEFAULT '',
	aggregateVersion int(11) NOT NULL,
	revision int(11) NOT NULL,
	aggregateType varchar(64) NOT NULL DEFAULT '',
	time bigint(20) NOT NULL,
	data blob NOT NULL,
	PRIMARY KEY (aggregateID,aggregateVersion)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;