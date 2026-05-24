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

package shanghai

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// ExecutionPayloadBody is the per-block body returned by
// engine_getPayloadBodiesByHashV1 / engine_getPayloadBodiesByRangeV1. It
// corresponds to the SSZ container `ExecutionPayloadBodyV1`.
//
// Withdrawals is `null` (Go nil slice) for pre-Shanghai blocks in JSON. In
// the SSZ encoding withdrawals is always present (an empty list represents
// the pre-Shanghai case), so a nil withdrawals slice round-trips through
// SSZ as an empty (non-nil, len-0) slice.
type ExecutionPayloadBody struct {
	Transactions []paris.Transaction `dynssz-max:"MAX_TRANSACTIONS_PER_PAYLOAD,MAX_BYTES_PER_TRANSACTION" ssz-max:"1048576,1073741824" ssz-size:"?,?" json:"transactions"`
	Withdrawals  []*Withdrawal       `dynssz-max:"MAX_WITHDRAWALS_PER_PAYLOAD"                            ssz-max:"16"                                 json:"withdrawals"`
}

// MarshalJSON implements json.Marshaler.
func (b *ExecutionPayloadBody) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	type alias ExecutionPayloadBody

	return json.Marshal((*alias)(b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *ExecutionPayloadBody) UnmarshalJSON(input []byte) error {
	type alias ExecutionPayloadBody

	if err := json.Unmarshal(input, (*alias)(b)); err != nil {
		return errors.Wrap(err, "ExecutionPayloadBody")
	}

	return nil
}

// String returns a JSON representation.
func (b *ExecutionPayloadBody) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
