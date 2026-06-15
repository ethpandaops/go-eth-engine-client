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
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// forkSegment returns the Marius URL segment for a fork (`paris`,
// `shanghai`, ...).
func forkSegment(v version.DataVersion) (string, error) {
	switch v {
	case version.DataVersionParis:
		return "paris", nil
	case version.DataVersionShanghai:
		return "shanghai", nil
	case version.DataVersionCancun:
		return "cancun", nil
	case version.DataVersionPrague:
		return "prague", nil
	case version.DataVersionOsaka:
		return "osaka", nil
	case version.DataVersionAmsterdam:
		return "amsterdam", nil
	default:
		return "", errors.Errorf("rest: unsupported fork %s", v)
	}
}

// forkURL returns "<address><basePath>/<fork>/<path>" with no trailing
// slash. The Marius spec forbids trailing slashes on the canonical form.
func (s *Service) forkURL(v version.DataVersion, suffix string) (string, error) {
	fork, err := forkSegment(v)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(s.address, "/") + s.basePath + "/" + fork + "/" + strings.Trim(suffix, "/"), nil
}

// unscopedURL returns "<address><basePath>/<path>" -- for unscoped
// endpoints (`/identity`, `/capabilities`) or the `/blobs/vN` family which
// are independently versioned (not fork-scoped).
func (s *Service) unscopedURL(suffix string) string {
	return strings.TrimRight(s.address, "/") + s.basePath + "/" + strings.Trim(suffix, "/")
}

// blobsURL returns the `/blobs/v<n>` endpoint URL.
func (s *Service) blobsURL(n int) string {
	return s.unscopedURL(fmt.Sprintf("blobs/v%d", n))
}

// payloadIDPath encodes a PayloadID as the lowercase 0x-prefixed hex form
// used in the GET /{fork}/payloads/{id} path segment.
func payloadIDPath(id paris.PayloadID) string {
	return fmt.Sprintf("%#x", id[:])
}
