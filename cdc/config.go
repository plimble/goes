package main

import (
	"os"

	"github.com/plimble/validator"
	"gopkg.in/plimble/goconf.v1"
)

// Config model
type Config struct {
	PullMsInternal    int    `envconfig:"PULL_MS_INTERVAL"`
	UndispatchedTable string `envconfig:"UNDISPATCHED_TABLE"`
	MysqlDataSource   string `envconfig:"MYSQL_DATASOURCE"`

	Nats struct {
		URL       string `envconfig:"URL"`
		ClusterID string `envconfig:"CLUSTER_ID"`
	}

	Aws struct {
		KeyID            string
		Secret           string
		DynamodbEndpoint string `envconfig:"DYNAMODB_ENDPOINT"`
		Region           string `envconfig:"REGION"`
	}
}

var config *Config

func defaultConfig() *Config {
	c := &Config{}
	c.PullMsInternal = 100

	c.Aws.KeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	c.Aws.Secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	c.Aws.Region = "ap-southeast-1"

	return c
}

func validteConfig(c *Config) error {
	v := validator.New()
	v.RequiredInt(c.PullMsInternal, "PullMsInternal")
	v.RequiredString(c.MysqlDataSource, "MysqlDataSource")
	v.RequiredString(c.UndispatchedTable, "UndispatchedTable")
	v.RequiredString(c.Nats.ClusterID, "Nats.ClusterID")
	v.RequiredString(c.Nats.URL, "Nats.URL")

	v.RequiredString(c.Aws.KeyID, "Aws.KeyId")
	v.RequiredString(c.Aws.Secret, "Aws.Secret")
	v.RequiredString(c.Aws.Region, "Aws.Region")

	return v.GetError()
}

// Get config
func Get() *Config {
	if config == nil {
		config = defaultConfig()
		goconf.Parse(config,
			goconf.WithEnv("youex"),
		)
		if err := validteConfig(config); err != nil {
			panic(err)
		}
	}

	return config
}
