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
	"context"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// fcuParams is the resolved set of forkchoice-update parameters extracted
// from a versioned request: the (constant-shape) forkchoice state and the
// optional, fork-specific payload attributes.
type fcuParams struct {
	state *paris.ForkchoiceState
	attrs any
}

// ForkchoiceUpdated updates the fork choice and optionally begins a payload
// build, dispatching to the engine_forkchoiceUpdated method version that
// matches request.Version.
//
//nolint:gocyclo // Per-version dispatch is an inherently wide switch.
func (s *Service) ForkchoiceUpdated(
	ctx context.Context,
	request *spec.VersionedForkchoiceUpdatedRequest,
) (*paris.ForkchoiceUpdatedResponse, error) {
	if request == nil {
		return nil, errors.New("ForkchoiceUpdated: nil request")
	}

	var (
		method string
		fcu    fcuParams
		err    error
	)

	switch request.Version {
	case version.DataVersionParis:
		method = "engine_forkchoiceUpdatedV1"
		if request.Paris == nil {
			return nil, errors.New("ForkchoiceUpdated: nil paris request")
		}

		fcu.state = request.Paris.ForkchoiceState
		if request.Paris.PayloadAttributes != nil {
			fcu.attrs = request.Paris.PayloadAttributes
		}
	case version.DataVersionShanghai:
		method = "engine_forkchoiceUpdatedV2"
		if request.Shanghai == nil {
			return nil, errors.New("ForkchoiceUpdated: nil shanghai request")
		}

		fcu.state = request.Shanghai.ForkchoiceState
		if request.Shanghai.PayloadAttributes != nil {
			fcu.attrs = request.Shanghai.PayloadAttributes
		}
	case version.DataVersionCancun, version.DataVersionPrague, version.DataVersionOsaka:
		method = "engine_forkchoiceUpdatedV3"

		fcu, err = cancunForkchoiceParams(request)
		if err != nil {
			return nil, err
		}
	case version.DataVersionAmsterdam:
		method = "engine_forkchoiceUpdatedV4"
		if request.Amsterdam == nil {
			return nil, errors.New("ForkchoiceUpdated: nil amsterdam request")
		}

		fcu.state = request.Amsterdam.ForkchoiceState
		if request.Amsterdam.PayloadAttributes != nil {
			fcu.attrs = request.Amsterdam.PayloadAttributes
		}
	default:
		return nil, errors.Errorf("ForkchoiceUpdated: unsupported version %s", request.Version)
	}

	if fcu.state == nil {
		return nil, errors.New("ForkchoiceUpdated: nil forkchoice state")
	}

	params := []any{fcu.state, fcu.attrs}

	// engine_forkchoiceUpdatedV4 takes an additional custodyColumns
	// parameter (null when the CL provides no custody set).
	if request.Version == version.DataVersionAmsterdam {
		var custody any
		if request.CustodyColumns != nil {
			custody = request.CustodyColumns
		}

		params = append(params, custody)
	}

	response := &paris.ForkchoiceUpdatedResponse{}
	if err := s.call(ctx, method, params, response); err != nil {
		return nil, err
	}

	return response, nil
}

// ForkchoiceUpdatedAgnostic is the fork-agnostic-typed variant of
// ForkchoiceUpdated.
func (s *Service) ForkchoiceUpdatedAgnostic(
	ctx context.Context,
	request *all.ForkchoiceUpdatedRequest,
) (*paris.ForkchoiceUpdatedResponse, error) {
	if request == nil {
		return nil, errors.New("ForkchoiceUpdated: nil request")
	}

	versioned, err := request.ToVersioned()
	if err != nil {
		return nil, err
	}

	return s.ForkchoiceUpdated(ctx, versioned)
}

// cancunForkchoiceParams extracts the V3 (cancun/prague/osaka) forkchoice
// parameters, which all share cancun's request shape.
func cancunForkchoiceParams(request *spec.VersionedForkchoiceUpdatedRequest) (fcuParams, error) {
	req := request.Cancun

	switch request.Version {
	case version.DataVersionPrague:
		req = request.Prague
	case version.DataVersionOsaka:
		req = request.Osaka
	}

	if req == nil {
		return fcuParams{}, errors.Errorf("ForkchoiceUpdated: nil %s request", request.Version)
	}

	out := fcuParams{state: req.ForkchoiceState}
	if req.PayloadAttributes != nil {
		out.attrs = req.PayloadAttributes
	}

	return out, nil
}
