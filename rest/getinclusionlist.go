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
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetInclusionList returns the EL's local-view inclusion list for the given
// block hash (bogota+, EIP-7805). The REST endpoint will live at
// GET /engine/v2/bogota/inclusion-list/{blockHash}, returning a bare
// SSZ `List[Transaction, MAX_BYTES_PER_INCLUSION_LIST/1]` over
// application/octet-stream. An HTTP 204 means "EL declines to serve" and
// surfaces as an empty (non-nil) slice.
//
// TODO: implement once the Marius spec finalises the endpoint path and
// SSZ container for inclusion-list responses (currently only specified in
// JSON-RPC form by execution-apis PR #609).
func (s *Service) GetInclusionList(
	_ context.Context,
	_ paris.Hash32,
) ([]paris.Transaction, error) {
	return nil, ErrNotImplemented
}
