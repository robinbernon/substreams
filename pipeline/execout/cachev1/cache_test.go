package cachev1

import (
	"testing"

	"github.com/streamingfast/logging"

	"github.com/streamingfast/substreams/block"
	"github.com/stretchr/testify/require"
)

func TestOutputCache_listContinuousCacheRanges(t *testing.T) {
	testCases := []struct {
		name           string
		fromBlock      uint64
		cachedRanges   block.Ranges
		expectedOutput string
	}{
		{
			name:           "sunny path",
			cachedRanges:   block.Ranges{{100, 200}, {200, 300}, {300, 400}},
			expectedOutput: "[100, 200),[200, 300),[300, 400)",
		},
		{
			name:           "sunny path with from",
			fromBlock:      99,
			cachedRanges:   block.Ranges{{100, 200}, {200, 300}, {300, 400}},
			expectedOutput: "[100, 200),[200, 300),[300, 400)",
		},
		{
			name:           "one",
			cachedRanges:   block.Ranges{{100, 200}},
			expectedOutput: "[100, 200)",
		},
		{
			name:           "none",
			cachedRanges:   nil,
			expectedOutput: "",
		},
		{
			name:           "split",
			cachedRanges:   block.Ranges{{100, 200}, {200, 300}, {400, 500}},
			expectedOutput: "[100, 200),[200, 300)",
		},
		{
			name:           "split and from",
			fromBlock:      300,
			cachedRanges:   block.Ranges{{100, 200}, {300, 400}},
			expectedOutput: "[300, 400)",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ranges := listContinuousCacheRanges(c.cachedRanges, c.fromBlock)
			result := ranges.String()
			require.Equal(t, c.expectedOutput, result)
		})
	}
}

func TestOutputCache_Delete(t *testing.T) {
	testCases := []struct {
		name         string
		kv           outputKV
		keysToDelete []string
		expectedKv   outputKV
	}{
		{
			name: "delete one block id from output cache",
			kv: map[string]*CacheItem{
				"1": {
					BlockNum: 1,
				},
				"2": {
					BlockNum: 2,
				},
			},
			keysToDelete: []string{"2"},
			expectedKv: map[string]*CacheItem{
				"1": {
					BlockNum: 1,
				},
			},
		},
		{
			name: "delete two block ids from output cache",
			kv: map[string]*CacheItem{
				"1": {
					BlockNum: 1,
				},
				"2": {
					BlockNum: 2,
				},
				"3": {
					BlockNum: 3,
				},
				"4": {
					BlockNum: 4,
				},
			},
			keysToDelete: []string{"1", "2"},
			expectedKv: map[string]*CacheItem{
				"3": {
					BlockNum: 3,
				},
				"4": {
					BlockNum: 4,
				},
			},
		},
	}
	var zlog, _ = logging.PackageLogger("test", "github.com/streamingfast/substreams/pipeline")
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			outputCache := NewOutputCache("module1", nil, 10, zlog)
			outputCache.kv = test.kv
			for _, key := range test.keysToDelete {
				outputCache.Delete(key)
			}
			require.Equal(t, test.expectedKv, outputCache.kv)
		})
	}
}
