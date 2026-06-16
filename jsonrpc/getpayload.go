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
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetPayload retrieves a previously-initiated payload build, dispatching to
// the engine_getPayload method version implied by dataVersion. The result is
// returned as the memory-efficient spec.VersionedGetPayloadResponse.
//
// engine_getPayloadV1 (paris) returns a bare ExecutionPayload with no
// response container and is therefore not supported here.
func (s *Service) GetPayload(
	ctx context.Context,
	dataVersion version.DataVersion,
	payloadID paris.PayloadID,
) (*spec.VersionedGetPayloadResponse, error) {
	out := &spec.VersionedGetPayloadResponse{Version: dataVersion}
	params := []any{payloadID}

	var (
		method string
		target any
	)

	switch dataVersion {
	case version.DataVersionShanghai:
		out.Shanghai = &shanghai.GetPayloadResponse{}
		method, target = "engine_getPayloadV2", out.Shanghai
	case version.DataVersionCancun:
		out.Cancun = &cancun.GetPayloadResponse{}
		method, target = "engine_getPayloadV3", out.Cancun
	case version.DataVersionPrague:
		out.Prague = &prague.GetPayloadResponse{}
		method, target = "engine_getPayloadV4", out.Prague
	case version.DataVersionOsaka:
		out.Osaka = &osaka.GetPayloadResponse{}
		method, target = "engine_getPayloadV5", out.Osaka
	case version.DataVersionAmsterdam:
		out.Amsterdam = &amsterdam.GetPayloadResponse{}
		method, target = "engine_getPayloadV6", out.Amsterdam
	case version.DataVersionBogota:
		// Bogota reuses Amsterdam's V6 response container.
		out.Bogota = &amsterdam.GetPayloadResponse{}
		method, target = "engine_getPayloadV6", out.Bogota
	default:
		return nil, errors.Errorf(
			"GetPayload: unsupported version %s (paris V1 returns a bare payload)",
			dataVersion,
		)
	}

	if err := s.call(ctx, method, params, target); err != nil {
		return nil, err
	}

	return out, nil
}

// GetPayloadAgnostic is the fork-agnostic-typed variant of GetPayload.
func (s *Service) GetPayloadAgnostic(
	ctx context.Context,
	dataVersion version.DataVersion,
	payloadID paris.PayloadID,
) (*all.GetPayloadResponse, error) {
	versioned, err := s.GetPayload(ctx, dataVersion, payloadID)
	if err != nil {
		return nil, err
	}

	out := &all.GetPayloadResponse{}
	if err := out.FromVersioned(versioned); err != nil {
		return nil, err
	}

	return out, nil
}
