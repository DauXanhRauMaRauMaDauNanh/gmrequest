package gmrequest

import (
	"context"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/engine"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/profile"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/wasm"
)

// Re-export common types for ease of use
type Client = engine.Client
type ClientOptions = engine.ClientOptions
type Response = engine.Response
type DeviceProfile = profile.DeviceProfile
type Platform = profile.Platform
type WASMRunner = wasm.Runner

const (
	PlatformWindows = profile.PlatformWindows
	PlatformMacOS   = profile.PlatformMacOS
	PlatformLinux   = profile.PlatformLinux
	PlatformAndroid = profile.PlatformAndroid
	PlatformIOS     = profile.PlatformIOS
)

// NewClient creates a new high-stealth HTTP/TLS/H2 client
func NewClient(opts ...ClientOptions) (*Client, error) {
	return engine.NewClient(opts...)
}

// NewClientWithSeed creates a client with a deterministic device profile seed (like CloakBrowser's --fingerprint=seed)
func NewClientWithSeed(seed int64, platform Platform, proxyURL ...string) (*Client, error) {
	p := profile.GenerateFromSeed(seed, platform)
	opts := ClientOptions{
		Profile: p,
	}
	if len(proxyURL) > 0 {
		opts.ProxyURL = proxyURL[0]
	}
	return engine.NewClient(opts)
}

// NewWASMRunner creates a headless WebAssembly runtime for solving PoW/challenge modules in Go
func NewWASMRunner(ctx context.Context, p *DeviceProfile) (*WASMRunner, error) {
	return wasm.NewRunner(ctx, p)
}
