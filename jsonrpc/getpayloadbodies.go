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

package jsonrpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetPayloadBodiesByHash retrieves execution payload bodies for the given
// block hashes (engine_getPayloadBodiesByHashV1 for shanghai..osaka, V2 for
// amsterdam). Entries are nil for blocks the EL does not have.
func (s *Service) GetPayloadBodiesByHash(
	ctx context.Context,
	dataVersion version.DataVersion,
	blockHashes []paris.Hash32,
) ([]*spec.VersionedExecutionPayloadBody, error) {
	method, v2, err := payloadBodiesByHashMethod(dataVersion)
	if err != nil {
		return nil, err
	}

	return s.fetchPayloadBodies(ctx, method, dataVersion, v2, []any{hashList(blockHashes)})
}

// GetPayloadBodiesByHashAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetPayloadBodiesByHashAgnostic(
	ctx context.Context,
	dataVersion version.DataVersion,
	blockHashes []paris.Hash32,
) ([]*all.ExecutionPayloadBody, error) {
	versioned, err := s.GetPayloadBodiesByHash(ctx, dataVersion, blockHashes)
	if err != nil {
		return nil, err
	}

	return payloadBodiesToAgnostic(versioned)
}

// GetPayloadBodiesByRange retrieves execution payload bodies for a range of
// block numbers (engine_getPayloadBodiesByRangeV1 for shanghai..osaka, V2
// for amsterdam). Entries are nil for blocks the EL does not have.
func (s *Service) GetPayloadBodiesByRange(
	ctx context.Context,
	dataVersion version.DataVersion,
	start uint64,
	count uint64,
) ([]*spec.VersionedExecutionPayloadBody, error) {
	method, v2, err := payloadBodiesByRangeMethod(dataVersion)
	if err != nil {
		return nil, err
	}

	params := []any{jsonhex.QuantityU64(start), jsonhex.QuantityU64(count)}

	return s.fetchPayloadBodies(ctx, method, dataVersion, v2, params)
}

// GetPayloadBodiesByRangeAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetPayloadBodiesByRangeAgnostic(
	ctx context.Context,
	dataVersion version.DataVersion,
	start uint64,
	count uint64,
) ([]*all.ExecutionPayloadBody, error) {
	versioned, err := s.GetPayloadBodiesByRange(ctx, dataVersion, start, count)
	if err != nil {
		return nil, err
	}

	return payloadBodiesToAgnostic(versioned)
}

// fetchPayloadBodies issues the call and wraps each body in a
// VersionedExecutionPayloadBody for dataVersion. The JSON-RPC result is a
// bare array of (nullable) body objects.
func (s *Service) fetchPayloadBodies(
	ctx context.Context,
	method string,
	dataVersion version.DataVersion,
	v2 bool,
	params []any,
) ([]*spec.VersionedExecutionPayloadBody, error) {
	if v2 {
		var bodies []*amsterdam.ExecutionPayloadBody
		if _, err := s.callOptional(ctx, method, params, &bodies); err != nil {
			return nil, err
		}

		out := make([]*spec.VersionedExecutionPayloadBody, len(bodies))
		for i, b := range bodies {
			if b != nil {
				out[i] = &spec.VersionedExecutionPayloadBody{Version: dataVersion, Amsterdam: b}
			}
		}

		return out, nil
	}

	var bodies []*shanghai.ExecutionPayloadBody
	if _, err := s.callOptional(ctx, method, params, &bodies); err != nil {
		return nil, err
	}

	out := make([]*spec.VersionedExecutionPayloadBody, len(bodies))
	for i, b := range bodies {
		if b != nil {
			out[i] = wrapV1PayloadBody(dataVersion, b)
		}
	}

	return out, nil
}

// wrapV1PayloadBody stores a V1 (shanghai-shaped) body in the
// VersionedExecutionPayloadBody field matching dataVersion.
func wrapV1PayloadBody(
	dataVersion version.DataVersion,
	body *shanghai.ExecutionPayloadBody,
) *spec.VersionedExecutionPayloadBody {
	out := &spec.VersionedExecutionPayloadBody{Version: dataVersion}

	switch dataVersion {
	case version.DataVersionCancun:
		out.Cancun = body
	case version.DataVersionPrague:
		out.Prague = body
	case version.DataVersionOsaka:
		out.Osaka = body
	default:
		out.Shanghai = body
	}

	return out
}

// payloadBodiesToAgnostic converts versioned bodies into union bodies.
func payloadBodiesToAgnostic(
	versioned []*spec.VersionedExecutionPayloadBody,
) ([]*all.ExecutionPayloadBody, error) {
	out := make([]*all.ExecutionPayloadBody, len(versioned))

	for i, v := range versioned {
		if v == nil {
			continue
		}

		body := &all.ExecutionPayloadBody{}
		if err := body.FromVersioned(v); err != nil {
			return nil, err
		}

		out[i] = body
	}

	return out, nil
}

// payloadBodiesByHashMethod returns the engine_getPayloadBodiesByHash method
// name for the version and whether it is the V2 (amsterdam) variant.
func payloadBodiesByHashMethod(v version.DataVersion) (string, bool, error) {
	switch v {
	case version.DataVersionShanghai,
		version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka:
		return "engine_getPayloadBodiesByHashV1", false, nil
	case version.DataVersionAmsterdam:
		return "engine_getPayloadBodiesByHashV2", true, nil
	default:
		return "", false, errors.Errorf("GetPayloadBodiesByHash: unsupported version %s", v)
	}
}

// payloadBodiesByRangeMethod mirrors payloadBodiesByHashMethod for the
// ByRange family.
func payloadBodiesByRangeMethod(v version.DataVersion) (string, bool, error) {
	switch v {
	case version.DataVersionShanghai,
		version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka:
		return "engine_getPayloadBodiesByRangeV1", false, nil
	case version.DataVersionAmsterdam:
		return "engine_getPayloadBodiesByRangeV2", true, nil
	default:
		return "", false, errors.Errorf("GetPayloadBodiesByRange: unsupported version %s", v)
	}
}
