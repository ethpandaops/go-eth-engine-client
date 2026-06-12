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

package http_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/holiman/uint256"

	enginehttp "github.com/ethpandaops/go-eth-engine-client/http"
	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

func testSecret() []byte {
	s := make([]byte, 32)
	for i := range s {
		s[i] = byte(i + 1)
	}

	return s
}

// verifyJWT checks the Authorization header carries a valid HS256 token
// signed with secret.
func verifyJWT(t *testing.T, authHeader string, secret []byte) {
	t.Helper()

	token, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		t.Fatalf("missing Bearer prefix: %q", authHeader)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("malformed JWT")
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	wantSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != wantSig {
		t.Fatalf("invalid JWT signature")
	}

	claims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode claims: %v", err)
	}

	var c map[string]any
	if err := json.Unmarshal(claims, &c); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}

	if _, ok := c["iat"]; !ok {
		t.Fatalf("missing iat claim")
	}
}

type rpcEnvelope struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
	ID     uint64            `json:"id"`
}

func newTestServer(t *testing.T, secret []byte, handler func(method string, params []json.RawMessage) any) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyJWT(t, r.Header.Get("Authorization"), secret)

		body, _ := io.ReadAll(r.Body)

		var req rpcEnvelope
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("server decode request: %v", err)
		}

		result := handler(req.Method, req.Params)

		resultJSON, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("server marshal result: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  json.RawMessage(resultJSON),
		})
	}))
}

