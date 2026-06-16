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
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// ForkchoiceUpdatedRequest carries the SSZ-transport parameters of
// engine_forkchoiceUpdatedV5. The shape mirrors V4 but with
// PayloadAttributes upgraded to V5 (`inclusionListTransactions`).
//
// The JSON-RPC variant takes a third positional `custodyColumns` parameter,
// inherited from V4. It is NOT part of the SSZ container per the Marius
// transport and so is not modelled here; the higher-level client layer
// passes [amsterdam.CustodyColumns] alongside the request when speaking
// JSON-RPC.
type ForkchoiceUpdatedRequest struct {
	ForkchoiceState   *paris.ForkchoiceState
	PayloadAttributes *PayloadAttributes `ssz-type:"optional-list"`
}
