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

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
)

// BlobCell is the byte payload of a single EIP-7594 cell (one of
// CELLS_PER_EXT_BLOB per extended blob). The Engine API treats it as opaque
// bytes; the client does not parse it.
type BlobCell []byte

// String returns the lowercase 0x-prefixed hex representation.
func (c BlobCell) String() string {
	return fmt.Sprintf("%#x", []byte(c))
}

// MarshalJSON implements json.Marshaler.
func (c BlobCell) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}

	if len(c) == 0 {
		return []byte(`"0x"`), nil
	}

	return fmt.Appendf(nil, `"%#x"`, []byte(c)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *BlobCell) UnmarshalJSON(input []byte) error {
	if string(input) == "null" {
		*c = nil

		return nil
	}

	type alias []byte

	v := alias(*c)
	if err := json.Unmarshal(input, &v); err != nil {
		return errors.Wrap(err, "BlobCell")
	}

	*c = BlobCell(v)

	return nil
}

// BlobCellsAndProofs is the per-blob entry returned by
// engine_getBlobsV4. Cell entries and proofs may individually be `null`
// (in JSON) when the corresponding cell is missing.
//
// This type has no defined SSZ encoding in the PR #764 SSZ transport spec;
// engine_getBlobsV4 is JSON-RPC only at the time of writing.
type BlobCellsAndProofs struct {
	BlobCells []*BlobCell        `json:"blob_cells"`
	Proofs    []*cancun.KZGProof `json:"proofs"`
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
		return errors.Wrap(err, "BlobCellsAndProofs")
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
