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

package cancun

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetBlobsRequest carries the parameters of engine_getBlobsV1. It
// corresponds to the SSZ container `GetBlobsV1Request`.
type GetBlobsRequest struct {
	BlobVersionedHashes []paris.Hash32 `dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" ssz-size:"?,32"`
}

// GetBlobsResponse is the response from engine_getBlobsV1. It corresponds to
// the SSZ container `GetBlobsV1Response`. Entries are positionally aligned
// with the request's blob versioned hashes.
type GetBlobsResponse struct {
	BlobsAndProofs []*BlobAndProof `dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" json:"blobsAndProofs"`
}
