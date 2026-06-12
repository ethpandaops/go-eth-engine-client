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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// jsonRPCRequest is a JSON-RPC 2.0 request envelope.
type jsonRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

// jsonRPCResponse is a JSON-RPC 2.0 response envelope.
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCError   `json:"error"`
}

// jsonRPCError is a JSON-RPC 2.0 error object. It implements error.
type jsonRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *jsonRPCError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("engine API error %d: %s (%s)", e.Code, e.Message, string(e.Data))
	}

	return fmt.Sprintf("engine API error %d: %s", e.Code, e.Message)
}

// call performs a single authenticated JSON-RPC call. The result is
// unmarshaled into out (which may be nil to discard the result). A `null`
// result is treated as an error; use callOptional for methods where null is
// a valid response.
func (s *Service) call(ctx context.Context, method string, params []any, out any) error {
	raw, err := s.callRaw(ctx, method, params)
	if err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	if len(raw) == 0 || string(raw) == "null" {
		return errors.Errorf("%s: null result", method)
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return errors.Wrapf(err, "%s: decode result", method)
	}

	return nil
}

// callOptional behaves like call but tolerates a JSON-RPC `null` result,
// which several Engine API methods (e.g. engine_getBlobs while syncing) use
// to signal "no data". It reports whether a non-null result was decoded.
func (s *Service) callOptional(ctx context.Context, method string, params []any, out any) (bool, error) {
	raw, err := s.callRaw(ctx, method, params)
	if err != nil {
		return false, err
	}

	if len(raw) == 0 || string(raw) == "null" {
		return false, nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return false, errors.Wrapf(err, "%s: decode result", method)
	}

	return true, nil
}

// callRaw performs a single authenticated JSON-RPC call and returns the raw
// result bytes (which may be `null`).
func (s *Service) callRaw(ctx context.Context, method string, params []any) (json.RawMessage, error) {
	if params == nil {
		params = []any{}
	}

	reqBody, err := json.Marshal(&jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      s.nextID(),
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: marshal request", method)
	}

	token, err := s.signer.token(time.Now())
	if err != nil {
		return nil, errors.Wrapf(err, "%s: mint JWT", method)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.address, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrapf(err, "%s: build HTTP request", method)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	httpResp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: HTTP request failed", method)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: read response body", method)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("%s: unexpected HTTP status %d: %s", method, httpResp.StatusCode, string(body))
	}

	var rpcResp jsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, errors.Wrapf(err, "%s: decode JSON-RPC response", method)
	}

	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}

	return rpcResp.Result, nil
}
