package main

import (
	"fmt"
	"strings"
	"time"
)

// --- Middleware Chaining Pattern ---
//
// Middleware wraps a handler with additional behavior (logging, auth,
// rate limiting, tracing) without touching the core handler logic.
// Each middleware is a function that takes a Handler and returns a
// new Handler — it can run code before, after, or around the inner call.
//
// This works for any handler shape, not just net/http. Here we use a
// simple Request/Response to keep the example self-contained.
//
// Chain applies middlewares so the first listed runs outermost (first
// to intercept the request, last to see the response).

// Request and Response are minimal stand-ins for http.Request / http.Response.
type Request struct {
	Path   string
	User   string
	Body   string
}

type Response struct {
	Status int
	Body   string
}

// Handler is the core function type: takes a request, returns a response.
// In net/http this is the http.Handler interface; here it's a plain func.
type Handler func(req Request) Response

// Middleware wraps a Handler to add behavior before/after the inner call.
type Middleware func(next Handler) Handler

// --- Core handler ---

// coreHandler is the application logic with no cross-cutting concerns mixed in.
// Middlewares keep it clean.
func coreHandler(req Request) Response {
	return Response{
		Status: 200,
		Body:   fmt.Sprintf("Hello from %s, user=%s", req.Path, req.User),
	}
}

// --- Middlewares ---

// LoggerMiddleware logs the path and status code around every request.
// It sees the request before the inner handler and the response after.
func LoggerMiddleware(next Handler) Handler {
	return func(req Request) Response {
		fmt.Printf("[logger]  → %s\n", req.Path)
		resp := next(req) // call the next handler in the chain
		fmt.Printf("[logger]  ← %d\n", resp.Status)
		return resp
	}
}

// AuthMiddleware checks that the request has a non-empty User field.
// If auth fails, it short-circuits the chain — the inner handler is never called.
func AuthMiddleware(next Handler) Handler {
	return func(req Request) Response {
		if strings.TrimSpace(req.User) == "" {
			// Short-circuit: return 401 without calling the rest of the chain.
			fmt.Println("[auth]    DENIED — no user")
			return Response{Status: 401, Body: "unauthorized"}
		}
		fmt.Printf("[auth]    OK — user=%s\n", req.User)
		return next(req)
	}
}

// TimingMiddleware records how long the inner chain takes and annotates
// the response body. In production you'd emit a metric instead.
func TimingMiddleware(next Handler) Handler {
	return func(req Request) Response {
		start := time.Now()
		resp := next(req)
		elapsed := time.Since(start)
		fmt.Printf("[timing]  %s took %v\n", req.Path, elapsed)
		return resp
	}
}

// --- Chain builder ---

// Chain applies middlewares to a handler in order: the first middleware
// in the list is the outermost wrapper (first to run on the way in,
// last to run on the way out). Applied right-to-left so execution order
// is left-to-right.
//
//   Chain(h, Logger, Auth, Timing)
//   execution order: Logger → Auth → Timing → h → Timing → Auth → Logger
func Chain(h Handler, middlewares ...Middleware) Handler {
	// Iterate in reverse so the first middleware ends up as the outermost wrapper.
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// DemoMiddleware wires three middlewares around the core handler
// and fires two requests: one with a user (succeeds) and one without (denied).
func DemoMiddleware() {
	fmt.Println("=== Middleware Chaining ===")

	// Build the handler once; the chain is reusable for every request.
	handler := Chain(coreHandler, LoggerMiddleware, AuthMiddleware, TimingMiddleware)

	fmt.Println("-- request with auth --")
	resp := handler(Request{Path: "/api/hello", User: "alice", Body: ""})
	fmt.Printf("response: %d  %s\n", resp.Status, resp.Body)

	fmt.Println()
	fmt.Println("-- request without auth --")
	resp = handler(Request{Path: "/api/hello", User: "", Body: ""})
	fmt.Printf("response: %d  %s\n", resp.Status, resp.Body)
}
