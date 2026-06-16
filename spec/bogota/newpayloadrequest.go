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

package bogota

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
)

// NewPayloadRequest carries the parameters of engine_newPayloadV6. It
// corresponds to the SSZ container `NewPayloadV6Request`: V5 with a trailing
// `inclusionListTransactions` field (EIP-7805). The execution payload type
// is unchanged from V5 (amsterdam's `ExecutionPayloadV4`).
type NewPayloadRequest struct {
	ExecutionPayload            *amsterdam.ExecutionPayload
	ExpectedBlobVersionedHashes []paris.Hash32            `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK" ssz-max:"4096" ssz-size:"?,32"`
	ParentBeaconBlockRoot       paris.Hash32              `ssz-size:"32"`
	ExecutionRequests           []prague.ExecutionRequest `dynssz-max:"MAX_EXECUTION_REQUESTS,MAX_BYTES_PER_TRANSACTION"     ssz-max:"256,1073741824" ssz-size:"?,?"`
	InclusionListTransactions   []paris.Transaction       `dynssz-max:"MAX_BYTES_PER_INCLUSION_LIST,MAX_BYTES_PER_INCLUSION_LIST" ssz-max:"8192,8192" ssz-size:"?,?"`
}
