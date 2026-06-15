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

// GetPayloadBodiesByHash retrieves execution payload bodies for the given
// block hashes via POST /engine/v2/{fork}/bodies/hash. The {fork} segment
// selects both the response schema and the era of returned blocks:
// out-of-era requests come back as `available=false` per entry.
//
// TODO: implement once each fork's ExecutionPayloadBody has a BodyEntry
// wrapper (`available + body`) and the per-fork bodies response container
// is added.
func (s *Service) GetPayloadBodiesByHash(
	_ context.Context,
	_ version.DataVersion,
	_ []paris.Hash32,
) ([]*spec.VersionedExecutionPayloadBody, error) {
	return nil, ErrNotImplemented
}

// GetPayloadBodiesByHashAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetPayloadBodiesByHashAgnostic(
	_ context.Context,
	_ version.DataVersion,
	_ []paris.Hash32,
) ([]*all.ExecutionPayloadBody, error) {
	return nil, ErrNotImplemented
}

// GetPayloadBodiesByRange retrieves bodies for a block-number range via
// GET /engine/v2/{fork}/bodies?from=N&count=M. Out-of-range entries come
// back with `available=false`.
//
// TODO: implement once the per-fork bodies response container lands.
func (s *Service) GetPayloadBodiesByRange(
	_ context.Context,
	_ version.DataVersion,
	_ uint64,
	_ uint64,
) ([]*spec.VersionedExecutionPayloadBody, error) {
	return nil, ErrNotImplemented
}

// GetPayloadBodiesByRangeAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetPayloadBodiesByRangeAgnostic(
	_ context.Context,
	_ version.DataVersion,
	_ uint64,
	_ uint64,
) ([]*all.ExecutionPayloadBody, error) {
	return nil, ErrNotImplemented
}
