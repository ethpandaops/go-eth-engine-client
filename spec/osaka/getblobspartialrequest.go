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

package osaka

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetBlobsPartialRequest carries the parameters of engine_getBlobsV3 — the
// variant that supports partial responses (per-blob nulls when individual
// blobs are unavailable). It corresponds to the SSZ container
// `GetBlobsV3Request`.
//
// Use [GetBlobsRequest] (engine_getBlobsV2) for the all-or-nothing query
// style; use this type when the caller can tolerate missing entries.
type GetBlobsPartialRequest struct {
	BlobVersionedHashes []paris.Hash32 `dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" ssz-size:"?,32"`
}
