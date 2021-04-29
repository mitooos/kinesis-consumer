package consumer

import (
	"fmt"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"

	"github.com/mitooos/kinesis-consumer/kinesisiface"
)

// listShards pulls a list of shard IDs from the kinesis api
func listShards(ksis kinesisiface.KinesisAPI, streamName string) ([]types.Shard, error) {
	var ss []types.Shard
	var listShardsInput = &kinesis.ListShardsInput{
		StreamName: aws.String(streamName),
	}

	ctx := context.Background()

	for {
		resp, err := ksis.ListShards(ctx, listShardsInput)
		if err != nil {
			return nil, fmt.Errorf("ListShards error: %w", err)
		}
		ss = append(ss, resp.Shards...)

		if resp.NextToken == nil {
			return ss, nil
		}

		listShardsInput = &kinesis.ListShardsInput{
			NextToken: resp.NextToken,
		}
	}
}
