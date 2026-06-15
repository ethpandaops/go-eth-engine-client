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
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// NewPayload submits an execution payload for validation via
// POST /engine/v2/{fork}/payloads with an SSZ-encoded ExecutionPayloadEnvelope
// body. Returns a [paris.PayloadStatus] decoded from the SSZ response body.
//
// TODO: implement once the spec/* containers reflect the Marius changes
// (ExpectedBlobVersionedHashes removed from NewPayloadRequest;
// PayloadStatus.ValidationError moved to Optional[String]).
func (s *Service) NewPayload(
	_ context.Context,
	_ *spec.VersionedNewPayloadRequest,
) (*paris.PayloadStatus, error) {
	return nil, ErrNotImplemented
}

// NewPayloadAgnostic is the fork-agnostic-typed variant of NewPayload.
func (s *Service) NewPayloadAgnostic(
	_ context.Context,
	_ *all.NewPayloadRequest,
) (*paris.PayloadStatus, error) {
	return nil, ErrNotImplemented
}
