package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	consumer "github.com/mitooos/kinesis-consumer"
	storage "github.com/mitooos/kinesis-consumer/store/ddb"
)

// A myLogger provides a minimalistic logger satisfying the Logger interface.
type myLogger struct {
	logger *log.Logger
}

// Log logs the parameters to the stdlib logger. See log.Println.
func (l *myLogger) Log(args ...interface{}) {
	l.logger.Println(args...)
}

func main() {
	// Wrap myLogger around  apex logger
	logger := &myLogger{
		logger: log.New(os.Stdout, "consumer-example: ", log.LstdFlags),
	}


	stream := "test"
	app := "test_app"
	table := "kinesis-consumer-checkpoint"

	// New Kinesis and DynamoDB clients (if you need custom config)
	// client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
	config.WithRegion("us-west-2"),
	)
	if err != nil {
	// handle error
	log.Fatal(err)
	}

	myDdbClient := dynamodb.NewFromConfig(cfg)

	myKsis := kinesis.NewFromConfig(cfg)

	// ddb persitance
	ddb, err := storage.New(app, table, storage.WithDynamoClient(myDdbClient))
	if err != nil {
		logger.Log("checkpoint error: %v", err)
	}

	// expvar counter
	var counter = expvar.NewMap("counters")

	// consumer
	c, err := consumer.New(
		stream,
		consumer.WithStore(ddb),
		consumer.WithLogger(logger),
		consumer.WithCounter(counter),
		consumer.WithClient(myKsis),
		consumer.WithAggregation(true),
	)
	if err != nil {
		logger.Log("consumer error: %v", err)
	}

	// use cancel func to signal shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// trap SIGINT, wait to trigger shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancel()
	}()

	// scan stream
	err = c.Scan(ctx, func(r *consumer.Record) error {
		fmt.Println(string(r.Data))
		return nil // continue scanning
	})
	if err != nil {
		logger.Log("scan error: %v", err)
	}

	if err := ddb.Shutdown(); err != nil {
		logger.Log("storage shutdown error: %v", err)
	}
}