func newTestService(t *testing.T, srv *httptest.Server, secret []byte) *enginehttp.Service {
	t.Helper()

	svc, err := enginehttp.New(context.Background(),
		enginehttp.WithAddress(srv.URL),
		enginehttp.WithJWTSecret(secret),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return svc
}

func cancunPayload() *cancun.ExecutionPayload {
	return &cancun.ExecutionPayload{
		ParentHash:    paris.Hash32{0x11},
		FeeRecipient:  paris.Address{0x22},
		StateRoot:     paris.Hash32{0x33},
		ReceiptsRoot:  paris.Hash32{0x44},
		PrevRandao:    paris.Hash32{0x55},
		BlockNumber:   42,
		GasLimit:      30_000_000,
		GasUsed:       21_000,
		Timestamp:     1_700_000_000,
		ExtraData:     []byte{0xab},
		BaseFeePerGas: uint256.NewInt(1_000_000_000),
		BlockHash:     paris.Hash32{0x66},
		Transactions:  []paris.Transaction{},
		Withdrawals:   []*shanghai.Withdrawal{},
		BlobGasUsed:   131072,
		ExcessBlobGas: 0,
	}
}

func TestNewPayloadV3(t *testing.T) {
	secret := testSecret()

	var (
		gotMethod     string
		gotParamCount int
	)

	srv := newTestServer(t, secret, func(method string, params []json.RawMessage) any {
		gotMethod = method
		gotParamCount = len(params)

		return map[string]any{
			"status":          "VALID",
			"latestValidHash": "0x" + strings.Repeat("ab", 32),
			"validationError": nil,
		}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	req := &spec.VersionedNewPayloadRequest{
		Version: version.DataVersionCancun,
		Cancun: &cancun.NewPayloadRequest{
			ExecutionPayload:            cancunPayload(),
			ExpectedBlobVersionedHashes: []paris.Hash32{{0x01}},
			ParentBeaconBlockRoot:       paris.Hash32{0x02},
		},
	}

	status, err := svc.NewPayload(context.Background(), req)
	if err != nil {
		t.Fatalf("NewPayload: %v", err)
	}

	if gotMethod != "engine_newPayloadV3" {
		t.Fatalf("expected engine_newPayloadV3, got %s", gotMethod)
	}

	if gotParamCount != 3 {
		t.Fatalf("expected 3 params for V3, got %d", gotParamCount)
	}

	if status.Status != paris.PayloadValidationStatusValid {
		t.Fatalf("expected VALID, got %s", status.Status)
	}

	if status.LatestValidHash == nil {
		t.Fatalf("expected latestValidHash set")
	}
}

func TestForkchoiceUpdatedV3(t *testing.T) {
	secret := testSecret()

	var gotMethod string

	srv := newTestServer(t, secret, func(method string, _ []json.RawMessage) any {
		gotMethod = method

		return map[string]any{
			"payloadStatus": map[string]any{
				"status":          "VALID",
				"latestValidHash": "0x" + strings.Repeat("cd", 32),
				"validationError": nil,
			},
			"payloadId": "0x0102030405060708",
		}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	req := &spec.VersionedForkchoiceUpdatedRequest{
		Version: version.DataVersionCancun,
		Cancun: &cancun.ForkchoiceUpdatedRequest{
			ForkchoiceState: &paris.ForkchoiceState{
				HeadBlockHash:      paris.Hash32{0x01},
				SafeBlockHash:      paris.Hash32{0x02},
				FinalizedBlockHash: paris.Hash32{0x03},
			},
		},
	}

	resp, err := svc.ForkchoiceUpdated(context.Background(), req)
	if err != nil {
		t.Fatalf("ForkchoiceUpdated: %v", err)
	}

	if gotMethod != "engine_forkchoiceUpdatedV3" {
		t.Fatalf("expected engine_forkchoiceUpdatedV3, got %s", gotMethod)
	}

	if resp.PayloadStatus.Status != paris.PayloadValidationStatusValid {
		t.Fatalf("expected VALID status")
	}

	want := paris.PayloadID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	if resp.PayloadID == nil || *resp.PayloadID != want {
		t.Fatalf("payloadId mismatch: %s", resp.PayloadID)
	}
}

func TestGetPayloadV3(t *testing.T) {
	secret := testSecret()

	srv := newTestServer(t, secret, func(method string, params []json.RawMessage) any {
		if method != "engine_getPayloadV3" {
			t.Fatalf("expected engine_getPayloadV3, got %s", method)
		}

		var id string
		_ = json.Unmarshal(params[0], &id)

		if id != "0x0102030405060708" {
			t.Fatalf("unexpected payload id param: %s", id)
		}

		payloadJSON, _ := json.Marshal(cancunPayload())

		return map[string]any{
			"executionPayload":      json.RawMessage(payloadJSON),
			"blockValue":            "0x1bc16d674ec80000",
			"blobsBundle":           map[string]any{"commitments": []any{}, "proofs": []any{}, "blobs": []any{}},
			"shouldOverrideBuilder": false,
		}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	resp, err := svc.GetPayload(
		context.Background(),
		version.DataVersionCancun,
		paris.PayloadID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	)
	if err != nil {
		t.Fatalf("GetPayload: %v", err)
	}

	if resp.Version != version.DataVersionCancun || resp.Cancun == nil {
		t.Fatalf("expected cancun response, got %+v", resp)
	}

	if resp.Cancun.ExecutionPayload == nil || resp.Cancun.ExecutionPayload.BlockNumber != 42 {
		t.Fatalf("execution payload not parsed")
	}

	if resp.Cancun.BlockValue == nil ||
		resp.Cancun.BlockValue.Cmp(uint256.NewInt(2_000_000_000_000_000_000)) != 0 {
		t.Fatalf("block value mismatch: %v", resp.Cancun.BlockValue)
	}

	// Agnostic variant should yield the same data through the union type.
	ag, err := svc.GetPayloadAgnostic(
		context.Background(),
		version.DataVersionCancun,
		paris.PayloadID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	)
	if err != nil {
		t.Fatalf("GetPayloadAgnostic: %v", err)
	}

	if ag.ExecutionPayload == nil || ag.ExecutionPayload.BlockNumber != 42 {
		t.Fatalf("agnostic execution payload not parsed")
	}
}

func TestGetBlobsV1(t *testing.T) {
	secret := testSecret()

	var gotMethod string

	srv := newTestServer(t, secret, func(method string, _ []json.RawMessage) any {
		gotMethod = method

		// engine_getBlobs returns a bare array of blob-and-proof objects.
		bp, _ := json.Marshal(&cancun.BlobAndProof{})

		return []json.RawMessage{json.RawMessage(bp), json.RawMessage("null")}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	req := &spec.VersionedGetBlobsRequest{
		Version: version.DataVersionCancun,
		Cancun: &cancun.GetBlobsRequest{
			BlobVersionedHashes: []paris.Hash32{{0x01}, {0x02}},
		},
	}

	resp, err := svc.GetBlobs(context.Background(), req)
	if err != nil {
		t.Fatalf("GetBlobs: %v", err)
	}

	if gotMethod != "engine_getBlobsV1" {
		t.Fatalf("expected engine_getBlobsV1, got %s", gotMethod)
	}

	if resp.Cancun == nil || len(resp.Cancun.BlobsAndProofs) != 2 {
		t.Fatalf("expected 2 blob entries, got %+v", resp.Cancun)
	}

	if resp.Cancun.BlobsAndProofs[0] == nil {
		t.Fatalf("expected first blob present")
	}

	if resp.Cancun.BlobsAndProofs[1] != nil {
		t.Fatalf("expected second blob to be null (missing)")
	}
}
