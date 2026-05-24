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

// Hash32 is a 32-byte hash used throughout the Engine API for parent /
// state / receipts / block / prev-randao / parent-beacon-block roots.
type Hash32 [Hash32Length]byte

// String returns the lowercase 0x-prefixed hex representation of the hash.
func (h Hash32) String() string {
	return fmt.Sprintf("%#x", h[:])
}

// Format formats the hash.
func (h Hash32) Format(state fmt.State, v rune) {
	format := string(v)

	switch v {
	case 's':
		fmt.Fprint(state, h.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}

		fmt.Fprintf(state, "%"+format, h[:])
	default:
		fmt.Fprintf(state, "%"+format, h[:])
	}
}

// MarshalJSON implements json.Marshaler.
func (h Hash32) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, h[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *Hash32) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(h[:], input, Hash32Length, "hash32")
}

// MarshalText implements encoding.TextMarshaler so Hash32 can be used as a
// map key in JSON.
func (h Hash32) MarshalText() ([]byte, error) {
	return jsonhex.EncodeFixedText(h[:]), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *Hash32) UnmarshalText(input []byte) error {
	return jsonhex.DecodeFixedText(h[:], input, Hash32Length, "hash32")
}
