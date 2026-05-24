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

// ForkchoiceUpdatedResponse is the response from any version of
// engine_forkchoiceUpdated. The SSZ transport defines a single container
// `ForkchoiceUpdatedResponseV1` shared by all versions.
//
// PayloadID is nullable: it is set when a payload build process was
// initiated and unset (nil) when no build is in progress (for instance when
// payloadAttributes was nil or the head was not VALID). The dynssz
// `optional-list` annotation encodes the pointer as the canonical SSZ
// `List[Bytes8, 1]`.
type ForkchoiceUpdatedResponse struct {
	PayloadStatus PayloadStatus
	PayloadID     *PayloadID `ssz-type:"optional-list" ssz-size:"8"`
}

type forkchoiceUpdatedResponseJSON struct {
	PayloadStatus PayloadStatus `json:"payloadStatus"`
	PayloadID     *PayloadID    `json:"payloadId"`
}

// MarshalJSON implements json.Marshaler.
func (r *ForkchoiceUpdatedResponse) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&forkchoiceUpdatedResponseJSON{
		PayloadStatus: r.PayloadStatus,
		PayloadID:     r.PayloadID,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *ForkchoiceUpdatedResponse) UnmarshalJSON(input []byte) error {
	var data forkchoiceUpdatedResponseJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "ForkchoiceUpdatedResponse")
	}

	r.PayloadStatus = data.PayloadStatus
	r.PayloadID = data.PayloadID

	return nil
}

// String returns a JSON representation of the response.
func (r *ForkchoiceUpdatedResponse) String() string {
	out, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
