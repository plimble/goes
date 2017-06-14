package dynamodb

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/plimble/errors"
	"github.com/plimble/goes"
)

// UndispatchedEvent entity
type UndispatchedEvent struct {
	Node          int    `json:"node" dynamodbav:"n"`
	KeyTime       string `json:"keyTime" dynamodbav:"keyTime"`
	AggregateID   string `json:"aggregateID" dynamodbav:"aggregateID"`
	AggregateType string `json:"aggregateType" dynamodbav:"aggregateType"`
	EventID       string `json:"eventID" dynamodbav:"eventID"`
	EventType     string `json:"eventType" dynamodbav:"eventType"`
	Revision      int    `json:"revision" dynamodbav:"revision"`
	Time          int64  `json:"time" dynamodbav:"time"`
	Data          []byte `json:"data" dynamodbav:"data"`
}

// Storage dynamodb storage
type Storage struct {
	eventTable              *string
	undispatchCommitedTable *string
	snapshotTable           *string
	db                      *dynamodb.DynamoDB
}

// New dynamodb storage
func New(sess *session.Session, eventTable, undispatchCommitedTable, snapshotTable string) *Storage {
	return &Storage{aws.String(eventTable), aws.String(undispatchCommitedTable), aws.String(snapshotTable), dynamodb.New(sess)}
}

// Save event
func (s *Storage) Save(es []goes.Event) error {
	total := len(es)
	if total == 0 {
		return nil
	}
	firstEvent := es[0]

	eventstoreList := make([]*dynamodb.WriteRequest, total)
	undisdpatchedList := make([]*dynamodb.WriteRequest, total)

	for i, e := range es {
		revisionStr := strconv.Itoa(e.Revision)
		item := map[string]*dynamodb.AttributeValue{
			"aggregateID":   {S: aws.String(e.AggregateID)},
			"aggregateType": {S: aws.String(e.AggregateType)},
			"eventID":       {S: aws.String(e.EventID)},
			"eventType":     {S: aws.String(e.EventType)},
			"revision":      {N: aws.String(revisionStr)},
			"time":          {N: aws.String(strconv.FormatInt(e.Time, 10))},
			"data":          {B: e.Data},
		}

		eventstoreList[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}

		item["node"] = &dynamodb.AttributeValue{N: aws.String("1")}
		item["keyTime"] = &dynamodb.AttributeValue{S: aws.String(e.AggregateID + e.AggregateType + revisionStr)}
		undisdpatchedList[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}
	}

	res, err := s.db.Query(&dynamodb.QueryInput{
		TableName:            s.eventTable,
		ProjectionExpression: aws.String("aggregateID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":aggregateID": {S: aws.String(firstEvent.AggregateID)},
			":revision":    {N: aws.String(strconv.Itoa(firstEvent.Revision))},
		},
		KeyConditionExpression: aws.String("aggregateID = :aggregateID AND revision >= :revision"),
	})

	if err != nil {
		return err
	}

	if len(res.Items) > 0 {
		return errors.New("event conflict")
	}

	if _, err := s.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			*s.eventTable:              eventstoreList,
			*s.undispatchCommitedTable: eventstoreList,
		},
	}); err != nil {
		return err
	}

	return nil
}

// SaveSnapshot save snapshot
func (s *Storage) SaveSnapshot(snap *goes.Snapshot) error {
	_, err := s.db.PutItem(&dynamodb.PutItemInput{
		TableName: s.snapshotTable,
		Item: map[string]*dynamodb.AttributeValue{
			"aggregateID":      {S: aws.String(snap.AggregateID)},
			"aggregateType":    {S: aws.String(snap.AggregateType)},
			"aggregateVersion": {N: aws.String(strconv.Itoa(snap.AggregateVersion))},
			"revision":         {N: aws.String(strconv.Itoa(snap.Revision))},
			"time":             {N: aws.String(strconv.FormatInt(snap.Time, 10))},
			"data":             {B: snap.Data},
		},
	})

	return err
}

// GetLastEvent get all event
func (s *Storage) GetLastEvent(id string) ([]goes.Event, error) {
	res, err := s.db.Query(&dynamodb.QueryInput{
		TableName: s.eventTable,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":aggregateID": {S: aws.String(id)},
		},
		ScanIndexForward:       aws.Bool(true),
		KeyConditionExpression: aws.String("aggregateID = :aggregateID"),
	})

	if err != nil {
		return nil, err
	}

	events := make([]goes.Event, len(res.Items))
	if len(events) == 0 {
		return events, errors.NotFound("event not found")
	}

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &events)

	return events, err
}

// GetSnapshot get snapshot
func (s *Storage) GetSnapshot(id string, version int) (*goes.Snapshot, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: s.snapshotTable,
		Key: map[string]*dynamodb.AttributeValue{
			"aggregateID":      {S: aws.String(id)},
			"aggregateVersion": {N: aws.String(strconv.Itoa(version))},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, errors.NotFound("snapshot not found")
	}

	snap := &goes.Snapshot{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &snap)

	return snap, nil
}

// GetFromRevision get from revision
func (s *Storage) GetFromRevision(id string, from int) ([]goes.Event, error) {
	res, err := s.db.Query(&dynamodb.QueryInput{
		TableName: s.eventTable,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":aggregateID": {S: aws.String(id)},
			":revision":    {N: aws.String(strconv.Itoa(from))},
		},
		ScanIndexForward:       aws.Bool(true),
		KeyConditionExpression: aws.String("aggregateID = :aggregateID AND revision > :revision"),
	})

	if err != nil {
		return nil, err
	}

	events := make([]goes.Event, len(res.Items))
	if len(events) == 0 {
		return events, errors.NotFound("event not found")
	}

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &events)

	return events, err
}

// GetUndispatchedEvent get all undispatched event
func (s *Storage) GetUndispatchedEvent() ([]goes.Event, error) {
	res, err := s.db.Query(&dynamodb.QueryInput{
		TableName: s.undispatchCommitedTable,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":node": {N: aws.String("1")},
		},
		ScanIndexForward:       aws.Bool(true),
		KeyConditionExpression: aws.String("node = :node"),
	})

	if err != nil {
		return nil, err
	}

	events := make([]goes.Event, len(res.Items))
	if len(events) == 0 {
		return nil, errors.NotFound("event not found")
	}

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &events)

	return events, err
}

// MarkDispatchedEvent delete dispatched event
func (s *Storage) MarkDispatchedEvent(es []goes.Event) error {
	total := len(es)
	dispatchedList := make([]*dynamodb.WriteRequest, total)

	for i, e := range es {
		dispatchedList[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"node":    {N: aws.String("1")},
					"keyTime": {S: aws.String(e.AggregateID + e.AggregateType + strconv.Itoa(e.Revision))},
				},
			},
		}
	}

	if _, err := s.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			*s.undispatchCommitedTable: dispatchedList,
		},
	}); err != nil {
		return err
	}

	return nil
}
