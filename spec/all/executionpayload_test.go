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

package all_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/holiman/uint256"

	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

func sampleUnion(v version.DataVersion) *all.ExecutionPayload {
	e := &all.ExecutionPayload{
		Version:       v,
		ParentHash:    paris.Hash32{0x11},
		FeeRecipient:  paris.Address{0x22},
		StateRoot:     paris.Hash32{0x33},
		ReceiptsRoot:  paris.Hash32{0x44},
		PrevRandao:    paris.Hash32{0x55},
		BlockNumber:   42,
		GasLimit:      30_000_000,
		GasUsed:       21_000,
		Timestamp:     1_700_000_000,
		ExtraData:     []byte("ethpandaops"),
		BaseFeePerGas: uint256.NewInt(1_000_000_000),
		BlockHash:     paris.Hash32{0x66},
		Transactions:  []paris.Transaction{{0x02, 0xaa}, {0x01, 0xbb}},
	}

	if v >= version.DataVersionShanghai {
		e.Withdrawals = []*shanghai.Withdrawal{
			{Index: 1, ValidatorIndex: 7, Address: paris.Address{0x77}, Amount: 32_000_000_000},
		}
	}

	if v >= version.DataVersionCancun {
		e.BlobGasUsed = 131072
		e.ExcessBlobGas = 262144
	}

	if v >= version.DataVersionAmsterdam {
		e.BlockAccessList = amsterdam.BlockAccessList{0xde, 0xad, 0xbe, 0xef}
		e.SlotNumber = 9000
	}

	return e
}

func TestExecutionPayloadToFromView(t *testing.T) {
	for _, v := range []version.DataVersion{
		version.DataVersionParis,
		version.DataVersionShanghai,
		version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka,
		version.DataVersionAmsterdam,
	} {
		t.Run(v.String(), func(t *testing.T) {
			orig := sampleUnion(v)

			view, err := orig.ToView()
			if err != nil {
				t.Fatalf("ToView: %v", err)
			}

			switch v {
			case version.DataVersionParis:
				if _, ok := view.(*paris.ExecutionPayload); !ok {
					t.Fatalf("expected *paris.ExecutionPayload, got %T", view)
				}
			case version.DataVersionShanghai:
				if _, ok := view.(*shanghai.ExecutionPayload); !ok {
					t.Fatalf("expected *shanghai.ExecutionPayload, got %T", view)
				}
			case version.DataVersionCancun, version.DataVersionPrague, version.DataVersionOsaka:
				if _, ok := view.(*cancun.ExecutionPayload); !ok {
					t.Fatalf("expected *cancun.ExecutionPayload, got %T", view)
				}
			case version.DataVersionAmsterdam:
				if _, ok := view.(*amsterdam.ExecutionPayload); !ok {
					t.Fatalf("expected *amsterdam.ExecutionPayload, got %T", view)
				}
			}

			// Round-trip back into a union with Version pinned.
			back := &all.ExecutionPayload{Version: v}
			if err := back.FromView(view); err != nil {
				t.Fatalf("FromView: %v", err)
			}

			if back.BlockNumber != orig.BlockNumber || back.BaseFeePerGas.Cmp(orig.BaseFeePerGas) != 0 {
				t.Fatalf("view roundtrip mismatch: %+v vs %+v", back, orig)
			}

			if v >= version.DataVersionAmsterdam && back.SlotNumber != orig.SlotNumber {
				t.Fatalf("slot number lost: got %d want %d", back.SlotNumber, orig.SlotNumber)
			}
		})
	}
}

func TestExecutionPayloadToFromVersioned(t *testing.T) {
	orig := sampleUnion(version.DataVersionPrague)

	ver, err := orig.ToVersioned()
	if err != nil {
		t.Fatalf("ToVersioned: %v", err)
	}

	if ver.Version != version.DataVersionPrague {
		t.Fatalf("expected prague version, got %s", ver.Version)
	}

	if ver.Prague == nil {
		t.Fatalf("expected Prague field populated")
	}

	back := &all.ExecutionPayload{}
	if err := back.FromVersioned(ver); err != nil {
		t.Fatalf("FromVersioned: %v", err)
	}

	if back.Version != version.DataVersionPrague {
		t.Fatalf("expected prague version after FromVersioned, got %s", back.Version)
	}

	if back.BlockNumber != orig.BlockNumber || back.BlobGasUsed != orig.BlobGasUsed {
		t.Fatalf("versioned roundtrip mismatch")
	}
}

func TestExecutionPayloadJSONMatchesView(t *testing.T) {
	orig := sampleUnion(version.DataVersionCancun)

	unionJSON, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("union marshal: %v", err)
	}

	view, _ := orig.ToView()

	viewJSON, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("view marshal: %v", err)
	}

	if !bytes.Equal(unionJSON, viewJSON) {
		t.Fatalf("union JSON != view JSON\n union: %s\n view:  %s", unionJSON, viewJSON)
	}

	decoded := &all.ExecutionPayload{Version: version.DataVersionCancun}
	if err := json.Unmarshal(unionJSON, decoded); err != nil {
		t.Fatalf("union unmarshal: %v", err)
	}

	if decoded.BlobGasUsed != orig.BlobGasUsed {
		t.Fatalf("json roundtrip lost blobGasUsed")
	}
}

func TestExecutionPayloadSSZMatchesView(t *testing.T) {
	orig := sampleUnion(version.DataVersionAmsterdam)

	unionSSZ, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("union ssz marshal: %v", err)
	}

	view, _ := orig.ToView()

	viewSSZ, err := view.(*amsterdam.ExecutionPayload).MarshalSSZ()
	if err != nil {
		t.Fatalf("view ssz marshal: %v", err)
	}

	if !bytes.Equal(unionSSZ, viewSSZ) {
		t.Fatalf("union SSZ != view SSZ")
	}

	decoded := &all.ExecutionPayload{Version: version.DataVersionAmsterdam}
	if err := decoded.UnmarshalSSZ(unionSSZ); err != nil {
		t.Fatalf("union ssz unmarshal: %v", err)
	}

	if decoded.SlotNumber != orig.SlotNumber || decoded.BlockNumber != orig.BlockNumber {
		t.Fatalf("ssz roundtrip mismatch")
	}

	origHTR, err := orig.HashTreeRoot()
	if err != nil {
		t.Fatalf("htr: %v", err)
	}

	viewHTR, err := view.(*amsterdam.ExecutionPayload).HashTreeRoot()
	if err != nil {
		t.Fatalf("view htr: %v", err)
	}

	if origHTR != viewHTR {
		t.Fatalf("union HTR != view HTR")
	}
}
