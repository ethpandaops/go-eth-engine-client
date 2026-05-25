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
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// BlockAccessList is the RLP-encoded EIP-7928 block access list. The
// Engine API treats it as opaque bytes; the client does not parse it.
type BlockAccessList []byte

// String returns the lowercase 0x-prefixed hex representation.
func (b BlockAccessList) String() string {
	return fmt.Sprintf("%#x", []byte(b))
}

// MarshalJSON implements json.Marshaler.
func (b BlockAccessList) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte(`"0x"`), nil
	}

	return fmt.Appendf(nil, `"%#x"`, []byte(b)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockAccessList) UnmarshalJSON(input []byte) error {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return errors.New("blockAccessList: not a JSON string")
	}

	body := input[1 : len(input)-1]
	if len(body) < 2 || body[0] != '0' || (body[1] != 'x' && body[1] != 'X') {
		return errors.New("blockAccessList: missing 0x prefix")
	}

	hexBody := body[2:]
	if len(hexBody) == 0 {
		*b = nil

		return nil
	}

	decoded := make([]byte, hex.DecodedLen(len(hexBody)))

	n, err := hex.Decode(decoded, hexBody)
	if err != nil {
		return errors.Wrap(err, "blockAccessList: invalid hex")
	}

	*b = decoded[:n]

	return nil
}
