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

package paris_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/holiman/uint256"
)

func mustHash(s string) paris.Hash32 {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		panic(err)
	}

	var h paris.Hash32

	copy(h[:], b)

	return h
}

func mustAddress(s string) paris.Address {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		panic(err)
	}

	var a paris.Address

	copy(a[:], b)

	return a
}

func samplePayload() *paris.ExecutionPayload {
	return &paris.ExecutionPayload{
		ParentHash:    mustHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
		FeeRecipient:  mustAddress("0x2222222222222222222222222222222222222222"),
		StateRoot:     mustHash("0x3333333333333333333333333333333333333333333333333333333333333333"),
		ReceiptsRoot:  mustHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
		PrevRandao:    mustHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
		BlockNumber:   42,
		GasLimit:      30_000_000,
		GasUsed:       21_000,
		Timestamp:     1_700_000_000,
		ExtraData:     []byte("hello"),
		BaseFeePerGas: uint256.NewInt(1_000_000_000),
		BlockHash:     mustHash("0x6666666666666666666666666666666666666666666666666666666666666666"),
		Transactions: []paris.Transaction{
			[]byte{0x02, 0xff, 0xee},
			[]byte{0x01, 0xab, 0xcd, 0xef},
		},
	}
}

func TestExecutionPayloadJSONRoundtrip(t *testing.T) {
	orig := samplePayload()

	encoded, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if !bytes.Contains(encoded, []byte(`"blockNumber":"0x2a"`)) {
		t.Fatalf("expected hex-encoded blockNumber 0x2a in %s", encoded)
	}

	if !bytes.Contains(encoded, []byte(`"baseFeePerGas":"0x3b9aca00"`)) {
		t.Fatalf("expected hex-encoded baseFeePerGas in %s", encoded)
	}

	var decoded paris.ExecutionPayload
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, &decoded) {
		t.Fatalf("JSON roundtrip mismatch\n orig:    %+v\n decoded: %+v", orig, &decoded)
	}
}

func TestExecutionPayloadSSZRoundtrip(t *testing.T) {
	orig := samplePayload()

	encoded, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded paris.ExecutionPayload
	if err := decoded.UnmarshalSSZ(encoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, &decoded) {
		t.Fatalf("SSZ roundtrip mismatch\n orig:    %+v\n decoded: %+v", orig, &decoded)
	}
}

func TestPayloadStatusRoundtrip(t *testing.T) {
	hash := mustHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	orig := &paris.PayloadStatus{
		Status:          paris.PayloadValidationStatusInvalid,
		LatestValidHash: &hash,
		ValidationError: []byte("block tx invalid"),
	}

	encoded, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}

	if !bytes.Contains(encoded, []byte(`"status":"INVALID"`)) {
		t.Fatalf("expected INVALID status in %s", encoded)
	}

	if !bytes.Contains(encoded, []byte(`"validationError":"block tx invalid"`)) {
		t.Fatalf("expected validation error string in %s", encoded)
	}

	var decoded paris.PayloadStatus
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, &decoded) {
		t.Fatalf("JSON roundtrip mismatch\n orig:    %+v\n decoded: %+v", orig, &decoded)
	}

	sszBytes, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("ssz marshal: %v", err)
	}

	var sszDecoded paris.PayloadStatus
	if err := sszDecoded.UnmarshalSSZ(sszBytes); err != nil {
		t.Fatalf("ssz unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, &sszDecoded) {
		t.Fatalf("SSZ roundtrip mismatch\n orig:    %+v\n decoded: %+v", orig, &sszDecoded)
	}
}

func TestPayloadStatusNullableSSZ(t *testing.T) {
	orig := &paris.PayloadStatus{
		Status: paris.PayloadValidationStatusSyncing,
	}

	encoded, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("ssz marshal: %v", err)
	}

	var decoded paris.PayloadStatus
	if err := decoded.UnmarshalSSZ(encoded); err != nil {
		t.Fatalf("ssz unmarshal: %v", err)
	}

	if decoded.LatestValidHash != nil {
		t.Fatalf("expected nil LatestValidHash, got %v", *decoded.LatestValidHash)
	}

	if decoded.ValidationError != nil {
		t.Fatalf("expected nil ValidationError, got %q", decoded.ValidationError)
	}

	if decoded.Status != paris.PayloadValidationStatusSyncing {
		t.Fatalf("expected SYNCING status, got %s", decoded.Status)
	}
}

func TestForkchoiceUpdatedResponseNullablePayloadID(t *testing.T) {
	orig := &paris.ForkchoiceUpdatedResponse{
		PayloadStatus: paris.PayloadStatus{
			Status: paris.PayloadValidationStatusValid,
		},
	}

	encoded, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("ssz marshal: %v", err)
	}

	var decoded paris.ForkchoiceUpdatedResponse
	if err := decoded.UnmarshalSSZ(encoded); err != nil {
		t.Fatalf("ssz unmarshal: %v", err)
	}

	if decoded.PayloadID != nil {
		t.Fatalf("expected nil PayloadID, got %s", decoded.PayloadID)
	}

	id := paris.PayloadID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	orig.PayloadID = &id

	encoded2, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("ssz marshal: %v", err)
	}

	var decoded2 paris.ForkchoiceUpdatedResponse
	if err := decoded2.UnmarshalSSZ(encoded2); err != nil {
		t.Fatalf("ssz unmarshal: %v", err)
	}

	if decoded2.PayloadID == nil || *decoded2.PayloadID != id {
		t.Fatalf("PayloadID roundtrip mismatch: got %v", decoded2.PayloadID)
	}
}
