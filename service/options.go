package service

import (
	"github.com/streamingfast/substreams/pipeline"
	"github.com/streamingfast/substreams/wasm"
)

type Option func(*Service)

func WithWASMExtension(ext wasm.WASMExtensioner) Option {
	return func(s *Service) {
		s.wasmExtensions = append(s.wasmExtensions, ext)
	}
}

func WithPipelineOptions(f pipeline.PipelineOptioner) Option {
	return func(s *Service) {
		s.pipelineOptions = append(s.pipelineOptions, f)
	}
}

func WithPartialMode() Option {
	return func(s *Service) {
		s.partialModeEnabled = true
	}
}

func WithStoresSaveInterval(block uint64) Option {
	return func(s *Service) {
		s.storesSaveInterval = block
	}
}

func WithOutCacheSaveInterval(block uint64) Option {
	return func(s *Service) {
		s.outputCacheSaveBlockInterval = block
	}
}
