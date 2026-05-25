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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
)

// BlobsBundle is the SSZ container `BlobsBundleV2` returned with
// engine_getPayloadV5. Proofs is no longer a single proof per blob but a
// flat list of cell proofs (CELLS_PER_EXT_BLOB per blob), supporting
// EIP-7594 PeerDAS.
type BlobsBundle struct {
	Commitments []cancun.KZGCommitment `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK"                   ssz-max:"4096"   ssz-size:"?,48"    json:"commitments"`
	Proofs      []cancun.KZGProof      `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK*CELLS_PER_EXT_BLOB" ssz-max:"524288" ssz-size:"?,48"    json:"proofs"`
	Blobs       []cancun.Blob          `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK"                   ssz-max:"4096"   ssz-size:"?,131072" json:"blobs"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobsBundle) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	type alias BlobsBundle

	return json.Marshal((*alias)(b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobsBundle) UnmarshalJSON(input []byte) error {
	type alias BlobsBundle

	if err := json.Unmarshal(input, (*alias)(b)); err != nil {
		return errors.Wrap(err, "BlobsBundle")
	}

	return nil
}

// String returns a JSON representation.
func (b *BlobsBundle) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
