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

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// NewPayload submits an execution payload for validation, dispatching to the
// engine_newPayload method version that matches request.Version.
//
//nolint:gocyclo // Per-version dispatch is an inherently wide switch.
func (s *Service) NewPayload(
	ctx context.Context,
	request *spec.VersionedNewPayloadRequest,
) (*paris.PayloadStatus, error) {
	if request == nil {
		return nil, errors.New("NewPayload: nil request")
	}

	var (
		method string
		params []any
	)

	switch request.Version {
	case version.DataVersionParis:
		if request.Paris == nil || request.Paris.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil paris request")
		}

		method, params = "engine_newPayloadV1", []any{request.Paris.ExecutionPayload}
	case version.DataVersionShanghai:
		if request.Shanghai == nil || request.Shanghai.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil shanghai request")
		}

		method, params = "engine_newPayloadV2", []any{request.Shanghai.ExecutionPayload}
	case version.DataVersionCancun:
		if request.Cancun == nil || request.Cancun.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil cancun request")
		}

		r := request.Cancun
		method, params = "engine_newPayloadV3", []any{
			r.ExecutionPayload,
			hashList(r.ExpectedBlobVersionedHashes),
			r.ParentBeaconBlockRoot,
		}
	case version.DataVersionPrague:
		if request.Prague == nil || request.Prague.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil prague request")
		}

		method, params = "engine_newPayloadV4", newPayloadV4Params(request.Prague)
	case version.DataVersionOsaka:
		if request.Osaka == nil || request.Osaka.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil osaka request")
		}

		// Osaka reuses prague's V4 newPayload request shape.
		method, params = "engine_newPayloadV4", newPayloadV4Params(request.Osaka)
	case version.DataVersionAmsterdam:
		if request.Amsterdam == nil || request.Amsterdam.ExecutionPayload == nil {
			return nil, errors.New("NewPayload: nil amsterdam request")
		}

		r := request.Amsterdam
		method, params = "engine_newPayloadV5", []any{
			r.ExecutionPayload,
			hashList(r.ExpectedBlobVersionedHashes),
			r.ParentBeaconBlockRoot,
			executionRequestList(r.ExecutionRequests),
		}
	default:
		return nil, errors.Errorf("NewPayload: unsupported version %s", request.Version)
	}

	status := &paris.PayloadStatus{}
	if err := s.call(ctx, method, params, status); err != nil {
		return nil, err
	}

	return status, nil
}

// NewPayloadAgnostic is the fork-agnostic-typed variant of NewPayload.
func (s *Service) NewPayloadAgnostic(
	ctx context.Context,
	request *all.NewPayloadRequest,
) (*paris.PayloadStatus, error) {
	if request == nil {
		return nil, errors.New("NewPayload: nil request")
	}

	versioned, err := request.ToVersioned()
	if err != nil {
		return nil, err
	}

	return s.NewPayload(ctx, versioned)
}

// newPayloadV4Params builds the positional params for engine_newPayloadV4.
func newPayloadV4Params(r *prague.NewPayloadRequest) []any {
	return []any{
		r.ExecutionPayload,
		hashList(r.ExpectedBlobVersionedHashes),
		r.ParentBeaconBlockRoot,
		executionRequestList(r.ExecutionRequests),
	}
}

// hashList coerces a nil slice to a non-nil empty slice so it JSON-encodes as
// `[]` rather than `null` (the Engine API expects an array).
func hashList(in []paris.Hash32) []paris.Hash32 {
	if in == nil {
		return []paris.Hash32{}
	}

	return in
}

// executionRequestList coerces a nil slice to a non-nil empty slice.
func executionRequestList(in []prague.ExecutionRequest) []prague.ExecutionRequest {
	if in == nil {
		return []prague.ExecutionRequest{}
	}

	return in
}
