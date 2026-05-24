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

// SSZ list-size limit introduced in Shanghai.
const (
	// MaxWithdrawalsPerPayload is the maximum number of withdrawals in a
	// single ExecutionPayload (2**4, Capella).
	MaxWithdrawalsPerPayload = 16
	// MaxPayloadBodiesRequest is the maximum number of blocks in a
	// getPayloadBodiesByHash / getPayloadBodiesByRange request.
	MaxPayloadBodiesRequest = 32
)
