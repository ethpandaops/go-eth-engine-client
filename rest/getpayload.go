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
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetPayload retrieves a previously-initiated payload build via
// GET /engine/v2/{fork}/payloads/{payloadId}. Unlike the JSON-RPC path,
// paris is supported because the Marius spec defines a `BuiltPayload`
// response container for every fork.
//
// TODO: implement once BuiltPayload field order in spec/*/getpayloadresponse.go
// matches Marius ({payload, block_value, blobs_bundle, execution_requests,
// should_override_builder}).
func (s *Service) GetPayload(
	_ context.Context,
	_ version.DataVersion,
	_ paris.PayloadID,
) (*spec.VersionedGetPayloadResponse, error) {
	return nil, ErrNotImplemented
}

// GetPayloadAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetPayloadAgnostic(
	_ context.Context,
	_ version.DataVersion,
	_ paris.PayloadID,
) (*all.GetPayloadResponse, error) {
	return nil, ErrNotImplemented
}
