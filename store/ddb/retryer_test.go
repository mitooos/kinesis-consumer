package ddb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDefaultRetyer(t *testing.T) {
	retryableError := &types.ProvisionedThroughputExceededException{}
	// retryer is not nil and should returns according to what error is passed in.
	q := &DefaultRetryer{}
	if q.ShouldRetry(retryableError) != true {
		t.Errorf("expected ShouldRetry returns %v. got %v", false, q.ShouldRetry(retryableError))
	}

	nonRetryableError := &types.PointInTimeRecoveryUnavailableException{}
	shouldRetry := q.ShouldRetry(nonRetryableError)
	if shouldRetry != false {
		t.Errorf("expected ShouldRetry returns %v. got %v", true, shouldRetry)
	}
}
