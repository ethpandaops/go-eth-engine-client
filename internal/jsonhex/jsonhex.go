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

// Package jsonhex provides JSON marshaling helpers for the Engine API's
// 0x-prefixed hex wire format: `QUANTITY` (uint64 / uint256) and variable
// `bytes` strings. The types live behind module-internal exports so each
// fork package can reuse them without leaking helpers into the public API.
package jsonhex

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

// QuantityU64 marshals as the JSON-RPC `QUANTITY` form of a uint64:
// `"0x..."` with no leading zeros, or `"0x0"` for the zero value.
type QuantityU64 uint64

// MarshalJSON implements json.Marshaler.
func (q QuantityU64) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"0x%x"`, uint64(q)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (q *QuantityU64) UnmarshalJSON(input []byte) error {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return errors.New("quantity: not a JSON string")
	}

	body := input[1 : len(input)-1]
	if len(body) < 3 || body[0] != '0' || (body[1] != 'x' && body[1] != 'X') {
		return errors.New("quantity: missing 0x prefix")
	}

	v, err := strconv.ParseUint(string(body[2:]), 16, 64)
	if err != nil {
		return errors.Wrap(err, "quantity")
	}

	*q = QuantityU64(v)

	return nil
}

// QuantityU256 marshals as the JSON-RPC `QUANTITY` form of a 256-bit
// unsigned integer. A nil pointer marshals as `null`.
type QuantityU256 uint256.Int

// MarshalJSON implements json.Marshaler.
func (q *QuantityU256) MarshalJSON() ([]byte, error) {
	v := (*uint256.Int)(q)
	if v == nil {
		return []byte("null"), nil
	}

	return fmt.Appendf(nil, `"0x%x"`, v), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (q *QuantityU256) UnmarshalJSON(input []byte) error {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return errors.New("quantity256: not a JSON string")
	}

	body := input[1 : len(input)-1]
	if len(body) < 3 || body[0] != '0' || (body[1] != 'x' && body[1] != 'X') {
		return errors.New("quantity256: missing 0x prefix")
	}

	v, err := uint256.FromHex(string(body))
	if err != nil {
		return errors.Wrap(err, "quantity256")
	}

	*q = QuantityU256(*v)

	return nil
}

// Bytes is a variable-length byte string marshaled as `"0x..."` (with
// `"0x"` representing an empty slice).
type Bytes []byte

// MarshalJSON implements json.Marshaler.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte(`"0x"`), nil
	}

	return fmt.Appendf(nil, `"%#x"`, []byte(b)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return errors.New("bytes: not a JSON string")
	}

	body := input[1 : len(input)-1]
	if len(body) < 2 || body[0] != '0' || (body[1] != 'x' && body[1] != 'X') {
		return errors.New("bytes: missing 0x prefix")
	}

	hexBody := body[2:]
	if len(hexBody) == 0 {
		*b = nil

		return nil
	}

	decoded := make([]byte, hex.DecodedLen(len(hexBody)))

	n, err := hex.Decode(decoded, hexBody)
	if err != nil {
		return errors.Wrap(err, "bytes")
	}

	*b = decoded[:n]

	return nil
}

// DecodeFixed decodes a JSON-quoted, 0x-prefixed hex string of exactly
// length bytes into dst. `name` is used in error messages.
func DecodeFixed(dst, input []byte, length int, name string) error {
	if len(input) == 0 {
		return fmt.Errorf("%s: input missing", name)
	}

	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return fmt.Errorf("%s: not a JSON string", name)
	}

	return DecodeFixedText(dst, input[1:len(input)-1], length, name)
}

// DecodeFixedText decodes a raw 0x-prefixed hex string of exactly length
// bytes into dst.
func DecodeFixedText(dst, input []byte, length int, name string) error {
	if len(input) < 2 || input[0] != '0' || (input[1] != 'x' && input[1] != 'X') {
		return fmt.Errorf("%s: missing 0x prefix", name)
	}

	hexBody := input[2:]
	if len(hexBody) != length*2 {
		return fmt.Errorf("%s: expected %d hex chars, got %d", name, length*2, len(hexBody))
	}

	n, err := hex.Decode(dst, hexBody)
	if err != nil {
		return errors.Wrapf(err, "%s: invalid hex", name)
	}

	if n != length {
		return fmt.Errorf("%s: decoded %d bytes, expected %d", name, n, length)
	}

	return nil
}

// EncodeFixed returns the JSON-quoted, 0x-prefixed hex string of src.
func EncodeFixed(src []byte) []byte {
	return fmt.Appendf(nil, `"%#x"`, src)
}

// EncodeFixedText returns the unquoted 0x-prefixed hex string of src.
func EncodeFixedText(src []byte) []byte {
	out := make([]byte, 2+len(src)*2)
	copy(out, "0x")
	hex.Encode(out[2:], src)

	return out
}
