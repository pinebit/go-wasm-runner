package main

import (
	"context"
	"errors"

	"github.com/pinebit/go-boot/boot"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmService interface {
	boot.Service

	Run(ctx context.Context, wasm []byte, a, b uint64) (uint64, error)
}

type wasmService struct {
	runtime wazero.Runtime
}

func NewWasmService() WasmService {
	return &wasmService{}
}

func (s *wasmService) Start(ctx context.Context) error {
	s.runtime = wazero.NewRuntime(ctx)
	_, err := wasi_snapshot_preview1.Instantiate(ctx, s.runtime)
	return err
}

func (s *wasmService) Stop(ctx context.Context) error {
	return s.runtime.Close(ctx)
}

func (s *wasmService) Run(ctx context.Context, wasm []byte, a, b uint64) (uint64, error) {
	mod, err := s.runtime.Instantiate(ctx, wasm)
	if err != nil {
		return 0, err
	}

	addFunc := mod.ExportedFunction("add")
	if addFunc == nil {
		return 0, errors.New("no add function is exported")
	}

	r, err := addFunc.Call(ctx, a, b)
	if err != nil {
		return 0, err
	}
	if len(r) == 0 {
		return 0, errors.New("no data returned from run function")
	}
	return r[0], nil
}
