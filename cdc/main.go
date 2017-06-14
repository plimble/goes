package main

import (
	"log"
	"os"
	"os/signal"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/go-nats-streaming"
	"github.com/plimble/errors"
	"github.com/plimble/goes"
	"github.com/plimble/goes/storage/mysql"
)

// Context context
type Context struct {
	es []*goes.Event
	e  goes.Event
}

func main() {

	conf := Get()
	db, err := sqlx.Connect("mysql", conf.MysqlDataSource)
	if err != nil {
		log.Fatal(err)
	}

	storage := mysql.New(db, "", conf.UndispatchedTable, "")

	conn, err := stan.Connect(conf.Nats.ClusterID, "goes-cdc")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	var es []goes.Event
	var e goes.Event
	var data []byte

	pulldelay := time.Duration(conf.PullMsInternal) * time.Millisecond

	go func() {
		for {
			time.Sleep(pulldelay)
			es, err = storage.GetUndispatchedEvent()
			if err != nil {
				if !errors.IsNotFound(err) {
					log.Println(err)
				}
				continue
			}

			if err == nil {
				for _, e = range es {
					data, _ = e.MarshalMsg(nil)
					if err = conn.Publish("eventstore."+e.AggregateType, data); err != nil {
						log.Println(err)
						continue
					} else {
						log.Printf("Published %s %s %s %d", e.EventType, e.AggregateID, e.AggregateType, e.Revision)
					}
				}
			}

			if err == nil {
				if err = storage.MarkDispatchedEvent(es); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

	s := <-c
	log.Println("Closed:", s)
	conn.Close()
}
