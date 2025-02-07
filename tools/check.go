package tools

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/streamingfast/dstore"
	"github.com/streamingfast/substreams/block"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
	"github.com/streamingfast/substreams/store"
)

var checkCmd = &cobra.Command{
	Use:   "check <store_url>",
	Short: "checks the integrity of the kv files in a given store",
	Args:  cobra.ExactArgs(1),
	RunE:  checkE,
}

func init() {
	Cmd.AddCommand(checkCmd)
}

func checkE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	stateStore, _, err := newStore(args[0])
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	files, err := stateStore.ListSnapshotFiles(ctx)
	if err != nil {
		return fmt.Errorf("listing snapshots: %w", err)
	}

	var prevRange *block.Range
	for _, file := range files {
		if !file.Partial {
			continue
		}
		currentRange := block.NewRange(file.StartBlock, file.EndBlock)

		if prevRange == nil {
			prevRange = currentRange
			continue
		}

		if currentRange.StartBlock != prevRange.ExclusiveEndBlock {
			return fmt.Errorf("**hole found** between %d and %d", prevRange.ExclusiveEndBlock, currentRange.ExclusiveEndBlock)
		}

		prevRange = currentRange
	}

	return err
}

func newStore(storeURL string) (*store.FullKV, dstore.Store, error) {
	remoteStore, err := dstore.NewStore(storeURL, "", "", false)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create store from %s: %w", storeURL, err)
	}

	s, err := store.NewFullKV(
		"",
		0,
		"",
		pbsubstreams.Module_KindStore_UPDATE_POLICY_SET_IF_NOT_EXISTS,
		"",
		remoteStore, zap.NewNop(),
	)
	if err != nil {
		return nil, nil, err
	}

	return s, remoteStore, nil
}
