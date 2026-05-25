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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
)

// PayloadAttributes is the SSZ container `PayloadAttributesV4`: V3 extended
// with `slotNumber` and `targetGasLimit`.
type PayloadAttributes struct {
	Timestamp             uint64                 `json:"timestamp"`
	PrevRandao            paris.Hash32           `ssz-size:"32"                            json:"prevRandao"`
	SuggestedFeeRecipient paris.Address          `ssz-size:"20"                            json:"suggestedFeeRecipient"`
	Withdrawals           []*shanghai.Withdrawal `dynssz-max:"MAX_WITHDRAWALS_PER_PAYLOAD" ssz-max:"16" json:"withdrawals"`
	ParentBeaconBlockRoot paris.Hash32           `ssz-size:"32"                            json:"parentBeaconBlockRoot"`
	SlotNumber            uint64                 `json:"slotNumber"`
	TargetGasLimit        uint64                 `json:"targetGasLimit"`
}

type payloadAttributesJSON struct {
	Timestamp             jsonhex.QuantityU64    `json:"timestamp"`
	PrevRandao            paris.Hash32           `json:"prevRandao"`
	SuggestedFeeRecipient paris.Address          `json:"suggestedFeeRecipient"`
	Withdrawals           []*shanghai.Withdrawal `json:"withdrawals"`
	ParentBeaconBlockRoot paris.Hash32           `json:"parentBeaconBlockRoot"`
	SlotNumber            jsonhex.QuantityU64    `json:"slotNumber"`
	TargetGasLimit        jsonhex.QuantityU64    `json:"targetGasLimit"`
}

// MarshalJSON implements json.Marshaler.
func (p *PayloadAttributes) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&payloadAttributesJSON{
		Timestamp:             jsonhex.QuantityU64(p.Timestamp),
		PrevRandao:            p.PrevRandao,
		SuggestedFeeRecipient: p.SuggestedFeeRecipient,
		Withdrawals:           p.Withdrawals,
		ParentBeaconBlockRoot: p.ParentBeaconBlockRoot,
		SlotNumber:            jsonhex.QuantityU64(p.SlotNumber),
		TargetGasLimit:        jsonhex.QuantityU64(p.TargetGasLimit),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributes) UnmarshalJSON(input []byte) error {
	var data payloadAttributesJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "PayloadAttributes")
	}

	if data.Withdrawals == nil {
		return errors.New("PayloadAttributes: withdrawals missing")
	}

	p.Timestamp = uint64(data.Timestamp)
	p.PrevRandao = data.PrevRandao
	p.SuggestedFeeRecipient = data.SuggestedFeeRecipient
	p.Withdrawals = data.Withdrawals
	p.ParentBeaconBlockRoot = data.ParentBeaconBlockRoot
	p.SlotNumber = uint64(data.SlotNumber)
	p.TargetGasLimit = uint64(data.TargetGasLimit)

	return nil
}

// String returns a JSON representation.
func (p *PayloadAttributes) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
