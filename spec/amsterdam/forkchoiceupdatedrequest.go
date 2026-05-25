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

// ForkchoiceUpdatedRequest carries the SSZ-transport parameters of
// engine_forkchoiceUpdatedV4. PayloadAttributes uses V4 (slotNumber +
// targetGasLimit) and is optional, encoded as
// `List[PayloadAttributesV4, 1]`.
//
// The JSON-RPC variant of this method takes a third positional parameter,
// `custodyColumns` ([CustodyColumns] or null). It is NOT part of the SSZ
// transport's `ForkchoiceUpdatedV4Request` container per the PR764 spec
// and is therefore not modelled on this struct; the higher-level client
// layer passes [CustodyColumns] alongside the request when speaking
// JSON-RPC.
type ForkchoiceUpdatedRequest struct {
	ForkchoiceState   *paris.ForkchoiceState
	PayloadAttributes *PayloadAttributes `ssz-type:"optional-list"`
}
