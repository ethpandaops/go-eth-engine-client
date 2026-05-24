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

package paris

// ForkchoiceUpdatedRequest carries the parameters of
// engine_forkchoiceUpdatedV1. PayloadAttributes is optional; when nil, the
// call updates the forkchoice without starting a new build process.
//
// In the SSZ transport the two parameters are wrapped in the
// `ForkchoiceUpdatedV1Request` container; in JSON-RPC they are sent as the
// two positional elements of `params`.
type ForkchoiceUpdatedRequest struct {
	ForkchoiceState   *ForkchoiceState
	PayloadAttributes *PayloadAttributes
}
