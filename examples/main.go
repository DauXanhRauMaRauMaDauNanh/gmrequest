package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("       gmrequest: Full Stealth HTTP/2 & TLS Fingerprint Engine  ")
	fmt.Println("================================================================")

	// 1. Tạo Client với Seed cố định (Chuẩn CloakBrowser/RayoBrowse device coherence)
	seed := int64(987654321)
	client, err := gmrequest.NewClientWithSeed(seed, gmrequest.PlatformWindows)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	fmt.Printf("[+] Generated Coherent Profile:\n")
	fmt.Printf("    - Browser: %s v%s\n", client.Profile.BrowserName, client.Profile.BrowserVer)
	fmt.Printf("    - OS Platform: %s (v%s)\n", client.Profile.Platform, client.Profile.PlatformVersion)
	fmt.Printf("    - User-Agent: %s\n", client.Profile.UserAgent)
	fmt.Printf("    - Sec-CH-UA: %s\n", client.Profile.SecChUa)
	fmt.Printf("    - Screen Geometry: %dx%d (Available: %dx%d)\n",
		client.Profile.ScreenWidth, client.Profile.ScreenHeight,
		client.Profile.AvailWidth, client.Profile.AvailHeight)
	fmt.Printf("    - Hardware Concurrency: %d Cores | Memory: %d GB\n\n",
		client.Profile.HardwareConcurrency, client.Profile.DeviceMemory)

	// 2. Test request kiểm tra Headers gửi đi
	fmt.Println("[+] Sending Stealth Request to check header order & spoofing...")
	resp, err := client.Get("https://httpbin.org/headers")
	if err != nil {
		fmt.Printf("[-] Network test failed (check internet connection): %v\n", err)
	} else {
		fmt.Printf("[+] Status Code: %d\n", resp.StatusCode)
		fmt.Printf("[+] Response Body:\n%s\n", resp.String())
	}

	// 3. Test Headless WASM Execution Engine
	fmt.Println("\n[+] Initializing Pure Go Headless WASM Engine (Wazero)...")
	ctx := context.Background()
	wasmRunner, err := gmrequest.NewWASMRunner(ctx, client.Profile)
	if err != nil {
		log.Fatalf("Failed to init wasm runner: %v", err)
	}
	defer wasmRunner.Close(ctx)

	fmt.Println("[+] WASM Engine initialized successfully with Mocked Browser Environment (Crypto, Timing, Screen, Memory)!")
	fmt.Println("[+] Ready to execute WASM challenges directly in memory without launching Chromium.")
}
