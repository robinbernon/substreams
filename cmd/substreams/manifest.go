package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/streamingfast/substreams/manifest"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

// var manifestCmd = &cobra.Command{
// 	Use:          "manifest",
// 	SilenceUsage: true,
// }
var infoCmd = &cobra.Command{
	Use:          "info <manifest_file>",
	Short:        "Display package modules and docs",
	RunE:         runInfo,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	manifestPath := args[0]
	manifestReader := manifest.NewReader(manifestPath)
	pkg, err := manifestReader.Read()
	if err != nil {
		return fmt.Errorf("read manifest %q: %w", manifestPath, err)
	}

	graph, err := manifest.NewModuleGraph(pkg.Modules.Modules)
	if err != nil {
		return fmt.Errorf("creating module graph: %w", err)
	}

	fmt.Println("Package name:", pkg.PackageMeta[0].Name)
	fmt.Println("Version:", pkg.PackageMeta[0].Version)
	if doc := pkg.PackageMeta[0].Doc; doc != "" {
		fmt.Println("Doc: " + strings.Replace(doc, "\n", "\n  ", -1))
	}

	hashes := manifest.NewModuleHashes()

	fmt.Println("Modules:")
	fmt.Println("----")
	for modIdx, module := range pkg.Modules.Modules {
		fmt.Println("Filename:", module.Name)
		fmt.Println("Initial block:", module.InitialBlock)
		kind := module.GetKind()
		switch v := kind.(type) {
		case *pbsubstreams.Module_KindMap_:
			fmt.Println("Kind: map")
			fmt.Println("Output Type:", v.KindMap.OutputType)
		case *pbsubstreams.Module_KindStore_:
			fmt.Println("Kind: store")
			fmt.Println("Value Type:", v.KindStore.ValueType)
			fmt.Println("Update Policy:", v.KindStore.UpdatePolicy)
		default:
			fmt.Println("Kind: Unknown")
		}

		hashes.HashModule(pkg.Modules, module, graph)

		fmt.Println("Hash:", hashes.Get(module.Name))
		moduleMeta := pkg.ModuleMeta[modIdx]
		if moduleMeta != nil && moduleMeta.Doc != "" {
			fmt.Println("Doc: " + strings.Replace(moduleMeta.Doc, "\n", "\n  ", -1))
		}
		fmt.Println("")
	}

	return nil
}
