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
	"testing"

	"github.com/holiman/uint256"

	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

func TestNewPayloadRequestNestedSSZ(t *testing.T) {
	req := &all.NewPayloadRequest{
		Version:                     version.DataVersionPrague,
		ExecutionPayload:            sampleUnion(version.DataVersionCancun), // payload schema is V3 for prague
		ExpectedBlobVersionedHashes: []paris.Hash32{{0xaa}, {0xbb}},
		ParentBeaconBlockRoot:       paris.Hash32{0xcc},
		ExecutionRequests:           []prague.ExecutionRequest{{0x01, 0x02}},
	}
	req.ExecutionPayload.Version = version.DataVersionPrague

	// ToView should produce a *prague.NewPayloadRequest with a populated
	// *cancun.ExecutionPayload.
	view, err := req.ToView()
	if err != nil {
		t.Fatalf("ToView: %v", err)
	}

	pv, ok := view.(*prague.NewPayloadRequest)
	if !ok {
		t.Fatalf("expected *prague.NewPayloadRequest, got %T", view)
	}

	if pv.ExecutionPayload == nil || pv.ExecutionPayload.BlockNumber != 42 {
		t.Fatalf("nested execution payload not copied into view")
	}

	if len(pv.ExecutionRequests) != 1 {
		t.Fatalf("execution requests not copied")
	}

	// SSZ roundtrip via the union vs the concrete view.
	unionSSZ, err := req.MarshalSSZ()
	if err != nil {
		t.Fatalf("union ssz: %v", err)
	}

	viewSSZ, err := pv.MarshalSSZ()
	if err != nil {
		t.Fatalf("view ssz: %v", err)
	}

	if !bytes.Equal(unionSSZ, viewSSZ) {
		t.Fatalf("union SSZ != view SSZ")
	}

	decoded := &all.NewPayloadRequest{Version: version.DataVersionPrague}
	if err := decoded.UnmarshalSSZ(unionSSZ); err != nil {
		t.Fatalf("union ssz decode: %v", err)
	}

	if decoded.ExecutionPayload == nil {
		t.Fatalf("decoded payload nil")
	}

	// Version must propagate into the nested union child.
	if decoded.ExecutionPayload.Version != version.DataVersionPrague {
		t.Fatalf("nested version not propagated: got %s", decoded.ExecutionPayload.Version)
	}

	if decoded.ExecutionPayload.BlockNumber != 42 {
		t.Fatalf("decoded nested block number mismatch: %d", decoded.ExecutionPayload.BlockNumber)
	}
}

func TestGetPayloadResponseVersionedRoundtrip(t *testing.T) {
	resp := &all.GetPayloadResponse{
		Version:          version.DataVersionCancun,
		ExecutionPayload: sampleUnion(version.DataVersionCancun),
		BlockValue:       uint256.NewInt(12345),
		BlobsBundle:      &all.BlobsBundle{Version: version.DataVersionCancun},
	}
	resp.ExecutionPayload.Version = version.DataVersionCancun

	ver, err := resp.ToVersioned()
	if err != nil {
		t.Fatalf("ToVersioned: %v", err)
	}

	if ver.Cancun == nil {
		t.Fatalf("expected cancun response populated")
	}

	if ver.Cancun.BlockValue.Cmp(uint256.NewInt(12345)) != 0 {
		t.Fatalf("block value mismatch")
	}

	back := &all.GetPayloadResponse{}
	if err := back.FromVersioned(ver); err != nil {
		t.Fatalf("FromVersioned: %v", err)
	}

	if back.Version != version.DataVersionCancun {
		t.Fatalf("version mismatch: %s", back.Version)
	}

	if back.ExecutionPayload == nil || back.ExecutionPayload.Version != version.DataVersionCancun {
		t.Fatalf("nested payload version not propagated")
	}

	if back.ExecutionPayload.BlobGasUsed != resp.ExecutionPayload.BlobGasUsed {
		t.Fatalf("nested payload field mismatch")
	}
}
