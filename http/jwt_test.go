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

package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestParseJWTSecret(t *testing.T) {
	want := make([]byte, 32)
	for i := range want {
		want[i] = byte(i)
	}

	hexStr := "0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

	for _, in := range []string{hexStr, strings.TrimPrefix(hexStr, "0x"), "  " + hexStr + "\n"} {
		got, err := parseJWTSecret(in)
		if err != nil {
			t.Fatalf("parseJWTSecret(%q): %v", in, err)
		}

		if string(got) != string(want) {
			t.Fatalf("parseJWTSecret(%q) mismatch", in)
		}
	}

	if _, err := parseJWTSecret("0x1234"); err == nil {
		t.Fatal("expected error for short secret")
	}
}

func TestJWTTokenStructure(t *testing.T) {
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(i + 1)
	}

	signer, err := newJWTSigner(secret)
	if err != nil {
		t.Fatalf("newJWTSigner: %v", err)
	}

	now := time.Unix(1_700_000_000, 0)

	tok, err := signer.token(now)
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT segments, got %d", len(parts))
	}

	// Header must decode to HS256/JWT.
	hdr, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}

	if string(hdr) != `{"alg":"HS256","typ":"JWT"}` {
		t.Fatalf("unexpected header: %s", hdr)
	}

	// Claims must carry the iat we provided.
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode claims: %v", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}

	if iat, ok := claims["iat"].(float64); !ok || int64(iat) != now.Unix() {
		t.Fatalf("iat claim mismatch: %v", claims["iat"])
	}

	// Signature must verify against the secret.
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	wantSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != wantSig {
		t.Fatalf("signature mismatch")
	}
}

func TestJWTTokenFreshIAT(t *testing.T) {
	signer, _ := newJWTSigner(make([]byte, 32))

	t1, _ := signer.token(time.Unix(1000, 0))
	t2, _ := signer.token(time.Unix(2000, 0))

	if t1 == t2 {
		t.Fatal("tokens with different iat should differ")
	}
}
