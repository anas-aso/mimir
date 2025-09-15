// SPDX-License-Identifier: AGPL-3.0-only

package propagation

import (
	"context"
	"net/http"
)

// Propagator represents something that extracts auxiliary information to and from a request.
type Propagator interface {
	// ReadFromCarrier extracts auxiliary information from a request (represented by a carrier)
	// and returns a new context derived from ctx with that information included.
	ReadFromCarrier(ctx context.Context, carrier Carrier) (context.Context, error)

	// AddToCarrier adds auxiliary information to a request (represented by a carrier) from ctx.
	AddToCarrier(ctx context.Context, carrier Carrier) error
}

type MultiPropagator struct {
	Propagators []Propagator
}

func (m *MultiPropagator) ReadFromCarrier(ctx context.Context, carrier Carrier) (context.Context, error) {
	for _, p := range m.Propagators {
		var err error
		ctx, err = p.ReadFromCarrier(ctx, carrier)
		if err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func (m *MultiPropagator) AddToCarrier(ctx context.Context, carrier Carrier) error {
	for _, p := range m.Propagators {
		if err := p.AddToCarrier(ctx, carrier); err != nil {
			return err
		}
	}

	return nil
}

// Carrier represents a carrier of key-value pairs for a request, such as HTTP headers.
type Carrier interface {
	// Get returns the value with the given name, or an empty string if it is not present.
	Get(name string) string

	Set(name string, value string)
}

type MapCarrier map[string]string

func (m MapCarrier) Get(name string) string {
	return m[name]
}

func (m MapCarrier) Set(name string, value string) {
	m[name] = value
}

type HttpHeaderCarrier http.Header

func (h HttpHeaderCarrier) Get(name string) string {
	return http.Header(h).Get(name)
}

func (h HttpHeaderCarrier) Set(name string, value string) {
	http.Header(h).Set(name, value)
}
