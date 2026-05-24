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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
)

// PayloadAttributes carries the inputs to a payload build process. It
// corresponds to the SSZ container `PayloadAttributesV1` and is the
// `payloadAttributes` parameter of engine_forkchoiceUpdatedV1.
type PayloadAttributes struct {
	Timestamp             uint64  `json:"timestamp"`
	PrevRandao            Hash32  `ssz-size:"32" json:"prevRandao"`
	SuggestedFeeRecipient Address `ssz-size:"20" json:"suggestedFeeRecipient"`
}

type payloadAttributesJSON struct {
	Timestamp             jsonhex.QuantityU64 `json:"timestamp"`
	PrevRandao            Hash32              `json:"prevRandao"`
	SuggestedFeeRecipient Address             `json:"suggestedFeeRecipient"`
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
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributes) UnmarshalJSON(input []byte) error {
	var data payloadAttributesJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "PayloadAttributes")
	}

	p.Timestamp = uint64(data.Timestamp)
	p.PrevRandao = data.PrevRandao
	p.SuggestedFeeRecipient = data.SuggestedFeeRecipient

	return nil
}

// String returns a JSON representation of the attributes.
func (p *PayloadAttributes) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
