package pipeline

import (
	"fmt"
	"github.com/streamingfast/substreams"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

func (p *Pipeline) sendSnapshots() error {
	snapshotModules := p.reqCtx.Request().InitialStoreSnapshotForModules
	if len(snapshotModules) == 0 {
		return nil

	}

	for _, modName := range snapshotModules {
		store, found := p.storeMap.Get(modName)
		if !found {
			return fmt.Errorf("store %q not found", modName)
		}

		send := func(count uint64, total uint64, deltas []*pbsubstreams.StoreDelta) {
			data := &pbsubstreams.InitialSnapshotData{
				ModuleName: modName,
				Deltas: &pbsubstreams.StoreDeltas{
					Deltas: deltas,
				},
				SentKeys:  count,
				TotalKeys: total,
			}
			p.respFunc(substreams.NewSnapshotData(data))
		}

		var count uint64
		total := store.Length()
		var accum []*pbsubstreams.StoreDelta

		store.Iter(func(k string, v []byte) error {
			count++
			accum = append(accum, &pbsubstreams.StoreDelta{
				Operation: pbsubstreams.StoreDelta_CREATE,
				Key:       k,
				NewValue:  v,
			})

			if count%100 == 0 {
				send(count, total, accum)
				accum = nil
			}
			return nil
		})

		if len(accum) != 0 {
			send(count, total, accum)
		}
	}

	p.respFunc(substreams.NewSnapshotComplete())

	return nil
}
