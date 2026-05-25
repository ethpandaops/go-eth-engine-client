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
)

// CustodyColumns is the 16-byte bit-array passed as the 3rd parameter to
// engine_forkchoiceUpdatedV4. Each bit indicates whether the calling CL
// custodies the corresponding cell index (CELLS_PER_EXT_BLOB / 8 = 16).
// Sent over JSON-RPC only — the PR764 SSZ transport's
// `ForkchoiceUpdatedV4Request` container does not include this field.
type CustodyColumns [CustodyColumnsLength]byte

// String returns the lowercase 0x-prefixed hex representation.
func (c CustodyColumns) String() string {
	return fmt.Sprintf("%#x", c[:])
}

// MarshalJSON implements json.Marshaler.
func (c CustodyColumns) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, c[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *CustodyColumns) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(c[:], input, CustodyColumnsLength, "custodyColumns")
}
