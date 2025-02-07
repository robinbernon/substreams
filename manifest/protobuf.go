package manifest

import (
	"fmt"

	"github.com/jhump/protoreflect/desc/protoparse"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
	"github.com/streamingfast/substreams/pb/system"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func loadProtobufs(pkg *pbsubstreams.Package, manif *Manifest) error {
	// System protos
	systemFiles, err := readSystemProtobufs()
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, file := range systemFiles.File {
		pkg.ProtoFiles = append(pkg.ProtoFiles, file)
		seen[*file.Name] = true
	}

	var importPaths []string
	for _, imp := range manif.Protobuf.ImportPaths {
		importPaths = append(importPaths, manif.resolvePath(imp))
	}
	// User-specified protos
	parser := protoparse.Parser{
		ImportPaths:           importPaths,
		IncludeSourceCodeInfo: true,
	}

	for _, file := range manif.Protobuf.Files {
		if seen[file] {
			return fmt.Errorf("WARNING: proto file %s already exists in system protobufs, do not include in your manifest", file)
		}
	}
	customFiles, err := parser.ParseFiles(manif.Protobuf.Files...)
	if err != nil {
		return fmt.Errorf("error parsing proto files %q (import paths: %q): %w", manif.Protobuf.Files, importPaths, err)
	}
	for _, fd := range customFiles {
		pkg.ProtoFiles = append(pkg.ProtoFiles, fd.AsFileDescriptorProto())
	}

	return nil
}

func readSystemProtobufs() (*descriptorpb.FileDescriptorSet, error) {
	fds := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(system.ProtobufDescriptors, fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}
