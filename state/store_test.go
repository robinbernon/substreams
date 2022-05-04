package state

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/streamingfast/dstore"
)

type TestStore struct {
	*dstore.MockStore

	WriteStateFunc        func(ctx context.Context, content []byte, blockNum uint64) error
	WritePartialStateFunc func(ctx context.Context, content []byte, startBlockNum, endBlockNum uint64) error
}

func (io *TestStore) WritePartialState(ctx context.Context, content []byte, startBlockNum, endBlockNum uint64) error {
	if io.WritePartialStateFunc != nil {
		return io.WritePartialStateFunc(ctx, content, startBlockNum, endBlockNum)
	}
	return nil
}

func (io *TestStore) WriteState(ctx context.Context, content []byte, blockNum uint64) error {
	if io.WriteStateFunc != nil {
		return io.WriteStateFunc(ctx, content, blockNum)
	}
	return nil
}

func TestStateFileName(t *testing.T) {
	prefix := StateFilePrefix("test", 10000)
	require.Equal(t, "test-0000010000", prefix)

	stateFileName := StateFileName("test", 100, 10000)
	require.Equal(t, "test-0000010000-0000000100.kv", stateFileName)

	partialFileName := PartialFileName("test", 10000, 20000)
	require.Equal(t, "test-0000020000-0000010000.partial", partialFileName)
}
