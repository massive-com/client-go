package models

import (
	"net/http"
	"net/url"
	"time"
)

// RequestOptions are used to configure client calls.
type RequestOptions struct {
	// APIKey to pass with the request
	APIKey *string

	// Headers to apply to the request
	Headers http.Header

	// Query params to apply to the request
	QueryParams url.Values

	// Trace enables request tracing
	Trace bool

	// Limiter used to rate limit requests
	Limiter RateLimiter
}

// RequestOption changes the configuration of RequestOptions.
type RequestOption func(o *RequestOptions)

// APIKey sets an APIKey as an option.
func APIKey(id string) RequestOption {
	return func(o *RequestOptions) {
		o.APIKey = &id
	}
}

// Header sets a header as an option.
func Header(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Headers == nil {
			o.Headers = make(http.Header)
		}

		o.Headers.Add(key, value)
	}
}

// QueryParam sets a query param as an option.
func QueryParam(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.QueryParams == nil {
			o.QueryParams = make(url.Values)
		}

		o.QueryParams.Add(key, value)
	}
}

// Headers required to use the Launchpad product.
const (
	// HeaderEdgeID is a required Launchpad header. It identifies the Edge User requesting data.
	HeaderEdgeID = "X-Massive-Edge-ID"
	// HeaderEdgeIPAddress is a required Launchpad header. It denotes the originating IP Address of the Edge User requesting data.
	HeaderEdgeIPAddress = "X-Massive-Edge-IP-Address"
	// HeaderEdgeUserAgent is an optional Launchpad header. It denotes the originating UserAgent of the Edge User requesting data.
	HeaderEdgeUserAgent = "X-Massive-Edge-User-Agent"
)

// RequiredEdgeHeaders sets the required headers for the Launchpad product.
func RequiredEdgeHeaders(edgeID, edgeIPAddress string) RequestOption {
	return func(o *RequestOptions) {
		if o.Headers == nil {
			o.Headers = make(http.Header)
		}

		o.Headers.Add(HeaderEdgeID, edgeID)
		o.Headers.Add(HeaderEdgeIPAddress, edgeIPAddress)
	}
}

// EdgeUserAgent sets the Launchpad optional header denoting the Edge User's UserAgent.
func EdgeUserAgent(userAgent string) RequestOption {
	return Header(HeaderEdgeUserAgent, userAgent)
}

// WithTrace enables or disables request tracing.
func WithTrace(trace bool) RequestOption {
	return func(o *RequestOptions) {
		o.Trace = trace
	}
}

type RateLimiter interface {
	Take() time.Time
}

// WithRateLimiter sets Ratelimiter to limit request rate.
func WithRateLimiter(limiter RateLimiter) RequestOption {
	return func(o *RequestOptions) {
		o.Limiter = limiter
	}
}
