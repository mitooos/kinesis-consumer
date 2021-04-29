package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	consumer "github.com/mitooos/kinesis-consumer"
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

	stream := "test"

	// client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
  	config.WithRegion("us-west-2"),
	)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	client := kinesis.NewFromConfig(cfg)

	logger := &myLogger{
		logger: log.New(os.Stdout, "consumer-example: ", log.LstdFlags),
	}


	// consumer
	c, err := consumer.New(
		stream,
		consumer.WithClient(client),
		consumer.WithLogger(logger),
		consumer.WithAggregation(true),
	)
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	// scan
	ctx := trap()
	err = c.Scan(ctx, func(r *consumer.Record) error {
		fmt.Println(string(r.Data))
		return nil // continue scanning
	})
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
}

func trap() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigs
		log.Printf("received %s", sig)
		cancel()
	}()

	return ctx
}
