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

// PayloadValidationStatus enumerates the outcomes a payload can be put in by
// the Execution Layer.
type PayloadValidationStatus uint8

const (
	// PayloadValidationStatusValid is the SSZ enum value for "VALID".
	PayloadValidationStatusValid PayloadValidationStatus = 0
	// PayloadValidationStatusInvalid is the SSZ enum value for "INVALID".
	PayloadValidationStatusInvalid PayloadValidationStatus = 1
	// PayloadValidationStatusSyncing is the SSZ enum value for "SYNCING".
	PayloadValidationStatusSyncing PayloadValidationStatus = 2
	// PayloadValidationStatusAccepted is the SSZ enum value for "ACCEPTED".
	PayloadValidationStatusAccepted PayloadValidationStatus = 3
	// PayloadValidationStatusInclusionListUnsatisfied is the bogota-only
	// "INCLUSION_LIST_UNSATISFIED" status returned by engine_newPayloadV6 when
	// the payload is otherwise VALID but fails the EIP-7805 inclusion-list
	// constraints. Appended after ACCEPTED in the
	// `PayloadStatusInclusionListUnsatisfied` enum.
	PayloadValidationStatusInclusionListUnsatisfied PayloadValidationStatus = 4
	// PayloadValidationStatusInvalidBlockHash is the JSON-only legacy
	// "INVALID_BLOCK_HASH" status from Paris. It has no SSZ encoding
	// (ssz-encoding.md's status table omits it).
	PayloadValidationStatusInvalidBlockHash PayloadValidationStatus = 255
)

var payloadValidationStatusStrings = map[PayloadValidationStatus]string{
	PayloadValidationStatusValid:                    "VALID",
	PayloadValidationStatusInvalid:                  "INVALID",
	PayloadValidationStatusSyncing:                  "SYNCING",
	PayloadValidationStatusAccepted:                 "ACCEPTED",
	PayloadValidationStatusInclusionListUnsatisfied: "INCLUSION_LIST_UNSATISFIED",
	PayloadValidationStatusInvalidBlockHash:         "INVALID_BLOCK_HASH",
}

var payloadValidationStatusFromString = map[string]PayloadValidationStatus{
	"VALID":                      PayloadValidationStatusValid,
	"INVALID":                    PayloadValidationStatusInvalid,
	"SYNCING":                    PayloadValidationStatusSyncing,
	"ACCEPTED":                   PayloadValidationStatusAccepted,
	"INCLUSION_LIST_UNSATISFIED": PayloadValidationStatusInclusionListUnsatisfied,
	"INVALID_BLOCK_HASH":         PayloadValidationStatusInvalidBlockHash,
}

// String returns the wire-format string for the validation status.
func (s PayloadValidationStatus) String() string {
	if name, ok := payloadValidationStatusStrings[s]; ok {
		return name
	}

	return fmt.Sprintf("UNKNOWN(%d)", uint8(s))
}

// MarshalJSON implements json.Marshaler.
func (s PayloadValidationStatus) MarshalJSON() ([]byte, error) {
	name, ok := payloadValidationStatusStrings[s]
	if !ok {
		return nil, fmt.Errorf("PayloadValidationStatus: unknown value %d", uint8(s))
	}

	return fmt.Appendf(nil, `%q`, name), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *PayloadValidationStatus) UnmarshalJSON(input []byte) error {
	var name string
	if err := json.Unmarshal(input, &name); err != nil {
		return errors.Wrap(err, "PayloadValidationStatus")
	}

	v, ok := payloadValidationStatusFromString[name]
	if !ok {
		return fmt.Errorf("PayloadValidationStatus: unrecognised %q", name)
	}

	*s = v

	return nil
}

// PayloadStatus is the response from engine_newPayloadV1 and the
// `payloadStatus` field of engine_forkchoiceUpdatedV1's response. It
// corresponds to the SSZ container `PayloadStatusV1`.
//
//   - LatestValidHash uses the new dynssz `optional-list` annotation, which
//     encodes a `*Hash32` as the canonical SSZ `List[Bytes32, 1]` (nil →
//     empty list, non-nil → 1-element list).
//   - ValidationError is a plain `ByteList[MAX_ERROR_MESSAGE_LENGTH]`. The
//     JSON wire format treats it as a nullable UTF-8 string; a nil slice
//     marshals as JSON `null` and an empty (non-nil) slice marshals as the
//     empty string. SSZ has no nil/empty distinction — both encode as a
//     zero-length ByteList.
type PayloadStatus struct {
	Status          PayloadValidationStatus
	LatestValidHash *Hash32 `ssz-type:"optional-list" ssz-size:"32"`
	ValidationError []byte  `dynssz-max:"MAX_ERROR_MESSAGE_LENGTH" ssz-max:"1024"`
}

type payloadStatusJSON struct {
	Status          PayloadValidationStatus `json:"status"`
	LatestValidHash *Hash32                 `json:"latestValidHash"`
	ValidationError *string                 `json:"validationError"`
}

// MarshalJSON implements json.Marshaler.
func (p *PayloadStatus) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}

	out := &payloadStatusJSON{
		Status:          p.Status,
		LatestValidHash: p.LatestValidHash,
	}

	if p.ValidationError != nil {
		s := string(p.ValidationError)
		out.ValidationError = &s
	}

	return json.Marshal(out)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadStatus) UnmarshalJSON(input []byte) error {
	var data payloadStatusJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "PayloadStatus")
	}

	p.Status = data.Status
	p.LatestValidHash = data.LatestValidHash

	if data.ValidationError == nil {
		p.ValidationError = nil
	} else {
		p.ValidationError = []byte(*data.ValidationError)
	}

	return nil
}

// String returns a JSON representation of the status.
func (p *PayloadStatus) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
