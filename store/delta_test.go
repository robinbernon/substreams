package store

import (
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplyDelta(t *testing.T) {
	tests := []struct {
		name       string
		deltas     []*pbsubstreams.StoreDelta
		expectedKV map[string][]byte
	}{
		{
			name: "creates",
			deltas: []*pbsubstreams.StoreDelta{
				{
					Operation: pbsubstreams.StoreDelta_CREATE,
					Key:       "k1",
					NewValue:  []byte("v1"),
				},
				{
					Operation: pbsubstreams.StoreDelta_CREATE,
					Key:       "k2",
					NewValue:  []byte("v2"),
				},
			},
			expectedKV: map[string][]byte{
				"k1": []byte("v1"),
				"k2": []byte("v2"),
			},
		},
		{
			name: "update",
			deltas: []*pbsubstreams.StoreDelta{
				{
					Operation: pbsubstreams.StoreDelta_CREATE,
					Key:       "k1",
					NewValue:  []byte("v1"),
				},
				{
					Operation: pbsubstreams.StoreDelta_UPDATE,
					Key:       "k1",
					OldValue:  []byte("v1"),
					NewValue:  []byte("v2"),
				},
			},
			expectedKV: map[string][]byte{
				"k1": []byte("v2"),
			},
		},
		{
			name: "delete",
			deltas: []*pbsubstreams.StoreDelta{
				{
					Operation: pbsubstreams.StoreDelta_CREATE,
					Key:       "k1",
					NewValue:  []byte("v1"),
				},
				{
					Operation: pbsubstreams.StoreDelta_CREATE,
					Key:       "k2",
					NewValue:  []byte("v2"),
				},
				{
					Operation: pbsubstreams.StoreDelta_DELETE,
					Key:       "k1",
					OldValue:  []byte("v1"),
				},
			},
			expectedKV: map[string][]byte{
				"k2": []byte("v2"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &BaseStore{
				kv: make(map[string][]byte),
			}
			for _, delta := range test.deltas {
				s.ApplyDelta(delta)
			}
			assert.Equal(t, test.expectedKV, s.kv)
		})
	}
}

func Test_ApplyDeltasReverse(t *testing.T) {
	testCases := []struct {
		name       string
		store      *BaseStore
		expectedKV map[string][]byte
	}{
		{
			name: "reverse one delta",
			store: &BaseStore{
				deltas: []*pbsubstreams.StoreDelta{
					{
						Operation: pbsubstreams.StoreDelta_CREATE,
						Key:       "key_1",
						NewValue:  []byte{99},
					},
				},
				kv: map[string][]byte{
					"key_1": {99},
				},
			},
			expectedKV: map[string][]byte{},
		},
		{
			name: "reverse a delta when multiple deltas were applied",
			store: &BaseStore{
				deltas: []*pbsubstreams.StoreDelta{
					{
						Operation: pbsubstreams.StoreDelta_UPDATE,
						Key:       "key_1",
						OldValue:  []byte{99},
						NewValue:  []byte{100},
					},
				},
				kv: map[string][]byte{
					"key_1": {100},
				},
			},
			expectedKV: map[string][]byte{
				"key_1": {99},
			},
		},
		{
			name: "reverse multiple deltas",
			store: &BaseStore{
				deltas: []*pbsubstreams.StoreDelta{
					{
						Operation: pbsubstreams.StoreDelta_DELETE,
						Key:       "key_1",
						OldValue:  []byte{100},
					},
					{
						Operation: pbsubstreams.StoreDelta_UPDATE,
						Key:       "key_2",
						OldValue:  []byte{100},
						NewValue:  []byte{150},
					},
				},
				kv: map[string][]byte{
					"key_2": {150},
				},
			},
			expectedKV: map[string][]byte{
				"key_1": {100},
				"key_2": {100},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.store.ApplyDeltasReverse(test.store.deltas)
			require.Equal(t, test.expectedKV, test.store.kv)
		})
	}
}
