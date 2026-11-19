package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/DauXanhRauMaRauMaDauNanh/gmrequest/pkg/profile"
)

// Response wraps the HTTP response with convenient helpers
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Raw        *http.Response
}

func (r *Response) String() string {
	return string(r.Body)
}

func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// Client is the high-performance stealth request engine
type Client struct {
	httpClient tls_client.HttpClient
	Profile    *profile.DeviceProfile
}

type ClientOptions struct {
	ProxyURL        string
	TimeoutSeconds  int
	FollowRedirects bool
	Profile         *profile.DeviceProfile
}

// NewClient creates a new stealth client with full TLS/H2/Header fingerprint spoofing
func NewClient(opts ...ClientOptions) (*Client, error) {
	var opt ClientOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Profile == nil {
		opt.Profile = profile.Random()
	}

	if opt.TimeoutSeconds <= 0 {
		opt.TimeoutSeconds = 30
	}

	jar := tls_client.NewCookieJar()

	tlsOptions := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(opt.TimeoutSeconds),
		tls_client.WithClientProfile(opt.Profile.TLSProfile),
		tls_client.WithCookieJar(jar),
	}

	if !opt.FollowRedirects {
		tlsOptions = append(tlsOptions, tls_client.WithNotFollowRedirects())
	}

	if opt.ProxyURL != "" {
		tlsOptions = append(tlsOptions, tls_client.WithProxyUrl(opt.ProxyURL))
	}

	httpClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), tlsOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create stealth http client: %w", err)
	}

	return &Client{
		httpClient: httpClient,
		Profile:    opt.Profile,
	}, nil
}

// Get performs a stealth GET request with standard browser navigation headers
func (c *Client) Get(targetURL string, customHeaders ...map[string]string) (*Response, error) {
	return c.Do("GET", targetURL, nil, false, customHeaders...)
}

// PostJSON performs a stealth POST request with JSON payload and AJAX headers
func (c *Client) PostJSON(targetURL string, data interface{}, customHeaders ...map[string]string) (*Response, error) {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	return c.Do("POST", targetURL, bytes.NewReader(bodyBytes), true, customHeaders...)
}

// PostForm performs a stealth POST request with form urlencoded payload
func (c *Client) PostForm(targetURL string, formData url.Values, customHeaders ...map[string]string) (*Response, error) {
	return c.Do("POST", targetURL, strings.NewReader(formData.Encode()), true, customHeaders...)
}

// Do executes a raw stealth request with customized Chrome header ordering
func (c *Client) Do(method, targetURL string, body io.Reader, isAPI bool, customHeaders ...map[string]string) (*Response, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target url: %w", err)
	}

	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Apply browser headers based on request type
	if isAPI {
		ApplyChromeAPIHeaders(req.Header, parsedURL.Host, c.Profile, "application/json")
	} else {
		ApplyChromeNavHeaders(req.Header, parsedURL.Host, c.Profile)
	}

	// Merge any user-provided custom headers
	if len(customHeaders) > 0 {
		for k, v := range customHeaders[0] {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBytes,
		Raw:        resp,
	}, nil
}

// GetCookies returns the cookies stored in the jar for a given URL
func (c *Client) GetCookies(u *url.URL) []*http.Cookie {
	return c.httpClient.GetCookies(u)
}

// SetCookies adds cookies to the client's cookie jar
func (c *Client) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.httpClient.SetCookies(u, cookies)
}

// SetProxy updates the proxy dynamically
func (c *Client) SetProxy(proxyURL string) error {
	return c.httpClient.SetProxy(proxyURL)
}
