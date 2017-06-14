package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DeleteTables delete all tables
func (s *Storage) DeleteTables() {
	s.deleteTable(s.eventTable)
	s.deleteTable(s.snapshotTable)
	s.deleteTable(s.undispatchCommitedTable)
}

func (s *Storage) deleteTable(name *string) {
	_, err := s.db.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: name,
	})

	if err != nil {
		aerr := err.(awserr.Error)
		if aerr.Code() != dynamodb.ErrCodeResourceNotFoundException {
			panic(err)
		}
	}
}

// CreateTables all tables
func (s *Storage) CreateTables() {
	_, err := s.db.CreateTable(&dynamodb.CreateTableInput{
		TableName: s.eventTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("aggregateID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("revision"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("aggregateID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("revision"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		aerr := err.(awserr.Error)
		if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
			panic(err)
		}
	}

	_, err = s.db.CreateTable(&dynamodb.CreateTableInput{
		TableName: s.snapshotTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("aggregateID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("aggregateVersion"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("aggregateID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("aggregateVersion"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		aerr := err.(awserr.Error)
		if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
			panic(err)
		}
	}

	_, err = s.db.CreateTable(&dynamodb.CreateTableInput{
		TableName: s.undispatchCommitedTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("node"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("keyTime"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("node"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("keyTime"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

}
