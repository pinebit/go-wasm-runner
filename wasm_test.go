package main

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/add.wasm
var addWasm []byte

func TestWasmService_StartStop(t *testing.T) {
	ctx := context.Background()
	wasm := NewWasmService()

	err := wasm.Start(ctx)
	assert.NoError(t, err)

	err = wasm.Stop(ctx)
	assert.NoError(t, err)
}

func TestWasmService_Run(t *testing.T) {
	ctx := context.Background()
	wasm := NewWasmService()

	err := wasm.Start(ctx)
	assert.NoError(t, err)

	output, err := wasm.Run(ctx, addWasm, 3, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint64(8), output)

	err = wasm.Stop(ctx)
	assert.NoError(t, err)
}
