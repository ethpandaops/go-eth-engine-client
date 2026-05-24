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

package paris

import (
	"fmt"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
)

// PayloadID is an 8-byte identifier returned by engine_forkchoiceUpdated to
// reference a payload build process.
type PayloadID [PayloadIDLength]byte

// String returns the lowercase 0x-prefixed hex representation.
func (p PayloadID) String() string {
	return fmt.Sprintf("%#x", p[:])
}

// MarshalJSON implements json.Marshaler.
func (p PayloadID) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, p[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadID) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(p[:], input, PayloadIDLength, "payloadId")
}
