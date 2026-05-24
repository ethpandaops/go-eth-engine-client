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

// Bloom is a 256-byte logs bloom filter as embedded in an ExecutionPayload.
type Bloom [BloomLength]byte

// String returns the lowercase 0x-prefixed hex representation.
func (b Bloom) String() string {
	return fmt.Sprintf("%#x", b[:])
}

// MarshalJSON implements json.Marshaler.
func (b Bloom) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, b[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bloom) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(b[:], input, BloomLength, "bloom")
}
