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

package amsterdam

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// GetBlobsCellsRequest carries the parameters of engine_getBlobsV4 — the
// cell-based variant that fetches individual blob cells (rather than full
// blobs) and supports partial responses.
//
// JSON-RPC only at the time of writing; PR #764 does not yet define a
// matching SSZ container.
type GetBlobsCellsRequest struct {
	VersionedBlobHashes []paris.Hash32 `json:"versioned_blob_hashes"`
	IndicesBitarray     [16]byte       `json:"indices_bitarray"`
}
