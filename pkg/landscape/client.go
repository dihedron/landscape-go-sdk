package landscape

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"resty.dev/v3"
)

// ErrNotSupported is returned by entity service methods that don't have a
// corresponding endpoint in the Landscape REST API (v2). It's wrapped in
// the returned error so callers can check for it with errors.Is.
var ErrNotSupported = errors.New("landscape: operation not supported by this endpoint")

// API wraps a resty.Client configured for the Landscape REST API v2.
// Each entity in the API (Computer, and others to come) is exposed as a
// public field holding its own service, e.g. Client.Computer. All of them
// share the same underlying resty client, so logging in once via Login
// authenticates every service.
type API struct {
	client *resty.Client

	mu    sync.RWMutex
	token string

	// Computer provides CRUD-shaped methods for the /computers endpoints.
	Computer *ComputerService
}

// Option configures a newly created API client, including authentication setup.
type Option func(*API)

// WithTimeout configures a custom timeout for the HTTP client.
func WithTimeout(timeout time.Duration) Option {
	return func(a *API) {
		a.client.SetTimeout(timeout)
	}
}

// WithClientTrustAnchors configures the root certificates for the HTTP client.
func WithClientTrustAnchors(paths ...string) Option {
	return func(a *API) {
		a.client.SetClientRootCertificates(paths...)
	}
}

// WithClientTrustAnchorsWatcher configures the root certificates for the HTTP client.
func WithClientTrustAnchorsWatcher(duration time.Duration, paths ...string) Option {
	return func(a *API) {
		a.client.SetClientRootCertificatesWatcher(&resty.CertWatcherOptions{PoolInterval: duration}, paths...)
	}
}

// WithClientCertificates configures the client certificates for the HTTP client;
// load the keypair from file (e.g. tlx.LoadX509KeyPair()) or use an in-memory certificate.
func WithClientCertificates(certs ...tls.Certificate) Option {
	return func(a *API) {
		a.client.SetCertificates(certs...)
	}
}

// WithClientCertificateFromFiles configures the client certificates for the HTTP client;
// it loads the keypair from the provided files.
func WithClientCertificateFromFiles(certPath string, keyPath string) Option {
	return func(a *API) {
		a.client.SetCertificateFromFile(certPath, keyPath)
	}
}

// WithClientCertificateFromString configures the client certificates for the HTTP client;
// it uses the provided certificate and key as strings.
func WithClientCertificateFromString(cert string, key string) Option {
	return func(a *API) {
		a.client.SetCertificateFromString(cert, key)
	}
}

// WithClientCertificatesWatcher configures the client certificates for the HTTP client.
func WithClientCertificatesWatcher(duration time.Duration, paths ...string) Option {
	return func(a *API) {
		a.client.SetClientRootCertificatesWatcher(&resty.CertWatcherOptions{PoolInterval: duration}, paths...)
	}
}

// WithTransport configures the transport for the HTTP client.
func WithTransport(transport http.RoundTripper) Option {
	return func(a *API) {
		a.client.SetTransport(transport)
	}
}

// WithRetry configures the retry count for the HTTP client.
func WithRetry(count int) Option {
	return func(a *API) {
		a.client.SetRetryCount(count)
	}
}

// WithDebug configures the debug mode for the HTTP client.
func WithDebug(enable bool) Option {
	return func(a *API) {
		a.client.SetDebug(enable)
	}
}

func WithTraceRequest(enabled bool) Option {
	return func(a *API) {
		a.client.SetTrace(enabled)
	}
}

func WithSaveResponse(enabled bool, path string) Option {
	return func(a *API) {
		a.client.SetResponseBodyUnlimitedReads(true).SetResponseSaveDirectory(path).SetResponseSaveToFile(enabled)
	}
}

// WithTokenAuth initializes the client with a JWT token and configures the underlying
// resty client to send it as the bearer token on every subsequent request.
func WithTokenAuth(token string) Option {
	return func(a *API) {
		a.mu.Lock()
		a.token = token
		a.mu.Unlock()
		a.client.SetAuthToken(token)
	}
}

// WithBasicAuth authenticates immediately using email/password and stores the returned
// token in the API, leaving the client ready for use without a separate login call.
func WithBasicAuth(email string, password string, account string) Option {
	return func(a *API) {
		if email == "" && password == "" {
			return
		}
		if _, err := a.login(context.Background(), email, password, account); err != nil {
			slog.Error("failed to initialize Landscape client with credentials", "error", err)
		}
	}
}

// New creates a Landscape API client. baseURL is your server root, e.g.
// "https://landscape.example.com" — do not include "/api/v2".
func New(baseURL string, opts ...Option) *API {
	client := resty.
		New().
		SetBaseURL(baseURL+"/api/v2").
		SetHeader("Accept", "application/json").
		SetLoggerWarnLevel(true).
		SetTimeout(30 * time.Second)

	api := &API{
		client:   client,
		Computer: &ComputerService{client: client},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(api)
		}
	}
	return api
}

func (a *API) Close() error {
	if a.Computer != nil {
		a.Computer.client = nil
	}
	// TODO: add more...
	if a.client != nil {
		return a.client.Close()
	}
	return nil
}

// LoginRequest is the body sent to POST /api/v2/login.
// Landscape authenticates with an email address, not a bare username.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// Account is optional; omit it to use the user's default account.
	Account string `json:"account,omitempty"`
}

// LoginResponse is what Landscape returns on a successful login.
type LoginResponse struct {
	Token          string    `json:"token"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	CurrentAccount string    `json:"current_account"`
	Accounts       []Account `json:"accounts"`
}

// Account describes one Landscape account the user has access to.
type Account struct {
	Title   string `json:"title"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

// APIError models Landscape's JSON error body for non-2xx responses.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("landscape api error: %s: %s", e.Code, e.Message)
}

// login authenticates with email/password and stores the returned JWT on
// the underlying resty client. Every subsequent request made through this
// Client will automatically carry "Authorization: Bearer <token>", so you
// never need to attach it yourself on later calls.
func (a *API) login(ctx context.Context, email, password, account string) (*LoginResponse, error) {
	var result LoginResponse
	var failure APIError

	resp, err := a.client.R().
		SetContext(ctx).
		SetBody(LoginRequest{
			Email:    email,
			Password: password,
			Account:  account,
		}).
		SetResult(&result).
		SetResultError(&failure).
		Post("/login")
	if err != nil {
		slog.Error("error performing login request", "error", err)
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	if resp.IsStatusFailure() {
		if failure.Message != "" {
			return nil, &failure
		}
		return nil, fmt.Errorf("login failed: status %d: %s", resp.StatusCode(), resp.String())
	}

	a.mu.Lock()
	a.token = result.Token
	a.mu.Unlock()

	// This is the "store for later calls" step: resty attaches this
	// bearer token to every request this client issues from now on.
	a.client.SetAuthToken(result.Token)

	return &result, nil
}

// Token returns the currently stored JWT (empty string if not logged in).
func (a *API) Token() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.token
}
