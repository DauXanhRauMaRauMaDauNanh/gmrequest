package engine

import (
	http "github.com/bogdanfinn/fhttp"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/profile"
)

// Standard Chrome Header Order for navigation/GET requests
var ChromeNavHeaderOrder = []string{
	"host",
	"connection",
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"sec-ch-ua-platform",
	"upgrade-insecure-requests",
	"user-agent",
	"accept",
	"sec-fetch-site",
	"sec-fetch-mode",
	"sec-fetch-user",
	"sec-fetch-dest",
	"sec-ch-ua-arch",
	"sec-ch-ua-bitness",
	"sec-ch-ua-full-version",
	"sec-ch-ua-full-version-list",
	"sec-ch-ua-model",
	"sec-ch-ua-platform-version",
	"referer",
	"accept-encoding",
	"accept-language",
	"priority",
	"cookie",
}

// Standard Chrome Header Order for XHR/Fetch/API POST requests
var ChromeFetchHeaderOrder = []string{
	"host",
	"connection",
	"content-length",
	"sec-ch-ua-platform",
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"user-agent",
	"content-type",
	"accept",
	"origin",
	"sec-fetch-site",
	"sec-fetch-mode",
	"sec-fetch-dest",
	"referer",
	"accept-encoding",
	"accept-language",
	"priority",
	"cookie",
}

// Pseudo-header ordering for HTTP/2: :method, :authority, :scheme, :path
var ChromeH2PseudoHeaderOrder = []string{
	":method",
	":authority",
	":scheme",
	":path",
}

// ApplyChromeNavHeaders applies standard navigation headers matching a DeviceProfile
func ApplyChromeNavHeaders(headers http.Header, host string, p *profile.DeviceProfile) {
	headers.Set("host", host)
	headers.Set("sec-ch-ua", p.SecChUa)
	if p.Mobile {
		headers.Set("sec-ch-ua-mobile", "?1")
	} else {
		headers.Set("sec-ch-ua-mobile", "?0")
	}
	headers.Set("sec-ch-ua-platform", `"`+string(p.Platform)+`"`)
	headers.Set("upgrade-insecure-requests", "1")
	headers.Set("user-agent", p.UserAgent)
	headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Set("sec-fetch-site", "none")
	headers.Set("sec-fetch-mode", "navigate")
	headers.Set("sec-fetch-user", "?1")
	headers.Set("sec-fetch-dest", "document")
	headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	headers.Set("accept-language", p.AcceptLang)
	headers.Set("priority", "u=0, i")

	// Set exact header ordering key for fhttp
	headers[http.HeaderOrderKey] = ChromeNavHeaderOrder
	headers[http.PHeaderOrderKey] = ChromeH2PseudoHeaderOrder
}

// ApplyChromeAPIHeaders applies standard AJAX/JSON POST/GET headers matching a DeviceProfile
func ApplyChromeAPIHeaders(headers http.Header, host string, p *profile.DeviceProfile, contentType string) {
	headers.Set("host", host)
	headers.Set("sec-ch-ua-platform", `"`+string(p.Platform)+`"`)
	headers.Set("sec-ch-ua", p.SecChUa)
	if p.Mobile {
		headers.Set("sec-ch-ua-mobile", "?1")
	} else {
		headers.Set("sec-ch-ua-mobile", "?0")
	}
	headers.Set("user-agent", p.UserAgent)
	if contentType != "" {
		headers.Set("content-type", contentType)
	}
	headers.Set("accept", "application/json, text/plain, */*")
	headers.Set("sec-fetch-site", "same-origin")
	headers.Set("sec-fetch-mode", "cors")
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	headers.Set("accept-language", p.AcceptLang)
	headers.Set("priority", "u=1, i")

	headers[http.HeaderOrderKey] = ChromeFetchHeaderOrder
	headers[http.PHeaderOrderKey] = ChromeH2PseudoHeaderOrder
}
