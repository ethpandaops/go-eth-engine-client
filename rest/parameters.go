// Copyright © 2026 ethPandaOps.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type parameters struct {
	log                 logrus.FieldLogger
	address             string
	basePath            string
	timeout             time.Duration
	client              *http.Client
	jwtSecret           []byte
	jwtSecretFile       string
	jwtID               string
	clientVersionHeader string
}

// Parameter is a functional option for the rest Service.
type Parameter interface {
	apply(*parameters)
}

type parameterFunc func(*parameters)

func (f parameterFunc) apply(p *parameters) {
	f(p)
}

// WithLogger sets the logger for the service.
func WithLogger(log logrus.FieldLogger) Parameter {
	return parameterFunc(func(p *parameters) {
		p.log = log
	})
}

// WithAddress sets the base URL of the Engine API endpoint, e.g.
// "http://localhost:8551". The Marius spec mandates h2c (HTTP/2 cleartext)
// for plain http:// URLs.
func WithAddress(address string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.address = address
	})
}

// WithBasePath sets the path prefix for hot-path endpoints. Defaults to
// "/engine/v2" per the Marius spec.
func WithBasePath(basePath string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.basePath = basePath
	})
}

// WithTimeout sets the default per-request timeout.
func WithTimeout(timeout time.Duration) Parameter {
	return parameterFunc(func(p *parameters) {
		p.timeout = timeout
	})
}

// WithHTTPClient sets a custom *http.Client. When unset, a client with
// HTTP/2 (h2c for cleartext) and the configured timeout is used.
func WithHTTPClient(client *http.Client) Parameter {
	return parameterFunc(func(p *parameters) {
		p.client = client
	})
}

// WithJWTSecret sets the 32-byte JWT secret used to authenticate requests.
func WithJWTSecret(secret []byte) Parameter {
	return parameterFunc(func(p *parameters) {
		p.jwtSecret = secret
	})
}

// WithJWTSecretFile sets the path to a file holding the hex-encoded 32-byte
// JWT secret. Ignored when WithJWTSecret is also supplied.
func WithJWTSecretFile(path string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.jwtSecretFile = path
	})
}

// WithJWTID sets the optional `id` claim included in minted JWTs.
func WithJWTID(id string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.jwtID = id
	})
}

// WithClientVersionHeader sets the value of the X-Engine-Client-Version
// request header sent on every call. The Marius spec moves the CL client
// version off the JWT `clv` claim and onto this header.
func WithClientVersionHeader(value string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.clientVersionHeader = value
	})
}

const (
	defaultTimeout  = 10 * time.Second
	defaultBasePath = DefaultBasePath
)

// parseAndCheckParameters applies the options and validates them.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	p := &parameters{
		timeout:  defaultTimeout,
		basePath: defaultBasePath,
	}

	for _, param := range params {
		param.apply(p)
	}

	if p.log == nil {
		p.log = logrus.StandardLogger()
	}

	if p.address == "" {
		return nil, errors.New("no address specified")
	}

	if len(p.jwtSecret) == 0 && p.jwtSecretFile == "" {
		return nil, errors.New("no JWT secret specified (use WithJWTSecret or WithJWTSecretFile)")
	}

	return p, nil
}

// newDefaultHTTPClient returns a stdlib *http.Client. The Marius spec
// mandates HTTP/2 (h2c for cleartext); for a production deployment supply
// an h2c-capable client via [WithHTTPClient] -- e.g. one built around
// `golang.org/x/net/http2` with `AllowHTTP: true` -- so the same connection
// can be reused across requests. Go's default transport will negotiate h2
// against TLS endpoints automatically.
func newDefaultHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
