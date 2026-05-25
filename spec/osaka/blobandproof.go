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

// BlobAndProof is the SSZ container `BlobAndProofV2` returned by
// engine_getBlobsV2 / V3. Each blob carries CELLS_PER_EXT_BLOB cell proofs.
type BlobAndProof struct {
	Blob   cancun.Blob       `ssz-size:"131072"               json:"blob"`
	Proofs []cancun.KZGProof `dynssz-max:"CELLS_PER_EXT_BLOB" ssz-max:"128" ssz-size:"?,48" json:"proofs"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobAndProof) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	type alias BlobAndProof

	return json.Marshal((*alias)(b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobAndProof) UnmarshalJSON(input []byte) error {
	type alias BlobAndProof

	if err := json.Unmarshal(input, (*alias)(b)); err != nil {
		return errors.Wrap(err, "BlobAndProof")
	}

	return nil
}

// String returns a JSON representation.
func (b *BlobAndProof) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
