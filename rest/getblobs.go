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
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
)

// GetBlobs fetches blobs all-or-nothing via POST /engine/v2/blobs/v1
// (cancun) or /v2 (osaka). The Marius response wraps each entry in a
// `BlobV{1,2}Entry { available, contents }` container; an HTTP 204 signals
// "EL cannot serve / all-or-nothing miss" and surfaces as a nil response
// slice.
//
// TODO: implement once cancun/osaka GetBlobsResponse adopts the
// BlobEntry-wrapped shape.
func (s *Service) GetBlobs(
	_ context.Context,
	_ *spec.VersionedGetBlobsRequest,
) (*spec.VersionedGetBlobsResponse, error) {
	return nil, ErrNotImplemented
}

// GetBlobsAgnostic is the fork-agnostic-typed variant.
func (s *Service) GetBlobsAgnostic(
	_ context.Context,
	_ *all.GetBlobsRequest,
) (*all.GetBlobsResponse, error) {
	return nil, ErrNotImplemented
}

// GetBlobsPartial fetches blobs with per-entry availability via
// POST /engine/v2/blobs/v3 (osaka). Missing blobs surface as
// `available=false` per entry.
//
// TODO: implement once osaka.GetBlobsPartialResponse adopts the
// BlobV2Entry wrapper (current shape uses the legacy nullable-list trick).
func (s *Service) GetBlobsPartial(
	_ context.Context,
	_ *osaka.GetBlobsPartialRequest,
) (*osaka.GetBlobsPartialResponse, error) {
	return nil, ErrNotImplemented
}

// GetBlobsCells fetches per-cell blob matrices via POST /engine/v2/blobs/v4
// (amsterdam). Cell-level availability is encoded as Optional[T] inside the
// per-blob BlobCellsAndProofs container.
//
// TODO: implement once amsterdam.GetBlobsCellsResponse adopts the
// BlobV4Entry wrapper.
func (s *Service) GetBlobsCells(
	_ context.Context,
	_ *amsterdam.GetBlobsCellsRequest,
) (*amsterdam.GetBlobsCellsResponse, error) {
	return nil, ErrNotImplemented
}
