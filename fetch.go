// +build js,wasm

// Package fetch is a WIP Web Assembly fetch wrapper that avoids importing net/http.
package fetch

import (
	"io"
	"io/ioutil"
	"syscall/js"
)

// cache enums
const (
	CacheDefault      = "default"
	CacheNoStore      = "no-store"
	CacheReload       = "reload"
	CacheNone         = "no-cache"
	CacheForce        = "force-cache"
	CacheOnlyIfCached = "only-if-cached"
)

// credentials enums
const (
	CredentialsOmit       = "omit"
	CredentialsSameOrigin = "same-origin"
	CredentialsInclude    = "include"
)

// Common HTTP methods.
//
// Unless otherwise noted, these are defined in RFC 7231 section 4.3.
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

// Mode enums
const (
	ModeSameOrigin = "same-origin"
	ModeNoCORS     = "no-cors"
	ModeCORS       = "cors"
	ModeNavigate   = "navigate"
)

// Redirect enums
const (
	RedirectFollow = "follow"
	RedirectError  = "error"
	RedirectManual = "manual"
)

// Referrer enums
const (
	ReferrerNone   = "no-referrer"
	ReferrerClient = "client"
)

// ReferrerPolicy enums
const (
	ReferrerPolicyNone        = "no-referrer"
	ReferrerPolicyNoDowngrade = "no-referrer-when-downgrade"
	ReferrerPolicyOrigin      = "origin"
	ReferrerPolicyCrossOrigin = "origin-when-cross-origin"
	ReferrerPolicyUnsafeURL   = "unsafe-url"
)

// Opts opts
type Opts struct {
	// Method is the http verb (constants are copied from net/http to avoid import)
	Method string

	// Headers is a map of http headers to send.
	Headers map[string]string

	// Body is the body request
	Body io.ReadCloser

	// Mode docs https://developer.mozilla.org/en-US/docs/Web/API/Request/mode
	Mode string

	// Credentials docs https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
	Credentials string

	// Cache docs https://developer.mozilla.org/en-US/docs/Web/API/Request/cache
	Cache string

	// Redirect docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	Redirect string

	// Referrer docs https://developer.mozilla.org/en-US/docs/Web/API/Request/referrer
	Referrer string

	// ReferrerPolicy docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	ReferrerPolicy string

	// Integrity docs https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
	Integrity string

	// KeepAlive docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	KeepAlive bool

	// Signal docs https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal
	Signal chan struct{}
}

// oof.
func mapOpts(opts *Opts) map[string]interface{} {
	mp := map[string]interface{}{
		"Method":         opts.Method,
		"Headers":        mapHeaders(opts.Headers),
		"Mode":           opts.Mode,
		"Credentials":    opts.Credentials,
		"Cache":          opts.Cache,
		"Redirect":       opts.Redirect,
		"Referrer":       opts.Referrer,
		"ReferrerPolicy": opts.ReferrerPolicy,
		"Integrity":      opts.Integrity,
		"KeepAlive":      opts.KeepAlive,
	}
	if opts.Signal != nil {
		// TODO: do signal
	}
	if opts.Body != nil {
		bts, err := ioutil.ReadAll(opts.Body)
		if err != nil {
			panic(err) // TODO: return err
		}

		mp["body"] = string(bts)
	}

	return mp
}

func mapHeaders(mp map[string]string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range mp {
		newMap[k] = v
	}
	return newMap
}

// Fetch fetches
func Fetch(url string, opts *Opts) []byte {
	ch := make(chan string)
	js.Global().Call("fetch", url, mapOpts(opts)).Call("then", js.NewCallback(func(args []js.Value) {
		args[0].Call("text").Call("then", js.NewCallback(func(args []js.Value) {
			ch <- args[0].String()
		}))
	}))

	return []byte(<-ch)
}
