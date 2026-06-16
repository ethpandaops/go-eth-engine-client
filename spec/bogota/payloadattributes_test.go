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

package bogota_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethpandaops/go-eth-engine-client/spec/bogota"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
)

// TestPayloadAttributesSSZRoundtrip checks that PayloadAttributesV5 (V4 +
// inclusionListTransactions) survives an SSZ encode/decode roundtrip.
func TestPayloadAttributesSSZRoundtrip(t *testing.T) {
	orig := &bogota.PayloadAttributes{
		Timestamp:             1700000000,
		PrevRandao:            paris.Hash32{0xaa},
		SuggestedFeeRecipient: paris.Address{0xbb},
		Withdrawals: []*shanghai.Withdrawal{
			{Index: 1, ValidatorIndex: 2, Address: paris.Address{0xcc}, Amount: 1_000_000},
		},
		ParentBeaconBlockRoot: paris.Hash32{0xdd},
		SlotNumber:            42,
		TargetGasLimit:        30_000_000,
		InclusionListTransactions: []paris.Transaction{
			{0x02, 0xf8, 0x01}, // a couple of opaque EIP-2718 prefix bytes
			{0xee, 0xff},
		},
	}

	encoded, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("MarshalSSZ: %v", err)
	}

	var decoded bogota.PayloadAttributes
	if err := decoded.UnmarshalSSZ(encoded); err != nil {
		t.Fatalf("UnmarshalSSZ: %v", err)
	}

	if len(decoded.InclusionListTransactions) != 2 {
		t.Fatalf("inclusion list length mismatch: %d", len(decoded.InclusionListTransactions))
	}

	if !bytes.Equal(decoded.InclusionListTransactions[0], orig.InclusionListTransactions[0]) {
		t.Fatalf("inclusion list tx 0 mismatch: %x", decoded.InclusionListTransactions[0])
	}

	if decoded.SlotNumber != orig.SlotNumber || decoded.TargetGasLimit != orig.TargetGasLimit {
		t.Fatalf("V4 fields not preserved: slot=%d target=%d", decoded.SlotNumber, decoded.TargetGasLimit)
	}
}

// TestPayloadAttributesJSONRoundtrip checks the JSON wire format includes
// inclusionListTransactions as a hex-encoded byte array.
func TestPayloadAttributesJSONRoundtrip(t *testing.T) {
	orig := &bogota.PayloadAttributes{
		Timestamp:                 1,
		Withdrawals:               []*shanghai.Withdrawal{},
		ParentBeaconBlockRoot:     paris.Hash32{0x01},
		SlotNumber:                7,
		TargetGasLimit:            8,
		InclusionListTransactions: []paris.Transaction{{0xde, 0xad, 0xbe, 0xef}},
	}

	raw, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	if !bytes.Contains(raw, []byte(`"inclusionListTransactions":["0xdeadbeef"]`)) {
		t.Fatalf("unexpected JSON encoding: %s", raw)
	}

	var decoded bogota.PayloadAttributes
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}

	if len(decoded.InclusionListTransactions) != 1 ||
		!bytes.Equal(decoded.InclusionListTransactions[0], []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("JSON roundtrip mismatch: %#v", decoded.InclusionListTransactions)
	}
}
