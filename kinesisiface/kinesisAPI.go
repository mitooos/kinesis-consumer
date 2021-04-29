package kinesisiface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type KinesisAPI interface {
	GetRecords(context.Context, *kinesis.GetRecordsInput, ...func(*kinesis.Options)) (*kinesis.GetRecordsOutput, error)
	GetShardIterator(context.Context, *kinesis.GetShardIteratorInput, ...func(*kinesis.Options)) (*kinesis.GetShardIteratorOutput, error)
	ListShards(context.Context, *kinesis.ListShardsInput, ...func(*kinesis.Options)) (*kinesis.ListShardsOutput, error)
}
