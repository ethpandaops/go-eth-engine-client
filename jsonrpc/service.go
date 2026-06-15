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

// Package jsonrpc provides the legacy JSON-RPC implementation of the Engine
// API client interfaces defined in the root engine package. It authenticates
// with the JWT (HS256) scheme required by the Engine API and dispatches each
// fork-agnostic request to the appropriate versioned `engine_*` JSON-RPC
// method.
//
// For the upcoming REST + SSZ transport (execution-apis PR #793 -- Marius
// spec), see the sibling `rest` package.
package jsonrpc

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	engine "github.com/ethpandaops/go-eth-engine-client"
)

// Compile-time assertions that the Service satisfies the engine interfaces.
var (
	_ engine.Service                   = (*Service)(nil)
	_ engine.NewPayloadProvider        = (*Service)(nil)
	_ engine.ForkchoiceUpdatedProvider = (*Service)(nil)
	_ engine.GetPayloadProvider        = (*Service)(nil)
	_ engine.GetBlobsProvider          = (*Service)(nil)
	_ engine.GetPayloadBodiesProvider  = (*Service)(nil)
	_ engine.CapabilitiesProvider      = (*Service)(nil)
	_ engine.ClientVersionProvider     = (*Service)(nil)
)

// Service is a JSON-RPC Engine API client.
type Service struct {
	log     logrus.FieldLogger
	address string
	client  *http.Client
	timeout time.Duration
	signer  *jwtSigner

	idMu  sync.Mutex
	idSeq uint64
}

// New creates a new Engine API client connecting over JSON-RPC with JWT
// authentication.
func New(_ context.Context, params ...Parameter) (*Service, error) {
	p, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	secret := p.jwtSecret
	if len(secret) == 0 {
		secret, err = loadJWTSecret(p.jwtSecretFile)
		if err != nil {
			return nil, err
		}
	}

	signer, err := newJWTSigner(secret)
	if err != nil {
		return nil, err
	}

	signer.id = p.jwtID
	signer.clientVersion = p.jwtClientVer

	client := p.client
	if client == nil {
		client = &http.Client{Timeout: p.timeout}
	}

	return &Service{
		log:     p.log.WithField("package", "github.com/ethpandaops/go-eth-engine-client/jsonrpc"),
		address: p.address,
		client:  client,
		timeout: p.timeout,
		signer:  signer,
	}, nil
}

// Address returns the address provided to the client.
func (s *Service) Address() string {
	return s.address
}

// nextID returns a monotonically-increasing JSON-RPC request id.
func (s *Service) nextID() uint64 {
	s.idMu.Lock()
	defer s.idMu.Unlock()

	s.idSeq++

	return s.idSeq
}
