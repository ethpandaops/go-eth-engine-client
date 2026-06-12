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
	"encoding/json"
	"fmt"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
)

// BlobCell is a single EIP-7594 blob cell: 64 field elements (2048 bytes).
type BlobCell [CellLength]byte

// String returns the lowercase 0x-prefixed hex representation.
func (c BlobCell) String() string {
	return fmt.Sprintf("%#x", c[:])
}

// MarshalJSON implements json.Marshaler.
func (c BlobCell) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, c[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *BlobCell) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(c[:], input, CellLength, "blobCell")
}

// BlobCellsAndProofs is the per-blob entry returned by engine_getBlobsV4 —
// the partial matrix of cells and their KZG proofs for a single blob. Cell
// and proof entries may individually be null when the corresponding cell is
// unavailable.
//
// The PR #764 SSZ transport spec does not define a container for this type;
// the SSZ schema here is defined following the same nullable-list
// convention used elsewhere in the spec:
//
//	class BlobCellsAndProofsV1(Container):
//	    blob_cells: List[List[ByteVector[CELL_LENGTH], 1], CELLS_PER_EXT_BLOB]
//	    proofs:     List[List[Bytes48, 1], CELLS_PER_EXT_BLOB]
type BlobCellsAndProofs struct {
	BlobCells []*BlobCell        `ssz-type:"list,optional-list" dynssz-max:"CELLS_PER_EXT_BLOB" ssz-max:"128" json:"blob_cells"`
	Proofs    []*cancun.KZGProof `ssz-type:"list,optional-list" dynssz-max:"CELLS_PER_EXT_BLOB" ssz-max:"128" json:"proofs"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobCellsAndProofs) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	type alias BlobCellsAndProofs

	return json.Marshal((*alias)(b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobCellsAndProofs) UnmarshalJSON(input []byte) error {
	type alias BlobCellsAndProofs

	if err := json.Unmarshal(input, (*alias)(b)); err != nil {
		return fmt.Errorf("BlobCellsAndProofs: %w", err)
	}

	return nil
}

// String returns a JSON representation.
func (b *BlobCellsAndProofs) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
