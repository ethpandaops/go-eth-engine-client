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

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// Withdrawal is a single beacon-chain validator withdrawal pushed onto the
// execution layer. It corresponds to the SSZ container `WithdrawalV1` and
// the JSON `WithdrawalV1` schema.
type Withdrawal struct {
	Index          uint64        `json:"index"`
	ValidatorIndex uint64        `json:"validatorIndex"`
	Address        paris.Address `ssz-size:"20" json:"address"`
	Amount         uint64        `json:"amount"`
}

type withdrawalJSON struct {
	Index          jsonhex.QuantityU64 `json:"index"`
	ValidatorIndex jsonhex.QuantityU64 `json:"validatorIndex"`
	Address        paris.Address       `json:"address"`
	Amount         jsonhex.QuantityU64 `json:"amount"`
}

// MarshalJSON implements json.Marshaler.
func (w *Withdrawal) MarshalJSON() ([]byte, error) {
	if w == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&withdrawalJSON{
		Index:          jsonhex.QuantityU64(w.Index),
		ValidatorIndex: jsonhex.QuantityU64(w.ValidatorIndex),
		Address:        w.Address,
		Amount:         jsonhex.QuantityU64(w.Amount),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Withdrawal) UnmarshalJSON(input []byte) error {
	var data withdrawalJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "Withdrawal")
	}

	w.Index = uint64(data.Index)
	w.ValidatorIndex = uint64(data.ValidatorIndex)
	w.Address = data.Address
	w.Amount = uint64(data.Amount)

	return nil
}

// String returns a JSON representation.
func (w *Withdrawal) String() string {
	out, err := json.Marshal(w)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
