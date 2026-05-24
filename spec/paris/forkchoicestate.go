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
)

// ForkchoiceState is the first parameter to engine_forkchoiceUpdated*. It
// corresponds to the SSZ container `ForkchoiceStateV1`.
type ForkchoiceState struct {
	HeadBlockHash      Hash32 `ssz-size:"32" json:"headBlockHash"`
	SafeBlockHash      Hash32 `ssz-size:"32" json:"safeBlockHash"`
	FinalizedBlockHash Hash32 `ssz-size:"32" json:"finalizedBlockHash"`
}

// MarshalJSON implements json.Marshaler.
func (f *ForkchoiceState) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}

	type alias ForkchoiceState

	return json.Marshal((*alias)(f))
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkchoiceState) UnmarshalJSON(input []byte) error {
	type alias ForkchoiceState

	if err := json.Unmarshal(input, (*alias)(f)); err != nil {
		return errors.Wrap(err, "ForkchoiceState")
	}

	return nil
}

// String returns a JSON representation of the state.
func (f *ForkchoiceState) String() string {
	out, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
