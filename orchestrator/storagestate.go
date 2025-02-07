package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/abourget/llerrgroup"
	"github.com/streamingfast/substreams/store"
)

type StorageState struct {
	sync.Mutex
	Snapshots map[string]*Snapshots
}

func NewStorageState() *StorageState {
	return &StorageState{
		Snapshots: map[string]*Snapshots{},
	}
}

func (s *StorageState) String() string {
	var out []string
	for k, v := range s.Snapshots {
		out = append(out, fmt.Sprintf("store=%s (%s)", k, v))
	}
	return strings.Join(out, ", ")
}

func FetchStorageState(ctx context.Context, storeMaps *store.Map) (out *StorageState, err error) {
	out = NewStorageState()
	eg := llerrgroup.New(10)

	for name, storeObj := range storeMaps.All() {
		if eg.Stop() {
			break
		}
		storeName := name
		store := storeObj
		eg.Go(func() error {

			snapshots, err := listSnapshots(ctx, store)
			if err != nil {
				return err
			}
			out.Lock()
			out.Snapshots[storeName] = snapshots
			out.Unlock()
			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return nil, fmt.Errorf("running list snapshots: %w", err)
	}
	return
}
