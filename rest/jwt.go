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

package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// jwtSecretLength is the required length of an Engine API JWT secret.
const jwtSecretLength = 32

// jwtSigner mints the short-lived HS256 bearer tokens required by the Engine
// API authentication scheme. The Marius spec keeps `iat` (±60s window) and
// the optional `id` claim but drops `clv` -- the CL client version moves to
// the X-Engine-Client-Version request header (see [WithClientVersionHeader]).
type jwtSigner struct {
	secret []byte
	header string
	id     string
}

// newJWTSigner validates the secret and returns a signer. The secret must be
// exactly 32 bytes.
func newJWTSigner(secret []byte) (*jwtSigner, error) {
	if len(secret) != jwtSecretLength {
		return nil, errors.Errorf("JWT secret must be %d bytes, got %d", jwtSecretLength, len(secret))
	}

	hdr := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	return &jwtSigner{
		secret: secret,
		header: hdr,
	}, nil
}

// loadJWTSecret reads a hex-encoded (optionally 0x-prefixed) 32-byte JWT
// secret from path.
func loadJWTSecret(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read JWT secret file")
	}

	return parseJWTSecret(string(raw))
}

// parseJWTSecret decodes a hex-encoded (optionally 0x-prefixed, whitespace
// tolerated) 32-byte JWT secret.
func parseJWTSecret(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")

	secret, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.Wrap(err, "decode JWT secret hex")
	}

	if len(secret) != jwtSecretLength {
		return nil, errors.Errorf("JWT secret must be %d bytes, got %d", jwtSecretLength, len(secret))
	}

	return secret, nil
}

// token mints a fresh HS256 JWT for the given issue time. Unlike the
// JSON-RPC signer, no `clv` claim is included; CL version is carried by the
// X-Engine-Client-Version request header instead.
func (s *jwtSigner) token(now time.Time) (string, error) {
	claims := map[string]any{"iat": now.Unix()}
	if s.id != "" {
		claims["id"] = s.id
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", errors.Wrap(err, "marshal JWT claims")
	}

	signingInput := s.header + "." + base64.RawURLEncoding.EncodeToString(payload)

	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, nil
}
