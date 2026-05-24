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

// Address is a 20-byte execution-layer account address. JSON marshaling
// emits the lowercase 0x-prefixed form. EIP-55 mixed-case input is accepted
// on unmarshal but is not produced on marshal.
type Address [AddressLength]byte

// String returns the lowercase 0x-prefixed hex representation.
func (a Address) String() string {
	return fmt.Sprintf("%#x", a[:])
}

// Format formats the address.
func (a Address) Format(state fmt.State, v rune) {
	format := string(v)

	switch v {
	case 's':
		fmt.Fprint(state, a.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}

		fmt.Fprintf(state, "%"+format, a[:])
	default:
		fmt.Fprintf(state, "%"+format, a[:])
	}
}

// MarshalJSON implements json.Marshaler.
func (a Address) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, a[:]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Address) UnmarshalJSON(input []byte) error {
	return jsonhex.DecodeFixed(a[:], input, AddressLength, "address")
}

// MarshalText implements encoding.TextMarshaler.
func (a Address) MarshalText() ([]byte, error) {
	return jsonhex.EncodeFixedText(a[:]), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (a *Address) UnmarshalText(input []byte) error {
	return jsonhex.DecodeFixedText(a[:], input, AddressLength, "address")
}
