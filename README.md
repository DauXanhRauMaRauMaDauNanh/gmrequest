# gmrequest (Full Stealth HTTP/2 & TLS Engine for Go)

`gmrequest` là thư viện Golang gửi request siêu tàng hình (**Ultra-Stealth Request Engine**) không cần khởi chạy Chromium hay Headless Browser, kết hợp tinh hoa từ:
1. **CloakBrowser & RayoBrowse**: Kiến trúc **Coherent Device Profiling** (đồng bộ User-Agent, Client Hints, Platform, Geometry màn hình, Hardware Concurrency dựa trên Seed cố định).
2. **CloudScraper Go / uTLS**: Khắc phục các điểm yếu chí mạng của `cloudscraper_go` (User-Agent Chrome 50 từ 2016, thiếu Client Hints, không hỗ trợ zstd/br, thiếu HTTP/2 frame control).
3. **Headless WASM Engine (Wazero)**: Tích hợp runtime WebAssembly JIT thuần Go để giải toán mã hóa / Proof-of-Work (PoW) trực tiếp trong RAM mà không cần mở trình duyệt.

---

## Tính năng nổi bật

* **Zero Browser Overhead**: Chạy thuần Go, tốn cực ít RAM, xử lý hàng nghìn request/giây.
* **JA3 / JA4 & TLS 1.3 Spoofing**: Giả lập ClientHello chuẩn Google Chrome 130+ (Cipher Suites, GREASE, Supported Groups X25519Kyber768, ALPN h2).
* **HTTP/2 Akamai Fingerprint**: SETTINGS frame, WINDOW_UPDATE (15663105), Pseudo-header order (`:method`, `:authority`, `:scheme`, `:path`).
* **Deterministic Header Ordering**: Sắp xếp thứ tự Headers theo đúng chuẩn Chrome phát ra để tránh bị WAF bắt bài thứ tự header.
* **Next-Gen Decompression**: Hỗ trợ giải nén tự động `gzip`, `deflate`, `br` (Brotli), và `zstd` (Zstandard).
* **Headless WASM Runner**: Chạy module WASM challenge của anti-bot với môi trường giả lập (WebCrypto, Timing, Memory, Screen).

---

## Cài đặt & Sử dụng

### Cài đặt

```bash
go get github.com/DauXanhRauMaRauMaDauNanh/gmrequest
```

### 1. Khởi tạo Client ngẫu nhiên hoặc theo Seed

```go
package main

import (
	"fmt"
	"log"

	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest"
)

func main() {
	// Khởi tạo client với Seed (giống --fingerprint của CloakBrowser)
	client, err := gmrequest.NewClientWithSeed(123456789, gmrequest.PlatformWindows)
	if err != nil {
		log.Fatal(err)
	}

	// Gửi GET request
	resp, err := client.Get("https://httpbin.org/headers")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Body:", resp.String())
}
```

### 2. Gửi API POST JSON với AJAX Header Order

```go
data := map[string]interface{}{
    "action": "login",
    "token": "example",
}

resp, err := client.PostJSON("https://api.example.com/v1/auth", data)
```

### 3. Thực thi WASM Challenge trong Go (Không mở browser)

```go
ctx := context.Background()
runner, err := gmrequest.NewWASMRunner(ctx, client.Profile)
if err != nil {
    log.Fatal(err)
}
defer runner.Close(ctx)

// Chạy trực tiếp bytecode WASM tải về từ server
wasmBytes := []byte{...}
results, err := runner.RunModule(ctx, wasmBytes, "solve", 1000)
```
