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
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetBlobs fetches blobs from the EL blob pool all-or-nothing
// (engine_getBlobsV1 for cancun, V2 for osaka). The JSON-RPC result is a
// bare array of blob-and-proof objects (the body of the SSZ
// GetBlobsResponse container) and may be `null` when the EL is syncing, in
// which case the response's BlobsAndProofs is nil.
func (s *Service) GetBlobs(
	ctx context.Context,
	request *spec.VersionedGetBlobsRequest,
) (*spec.VersionedGetBlobsResponse, error) {
	if request == nil {
		return nil, errors.New("GetBlobs: nil request")
	}

	out := &spec.VersionedGetBlobsResponse{Version: request.Version}

	switch request.Version {
	case version.DataVersionCancun:
		if request.Cancun == nil {
			return nil, errors.New("GetBlobs: nil cancun request")
		}

		params := []any{hashList(request.Cancun.BlobVersionedHashes)}

		var blobs []*cancun.BlobAndProof
		if _, err := s.callOptional(ctx, "engine_getBlobsV1", params, &blobs); err != nil {
			return nil, err
		}

		out.Cancun = &cancun.GetBlobsResponse{BlobsAndProofs: blobs}
	case version.DataVersionOsaka:
		if request.Osaka == nil {
			return nil, errors.New("GetBlobs: nil osaka request")
		}

		params := []any{hashList(request.Osaka.BlobVersionedHashes)}

		var blobs []*osaka.BlobAndProof
		if _, err := s.callOptional(ctx, "engine_getBlobsV2", params, &blobs); err != nil {
			return nil, err
		}

		out.Osaka = &osaka.GetBlobsResponse{BlobsAndProofs: blobs}
	default:
		return nil, errors.Errorf("GetBlobs: unsupported version %s", request.Version)
	}

	return out, nil
}

// GetBlobsAgnostic is the fork-agnostic-typed variant of GetBlobs.
func (s *Service) GetBlobsAgnostic(
	ctx context.Context,
	request *all.GetBlobsRequest,
) (*all.GetBlobsResponse, error) {
	if request == nil {
		return nil, errors.New("GetBlobs: nil request")
	}

	versioned, err := request.ToVersioned()
	if err != nil {
		return nil, err
	}

	resp, err := s.GetBlobs(ctx, versioned)
	if err != nil {
		return nil, err
	}

	out := &all.GetBlobsResponse{}
	if err := out.FromVersioned(resp); err != nil {
		return nil, err
	}

	return out, nil
}

// GetBlobsPartial fetches blobs via engine_getBlobsV3 (osaka), the
// partial-response variant: a missing blob is a nil entry rather than a
// failed call. The JSON-RPC result is a bare array of (nullable)
// blob-and-proof objects and may itself be `null` when the EL is syncing.
func (s *Service) GetBlobsPartial(
	ctx context.Context,
	request *osaka.GetBlobsPartialRequest,
) (*osaka.GetBlobsPartialResponse, error) {
	if request == nil {
		return nil, errors.New("GetBlobsPartial: nil request")
	}

	params := []any{hashList(request.BlobVersionedHashes)}

	var blobs []*osaka.BlobAndProof
	if _, err := s.callOptional(ctx, "engine_getBlobsV3", params, &blobs); err != nil {
		return nil, err
	}

	return &osaka.GetBlobsPartialResponse{BlobsAndProofs: blobs}, nil
}

// GetBlobsCells fetches partial blob cell matrices via engine_getBlobsV4
// (amsterdam). The JSON-RPC result is a bare array of cell-and-proof
// objects (the body of the SSZ GetBlobsCellsResponse container) and may be
// `null` when the EL is syncing, in which case BlobsAndProofs is nil.
func (s *Service) GetBlobsCells(
	ctx context.Context,
	request *amsterdam.GetBlobsCellsRequest,
) (*amsterdam.GetBlobsCellsResponse, error) {
	if request == nil {
		return nil, errors.New("GetBlobsCells: nil request")
	}

	params := []any{hashList(request.VersionedBlobHashes), request.IndicesBitarray}

	var entries []*amsterdam.BlobCellsAndProofs
	if _, err := s.callOptional(ctx, "engine_getBlobsV4", params, &entries); err != nil {
		return nil, err
	}

	return &amsterdam.GetBlobsCellsResponse{BlobsAndProofs: entries}, nil
}
