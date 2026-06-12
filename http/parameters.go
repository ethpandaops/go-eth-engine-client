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

package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type parameters struct {
	log           logrus.FieldLogger
	address       string
	timeout       time.Duration
	client        *http.Client
	jwtSecret     []byte
	jwtSecretFile string
	jwtID         string
	jwtClientVer  string
}

// Parameter is a functional option for the http Service.
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
// "http://localhost:8551".
func WithAddress(address string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.address = address
	})
}

// WithTimeout sets the default per-request timeout.
func WithTimeout(timeout time.Duration) Parameter {
	return parameterFunc(func(p *parameters) {
		p.timeout = timeout
	})
}

// WithHTTPClient sets a custom *http.Client. When unset, a client with the
// configured timeout is used.
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

// WithJWTClientVersion sets the optional `clv` claim included in minted JWTs.
func WithJWTClientVersion(clv string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.jwtClientVer = clv
	})
}

const defaultTimeout = 10 * time.Second

// parseAndCheckParameters applies the options and validates them.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	p := &parameters{
		timeout: defaultTimeout,
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
