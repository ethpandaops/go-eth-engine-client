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

package prague

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// ExecutionRequest is an opaque EIP-7685 execution-layer triggered request.
// The first byte is the request type; the remainder is the request payload.
// The Engine API client does not parse the contents.
type ExecutionRequest []byte

// String returns the lowercase 0x-prefixed hex representation.
func (r ExecutionRequest) String() string {
	return fmt.Sprintf("%#x", []byte(r))
}

// MarshalJSON implements json.Marshaler.
func (r ExecutionRequest) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%#x"`, []byte(r)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *ExecutionRequest) UnmarshalJSON(input []byte) error {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return errors.New("executionRequest: not a JSON string")
	}

	body := input[1 : len(input)-1]
	if len(body) < 2 || body[0] != '0' || (body[1] != 'x' && body[1] != 'X') {
		return errors.New("executionRequest: missing 0x prefix")
	}

	hexBody := body[2:]
	decoded := make([]byte, hex.DecodedLen(len(hexBody)))

	n, err := hex.Decode(decoded, hexBody)
	if err != nil {
		return errors.Wrap(err, "executionRequest: invalid hex")
	}

	*r = decoded[:n]

	return nil
}
