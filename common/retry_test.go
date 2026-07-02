package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultRetryTimesAllowsChannelFallback(t *testing.T) {
	require.Equal(t, 2, RetryTimes)
}
