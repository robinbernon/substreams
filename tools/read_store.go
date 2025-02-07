package tools

import (
	"github.com/spf13/cobra"
	"github.com/streamingfast/substreams/manifest"
)

var storeCmd = &cobra.Command{
	Use:   "store get <manifest_path> <module_name> <block_id> <key>",
	Short: "BaseStore files in a common archive format",
	Long:  `BaseStore files in a common archive format`,
	RunE:  storeGetE,
	Args:  cobra.ExactArgs(4),
}

func init() {
	Cmd.AddCommand(storeCmd)
}

func storeGetE(cmd *cobra.Command, args []string) error {
	pkg, err := manifest.NewReader(args[0]).Read()
	if err != nil {
		return err
	}

	for _, m := range pkg.Modules.Modules {
		if m.Name == args[1] {

		}
	}

	return nil
}
