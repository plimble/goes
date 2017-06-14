package example

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/plimble/goes"
	encoder "github.com/plimble/goes/encoder/json"
	"github.com/plimble/goes/storage/dynamodb"
	"github.com/stretchr/testify/require"
)

type NameChanged struct {
	Name string
}

type AgeChanged struct {
	Age int
}

type Account struct {
	*goes.AggregateBase
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (a *Account) ApplyEvent(e goes.Event) error {
	switch e.EventType {
	case "NameChanged":
		event := &NameChanged{}
		if err := a.DecodeEvent(e.Data, event); err != nil {
			return err
		}

		a.Name = event.Name
	case "AgeChanged":
		event := &AgeChanged{}
		if err := a.DecodeEvent(e.Data, event); err != nil {
			return err
		}

		a.Age = event.Age
	}

	return nil
}

func (a *Account) ChangeName(name string) {
	a.AddEvent("NameChanged", &NameChanged{name})
}

func (a *Account) ChangeAge(age int) {
	a.AddEvent("AgeChanged", &AgeChanged{age})
}

func TestEventStore(t *testing.T) {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewEnvCredentials())
	config.WithRegion("ap-southeast-1")
	config.WithEndpoint("http://localhost:8000")
	sess, err := session.NewSession(config)
	require.NoError(t, err)

	storage := dynamodb.New(sess, "eventstore", "undispatch", "snapshot")
	storage.DeleteTables()
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
