package storegateway

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cortex_tsdb "github.com/grafana/mimir/pkg/storage/tsdb"
)

func TestChunkBytesPool_Get(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()
	p, err := newChunkBytesPool(cortex_tsdb.ChunkPoolDefaultMinBucketSize, cortex_tsdb.ChunkPoolDefaultMaxBucketSize, 0, reg)
	require.NoError(t, err)

	_, err = p.Get(cortex_tsdb.EstimatedMaxChunkSize - 1)
	require.NoError(t, err)

	_, err = p.Get(cortex_tsdb.EstimatedMaxChunkSize + 1)
	require.NoError(t, err)

	assert.NoError(t, testutil.GatherAndCompare(reg, bytes.NewBufferString(fmt.Sprintf(`
		# HELP cortex_bucket_store_chunk_pool_requested_bytes_total Total bytes requested to chunk bytes pool.
		# TYPE cortex_bucket_store_chunk_pool_requested_bytes_total counter
		cortex_bucket_store_chunk_pool_requested_bytes_total %d

		# HELP cortex_bucket_store_chunk_pool_returned_bytes_total Total bytes returned by the chunk bytes pool.
		# TYPE cortex_bucket_store_chunk_pool_returned_bytes_total counter
		cortex_bucket_store_chunk_pool_returned_bytes_total %d
	`, cortex_tsdb.EstimatedMaxChunkSize*2, cortex_tsdb.EstimatedMaxChunkSize*3))))
}
