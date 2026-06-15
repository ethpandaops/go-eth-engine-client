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

// Package rest is a scaffold for the REST + SSZ Engine API transport
// described in execution-apis PR #793 ("Marius spec"). The wire format moves
// from JSON-RPC over HTTP/1.1 to fork-scoped REST endpoints under
// `/engine/v2/{fork}/...` with SSZ request and response bodies over
// HTTP/2 cleartext (h2c). JWT (HS256) authentication remains, but the `clv`
// claim is dropped in favour of an `X-Engine-Client-Version` request header.
//
// At present every provider method returns [ErrNotImplemented]: the spec is
// still a draft (MAX_* constants and naming are flagged TBD), and several
// shared SSZ containers (PayloadStatus.ValidationError as Optional[String],
// BodyEntry / BlobEntry wrappers, fork-folded custody columns, the removal
// of ExpectedBlobVersionedHashes from the new-payload envelope, BuiltPayload
// field reordering) need to land in the `spec/*` packages before the wire
// implementations are useful.
//
// For the legacy JSON-RPC transport, see the sibling `jsonrpc` package.
package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	engine "github.com/ethpandaops/go-eth-engine-client"
)

// Compile-time assertions that the Service satisfies the engine interfaces.
// The provider methods are stubs returning ErrNotImplemented; the assertions
// pin the surface so the type checker flags interface drift while the
// implementation catches up.
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

// DefaultBasePath is the path prefix mandated by the Marius spec for all
// hot-path endpoints. Diagnostic endpoints (/identity, /capabilities) live
// directly under this prefix without a fork segment.
const DefaultBasePath = "/engine/v2"

// Service is a REST + SSZ Engine API client. It currently implements the
// engine provider surface as ErrNotImplemented stubs.
type Service struct {
	log      logrus.FieldLogger
	address  string
	basePath string
	client   *http.Client
	timeout  time.Duration
	signer   *jwtSigner
	clientID string
}

// New creates a new REST + SSZ Engine API client.
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

	client := p.client
	if client == nil {
		client = newDefaultHTTPClient(p.timeout)
	}

	return &Service{
		log:      p.log.WithField("package", "github.com/ethpandaops/go-eth-engine-client/rest"),
		address:  p.address,
		basePath: p.basePath,
		client:   client,
		timeout:  p.timeout,
		signer:   signer,
		clientID: p.clientVersionHeader,
	}, nil
}

// Address returns the address provided to the client.
func (s *Service) Address() string {
	return s.address
}
