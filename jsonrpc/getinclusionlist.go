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

package jsonrpc

import (
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetInclusionList obtains the EL's local-view inclusion list for the given
// block hash (engine_getInclusionListV1, bogota+, EIP-7805). The returned
// slice contains EIP-2718-encoded transaction byte strings whose total RLP
// byte length is bounded by [bogota.MaxBytesPerInclusionList]; the EL
// selects which mempool transactions to include.
func (s *Service) GetInclusionList(
	ctx context.Context,
	blockHash paris.Hash32,
) ([]paris.Transaction, error) {
	var result []paris.Transaction
	if err := s.call(ctx, "engine_getInclusionListV1", []any{blockHash}, &result); err != nil {
		return nil, err
	}

	if result == nil {
		result = []paris.Transaction{}
	}

	return result, nil
}
