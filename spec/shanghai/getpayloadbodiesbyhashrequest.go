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

package shanghai

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetPayloadBodiesByHashRequest carries the parameters of
// engine_getPayloadBodiesByHashV1. Used by both JSON-RPC and the SSZ
// transport's `GetPayloadBodiesByHashV1Request` container.
type GetPayloadBodiesByHashRequest struct {
	BlockHashes []paris.Hash32 `dynssz-max:"MAX_PAYLOAD_BODIES_REQUEST" ssz-max:"32" ssz-size:"?,32"`
}
