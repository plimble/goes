package example

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/plimble/goes"
	encoder "github.com/plimble/goes/encoder/json"
	"github.com/plimble/goes/storage/mysql"
	"github.com/stretchr/testify/require"
)

func TestMysqlEventStore(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root@tcp(127.0.0.1:3306)/goes")
	require.NoError(t, err)

	storage := mysql.New(db, "eventstore", "eventstore_log", "eventstore_snapshot")
	storage.DropTables()
	storage.CreateTables()

	es := goes.NewEventStore(storage, encoder.New(), goes.NewUUIDV4())

	id := "1111"
	a := &Account{}
	a.AggregateBase = es.NewAggregateBase(id, "Account", 1, a.ApplyEvent)

	a.ChangeName("john doe")
	a.ChangeAge(10)

	require.Equal(t, "john doe", a.Name)
	require.Equal(t, 10, a.Age)
	require.Equal(t, 2, a.Revision)
	require.Equal(t, 1, a.Version)

	err = es.Commit(a)
	require.NoError(t, err)

	a = &Account{}
	a.AggregateBase = es.NewAggregateBase(id, "Account", 1, a.ApplyEvent)
	err = es.GetStream(a)
	require.NoError(t, err)

	require.Equal(t, "john doe", a.Name)
	require.Equal(t, 10, a.Age)
	require.Equal(t, 2, a.Revision)
	require.Equal(t, 1, a.Version)

	a.ChangeName("tester")
	require.Equal(t, "tester", a.Name)

	err = es.SaveSnapshot(a)
	require.NoError(t, err)

	a.ChangeAge(20)

	err = es.Commit(a)
	require.NoError(t, err)

	a = &Account{}
	a.AggregateBase = es.NewAggregateBase(id, "Account", 1, a.ApplyEvent)
	err = es.GetStreamWithSnapshot(a)
	require.NoError(t, err)

	require.Equal(t, "tester", a.Name)
	require.Equal(t, 20, a.Age)
	require.Equal(t, 4, a.Revision)
	require.Equal(t, 1, a.Version)
}
