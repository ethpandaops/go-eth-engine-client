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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// contentTypeSSZ is the wire content type for hot-path bodies.
	contentTypeSSZ = "application/octet-stream"
	// contentTypeJSON is used for diagnostic endpoints (/identity,
	// /capabilities).
	contentTypeJSON = "application/json"
	// contentTypeProblem is the RFC 7807 error content type.
	contentTypeProblem = "application/problem+json"
	// headerClientVersion carries the CL client version on every request
	// (replaces the `clv` JWT claim of the legacy spec).
	headerClientVersion = "X-Engine-Client-Version"
)

// authHeaders sets the headers common to every Engine API request: a fresh
// JWT bearer, the X-Engine-Client-Version header (when configured), and the
// supplied Content-Type/Accept pair.
func (s *Service) authHeaders(req *http.Request, contentType string) error {
	token, err := s.signer.token(time.Now())
	if err != nil {
		return errors.Wrap(err, "mint JWT")
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if s.clientID != "" {
		req.Header.Set(headerClientVersion, s.clientID)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Accept", contentType)
	}

	return nil
}

// postSSZ POSTs an SSZ-encoded body to the URL and returns the response
// body when status is 200, the empty slice and (true, nil) when status is
// 204 ("EL cannot serve"), or an error otherwise.
func (s *Service) postSSZ(ctx context.Context, url string, body []byte) (data []byte, noContent bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, false, errors.Wrap(err, "build POST request")
	}

	if err := s.authHeaders(req, contentTypeSSZ); err != nil {
		return nil, false, err
	}

	return s.do(req)
}

// getSSZ GETs an SSZ-encoded body from the URL. See [postSSZ] for return
// semantics.
func (s *Service) getSSZ(ctx context.Context, url string) (data []byte, noContent bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "build GET request")
	}

	if err := s.authHeaders(req, contentTypeSSZ); err != nil {
		return nil, false, err
	}

	return s.do(req)
}

// getJSON GETs a JSON-encoded body and decodes it into out.
func (s *Service) getJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "build GET request")
	}

	if err := s.authHeaders(req, contentTypeJSON); err != nil {
		return err
	}

	data, _, err := s.do(req)
	if err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(data, out)
}

// do executes the prepared request and maps 200 → (body, false, nil),
// 204 → (nil, true, nil), any error status → a [*ProblemDetails] (or raw
// status) error.
func (s *Service) do(req *http.Request) ([]byte, bool, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, false, errors.Wrapf(err, "%s %s: request failed", req.Method, req.URL)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, false, errors.Wrap(err, "read response body")
		}

		return body, false, nil
	case http.StatusNoContent:
		// Drain and discard so the connection can be reused.
		_, _ = io.Copy(io.Discard, resp.Body)

		return nil, true, nil
	default:
		body, _ := io.ReadAll(resp.Body)
		if resp.Header.Get("Content-Type") == contentTypeProblem {
			var pd ProblemDetails
			if err := json.Unmarshal(body, &pd); err == nil && pd.Status != 0 {
				return nil, false, &pd
			}
		}

		return nil, false, errors.Errorf("%s %s: HTTP %d: %s", req.Method, req.URL, resp.StatusCode, string(body))
	}
}
