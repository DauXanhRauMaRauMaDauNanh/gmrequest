package profile

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bogdanfinn/tls-client/profiles"
)

// Platform represents the operating system platform
type Platform string

const (
	PlatformWindows Platform = "Windows"
	PlatformMacOS   Platform = "macOS"
	PlatformLinux   Platform = "Linux"
	PlatformAndroid Platform = "Android"
	PlatformIOS     Platform = "iOS"
)

// DeviceProfile encapsulates a coherent fingerprint profile
// inspired by CloakBrowser and RayoBrowse's device coherence architecture.
type DeviceProfile struct {
	Seed          int64
	BrowserName   string
	BrowserVer    string
	FullVer       string
	Platform      Platform
	PlatformVersion string
	Architecture  string
	Bitness       string
	Model         string
	Mobile        bool

	// Network & Identity
	UserAgent     string
	SecChUa       string
	SecChUaFull   string
	Locale        string
	Timezone      string
	AcceptLang    string

	// Screen & Geometry
	ScreenWidth   int
	ScreenHeight  int
	AvailWidth    int
	AvailHeight   int
	DevicePixelRatio float64
	HardwareConcurrency int
	DeviceMemory        int

	// TLS & Protocol Profile
	TLSProfile    profiles.ClientProfile
}

// GenerateFromSeed creates a deterministic, coherent device profile from a numeric seed.
// Just like CloakBrowser's --fingerprint=seed, the same seed will always produce identical signals.
func GenerateFromSeed(seed int64, platform Platform) *DeviceProfile {
	r := rand.New(rand.NewSource(seed))

	switch platform {
	case PlatformMacOS:
		return generateMacOSProfile(r, seed)
	default:
		return generateWindowsProfile(r, seed)
	}
}

func generateWindowsProfile(r *rand.Rand, seed int64) *DeviceProfile {
	// Chrome 131 - 133 series
	majorVer := 131 + r.Intn(3) // 131, 132, 133
	buildVer := 6778 + r.Intn(100)
	patchVer := 85 + r.Intn(120)
	fullVer := fmt.Sprintf("%d.0.%d.%d", majorVer, buildVer, patchVer)

	ua := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", fullVer)
	secChUa := fmt.Sprintf(`"Google Chrome";v="%d", "Chromium";v="%d", "Not_A Brand";v="24"`, majorVer, majorVer)
	secChUaFull := fmt.Sprintf(`"Google Chrome";v="%s", "Chromium";v="%s", "Not_A Brand";v="24.0.0.0"`, fullVer, fullVer)

	resolutions := [][2]int{
		{1920, 1080},
		{2560, 1440},
		{1536, 864},
		{1366, 768},
	}
	res := resolutions[r.Intn(len(resolutions))]

	concurrencyList := []int{8, 12, 16, 24, 32}
	concurrency := concurrencyList[r.Intn(len(concurrencyList))]

	memList := []int{8, 16, 32}
	mem := memList[r.Intn(len(memList))]

	return &DeviceProfile{
		Seed:                seed,
		BrowserName:         "Chrome",
		BrowserVer:          fmt.Sprintf("%d", majorVer),
		FullVer:             fullVer,
		Platform:            PlatformWindows,
		PlatformVersion:     "15.0.0", // Windows 11/10
		Architecture:        "x86",
		Bitness:             "64",
		Model:               "",
		Mobile:              false,
		UserAgent:           ua,
		SecChUa:             secChUa,
		SecChUaFull:         secChUaFull,
		Locale:              "en-US",
		Timezone:            "America/New_York",
		AcceptLang:          "en-US,en;q=0.9",
		ScreenWidth:         res[0],
		ScreenHeight:        res[1],
		AvailWidth:          res[0],
		AvailHeight:         res[1] - 40, // Windows taskbar offset
		DevicePixelRatio:    1.0,
		HardwareConcurrency: concurrency,
		DeviceMemory:        mem,
		TLSProfile:          profiles.Chrome_131,
	}
}

func generateMacOSProfile(r *rand.Rand, seed int64) *DeviceProfile {
	majorVer := 131 + r.Intn(3)
	buildVer := 6778 + r.Intn(100)
	patchVer := 85 + r.Intn(120)
	fullVer := fmt.Sprintf("%d.0.%d.%d", majorVer, buildVer, patchVer)

	ua := fmt.Sprintf("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", fullVer)
	secChUa := fmt.Sprintf(`"Google Chrome";v="%d", "Chromium";v="%d", "Not_A Brand";v="24"`, majorVer, majorVer)
	secChUaFull := fmt.Sprintf(`"Google Chrome";v="%s", "Chromium";v="%s", "Not_A Brand";v="24.0.0.0"`, fullVer, fullVer)

	return &DeviceProfile{
		Seed:                seed,
		BrowserName:         "Chrome",
		BrowserVer:          fmt.Sprintf("%d", majorVer),
		FullVer:             fullVer,
		Platform:            PlatformMacOS,
		PlatformVersion:     "14.5.0", // Sonoma
		Architecture:        "arm",
		Bitness:             "64",
		Model:               "",
		Mobile:              false,
		UserAgent:           ua,
		SecChUa:             secChUa,
		SecChUaFull:         secChUaFull,
		Locale:              "en-US",
		Timezone:            "America/Los_Angeles",
		AcceptLang:          "en-US,en;q=0.9",
		ScreenWidth:         1728,
		ScreenHeight:        1117,
		AvailWidth:          1728,
		AvailHeight:         1080,
		DevicePixelRatio:    2.0, // Retina
		HardwareConcurrency: 10,
		DeviceMemory:        16,
		TLSProfile:          profiles.Chrome_131,
	}
}

// Random returns a random coherent profile
func Random() *DeviceProfile {
	seed := time.Now().UnixNano()
	return GenerateFromSeed(seed, PlatformWindows)
}
