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

package jsonrpc_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethpandaops/go-eth-engine-client/spec/identification"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

func TestGetPayloadBodiesByHashV1(t *testing.T) {
	secret := testSecret()

	var (
		gotMethod  string
		gotnHashes int
	)

	srv := newTestServer(t, secret, func(method string, params []json.RawMessage) any {
		gotMethod = method

		var hashes []string
		_ = json.Unmarshal(params[0], &hashes)
		gotnHashes = len(hashes)

		// A present body followed by a null (unknown block).
		return []any{
			map[string]any{"transactions": []any{}, "withdrawals": []any{}},
			nil,
		}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	bodies, err := svc.GetPayloadBodiesByHash(
		context.Background(),
		version.DataVersionShanghai,
		[]paris.Hash32{{0x01}, {0x02}},
	)
	if err != nil {
		t.Fatalf("GetPayloadBodiesByHash: %v", err)
	}

	if gotMethod != "engine_getPayloadBodiesByHashV1" {
		t.Fatalf("expected V1 method, got %s", gotMethod)
	}

	if gotnHashes != 2 {
		t.Fatalf("expected 2 hashes sent, got %d", gotnHashes)
	}

	if len(bodies) != 2 {
		t.Fatalf("expected 2 bodies, got %d", len(bodies))
	}

	if bodies[0] == nil || bodies[0].Version != version.DataVersionShanghai || bodies[0].Shanghai == nil {
		t.Fatalf("first body not parsed: %+v", bodies[0])
	}

	if bodies[1] != nil {
		t.Fatalf("expected nil for unknown block")
	}

	// Agnostic variant.
	ag, err := svc.GetPayloadBodiesByHashAgnostic(
		context.Background(),
		version.DataVersionShanghai,
		[]paris.Hash32{{0x01}, {0x02}},
	)
	if err != nil {
		t.Fatalf("GetPayloadBodiesByHashAgnostic: %v", err)
	}

	if len(ag) != 2 || ag[0] == nil || ag[0].Version != version.DataVersionShanghai {
		t.Fatalf("agnostic bodies mismatch")
	}

	if ag[1] != nil {
		t.Fatalf("expected nil agnostic body for unknown block")
	}
}

func TestGetPayloadBodiesByRangeV1(t *testing.T) {
	secret := testSecret()

	var gotParams []string

	srv := newTestServer(t, secret, func(_ string, params []json.RawMessage) any {
		for _, p := range params {
			var s string
			_ = json.Unmarshal(p, &s)
			gotParams = append(gotParams, s)
		}

		return []any{map[string]any{"transactions": []any{}, "withdrawals": []any{}}}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	bodies, err := svc.GetPayloadBodiesByRange(context.Background(), version.DataVersionShanghai, 16, 1)
	if err != nil {
		t.Fatalf("GetPayloadBodiesByRange: %v", err)
	}

	// start/count must be hex QUANTITY.
	if len(gotParams) != 2 || gotParams[0] != "0x10" || gotParams[1] != "0x1" {
		t.Fatalf("expected hex range params [0x10 0x1], got %v", gotParams)
	}

	if len(bodies) != 1 || bodies[0] == nil {
		t.Fatalf("expected 1 body")
	}
}

func TestClientVersionV1(t *testing.T) {
	secret := testSecret()

	var gotMethod string

	srv := newTestServer(t, secret, func(method string, params []json.RawMessage) any {
		gotMethod = method

		// Echo a single EL client version back.
		var sent map[string]any
		_ = json.Unmarshal(params[0], &sent)

		return []any{
			map[string]any{
				"code":    "GE",
				"name":    "geth",
				"version": "1.14.0",
				"commit":  "0x12345678",
			},
		}
	})
	defer srv.Close()

	svc := newTestService(t, srv, secret)

	cl := &identification.ClientVersion{
		Code:    []byte("LH"),
		Name:    []byte("lighthouse"),
		Version: []byte("5.0.0"),
		Commit:  [4]byte{0xde, 0xad, 0xbe, 0xef},
	}

	out, err := svc.ClientVersion(context.Background(), cl)
	if err != nil {
		t.Fatalf("ClientVersion: %v", err)
	}

	if gotMethod != "engine_getClientVersionV1" {
		t.Fatalf("expected engine_getClientVersionV1, got %s", gotMethod)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 client version, got %d", len(out))
	}

	if string(out[0].Code) != "GE" || string(out[0].Name) != "geth" {
		t.Fatalf("client version mismatch: %s", out[0])
	}

	if out[0].Commit != [4]byte{0x12, 0x34, 0x56, 0x78} {
		t.Fatalf("commit mismatch: %x", out[0].Commit)
	}
}
