package main

import (
	"fmt"
	"time"
)

// --- Functional Options Pattern ---
//
// Instead of a long constructor parameter list or a separate config struct,
// options are expressed as functions of type func(*Server). The constructor
// applies each one in order, so callers only pass what they want to change
// and defaults handle the rest.
//
// This keeps the API stable: adding a new option never breaks existing callers.

// Server is the type being configured. Fields are unexported so the only
// way to set them is through the provided With* option functions.
type Server struct {
	host        string
	port        int
	timeout     time.Duration
	maxConns    int
	tlsEnabled  bool
}

// Option is a function that mutates a Server during construction.
// Callers compose behavior by passing any number of these to NewServer.
type Option func(*Server)

// WithHost overrides the default host.
func WithHost(host string) Option {
	return func(s *Server) { s.host = host }
}

// WithPort overrides the default port.
func WithPort(port int) Option {
	return func(s *Server) { s.port = port }
}

// WithTimeout sets the per-request timeout.
func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

// WithMaxConns limits how many simultaneous connections are accepted.
func WithMaxConns(n int) Option {
	return func(s *Server) { s.maxConns = n }
}

// WithTLS enables TLS on the server.
func WithTLS() Option {
	return func(s *Server) { s.tlsEnabled = true }
}

// NewServer constructs a Server with sensible defaults, then applies
// each caller-supplied option in order. New options can be added in
// the future without changing this signature or breaking callers.
func NewServer(opts ...Option) *Server {
	// Start with defaults — callers only override what they care about.
	s := &Server{
		host:     "localhost",
		port:     8080,
		timeout:  30 * time.Second,
		maxConns: 100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// DemoFunctionalOptions shows three different callers each configuring
// a Server differently without any of them affecting the others.
func DemoFunctionalOptions() {
	fmt.Println("=== Functional Options ===")

	// Default server — no options, all defaults apply.
	defaultSrv := NewServer()
	fmt.Printf("default:  %s:%d  timeout=%v  tls=%v\n",
		defaultSrv.host, defaultSrv.port, defaultSrv.timeout, defaultSrv.tlsEnabled)

	// Production server — override several options.
	prodSrv := NewServer(
		WithHost("0.0.0.0"),
		WithPort(443),
		WithTLS(),
		WithTimeout(10*time.Second),
		WithMaxConns(1000),
	)
	fmt.Printf("prod:     %s:%d  timeout=%v  tls=%v  maxConns=%d\n",
		prodSrv.host, prodSrv.port, prodSrv.timeout, prodSrv.tlsEnabled, prodSrv.maxConns)

	// Dev server — only change the port.
	devSrv := NewServer(WithPort(3000))
	fmt.Printf("dev:      %s:%d  timeout=%v  tls=%v\n",
		devSrv.host, devSrv.port, devSrv.timeout, devSrv.tlsEnabled)
}
