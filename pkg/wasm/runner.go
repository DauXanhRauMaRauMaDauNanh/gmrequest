package wasm

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/profile"
)

// Runner is a headless WebAssembly execution engine for computing anti-bot PoW & challenge algorithms
type Runner struct {
	runtime wazero.Runtime
	profile *profile.DeviceProfile
	mu      sync.Mutex
}

// NewRunner creates a new WASM execution context with mocked browser environment
func NewRunner(ctx context.Context, p *profile.DeviceProfile) (*Runner, error) {
	if p == nil {
		p = profile.Random()
	}

	config := wazero.NewRuntimeConfigCompiler() // High performance JIT compiler
	runtime := wazero.NewRuntimeWithConfig(ctx, config)

	// Instantiate WASI preview 1 support (standard for WebAssembly binaries)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		return nil, fmt.Errorf("failed to instantiate wasi: %w", err)
	}

	runner := &Runner{
		runtime: runtime,
		profile: p,
	}

	if err := runner.registerBrowserHostEnv(ctx); err != nil {
		return nil, fmt.Errorf("failed to register browser mock env: %w", err)
	}

	return runner, nil
}

// registerBrowserHostEnv mocks common browser APIs imported by challenge WASM modules
func (r *Runner) registerBrowserHostEnv(ctx context.Context) error {
	startTime := time.Now()

	_, err := r.runtime.NewHostModuleBuilder("env").
		// performance.now() mock
		NewFunctionBuilder().
		WithFunc(func() float64 {
			return float64(time.Since(startTime).Microseconds()) / 1000.0
		}).
		Export("performance_now").

		// Date.now() mock
		NewFunctionBuilder().
		WithFunc(func() int64 {
			return time.Now().UnixMilli()
		}).
		Export("date_now").

		// window.crypto.getRandomValues(buf) mock
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, ptr uint32, length uint32) {
			buf := make([]byte, length)
			_, _ = rand.Read(buf)
			mod.Memory().Write(ptr, buf)
		}).
		Export("crypto_get_random_values").

		// screen.width / screen.height mock aligned with DeviceProfile
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return uint32(r.profile.ScreenWidth)
		}).
		Export("screen_width").

		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return uint32(r.profile.ScreenHeight)
		}).
		Export("screen_height").

		// navigator.hardwareConcurrency mock
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return uint32(r.profile.HardwareConcurrency)
		}).
		Export("hardware_concurrency").

		// navigator.deviceMemory mock
		NewFunctionBuilder().
		WithFunc(func() uint32 {
			return uint32(r.profile.DeviceMemory)
		}).
		Export("device_memory").

		Instantiate(ctx)

	return err
}

// RunModule loads and runs a WASM challenge module, calling an exported challenge solver function
func (r *Runner) RunModule(ctx context.Context, wasmBytes []byte, funcName string, params ...uint64) ([]uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	compiled, err := r.runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm module: %w", err)
	}
	defer compiled.Close(ctx)

	mod, err := r.runtime.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm module: %w", err)
	}
	defer mod.Close(ctx)

	targetFunc := mod.ExportedFunction(funcName)
	if targetFunc == nil {
		return nil, fmt.Errorf("function %q not exported by wasm module", funcName)
	}

	return targetFunc.Call(ctx, params...)
}

// Close releases the WASM runtime resources
func (r *Runner) Close(ctx context.Context) error {
	return r.runtime.Close(ctx)
}
