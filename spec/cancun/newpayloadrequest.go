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

package cancun

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// NewPayloadRequest carries the parameters of engine_newPayloadV3. It
// corresponds to the SSZ container `NewPayloadV3Request` and is the only
// new-payload variant with multiple parameters beyond the execution payload.
type NewPayloadRequest struct {
	ExecutionPayload            *ExecutionPayload
	ExpectedBlobVersionedHashes []paris.Hash32 `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK" ssz-max:"4096" ssz-size:"?,32"`
	ParentBeaconBlockRoot       paris.Hash32   `ssz-size:"32"`
}
