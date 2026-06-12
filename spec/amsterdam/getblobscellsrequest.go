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

package amsterdam

import (
	"fmt"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// IndicesBitarray is the 16-byte bit-array selecting which cell indices to
// retrieve in an engine_getBlobsV4 request (CELLS_PER_EXT_BLOB / 8 = 16).
type IndicesBitarray [16]byte

// MarshalJSON implements json.Marshaler.
func (b IndicesBitarray) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, b[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *IndicesBitarray) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(b[:], input, 16, "indicesBitarray")
}

// GetBlobsCellsRequest carries the parameters of engine_getBlobsV4 — the
// cell-based variant that fetches individual blob cells (rather than full
// blobs) and supports partial responses.
//
// The PR #764 SSZ transport spec does not define a container for this
// method; the SSZ schema here follows the same convention as the other
// getBlobs requests.
type GetBlobsCellsRequest struct {
	VersionedBlobHashes []paris.Hash32  `dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" ssz-size:"?,32" json:"versioned_blob_hashes"`
	IndicesBitarray     IndicesBitarray `ssz-size:"16"                                                      json:"indices_bitarray"`
}
