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

// Package bogota holds the Engine API spec containers introduced in the
// Bogota fork (EIP-7805 inclusion lists). The new wire-format change is the
// `inclusionListTransactions` field, which lands as a 5th positional
// parameter on engine_newPayloadV6 and as a new field on
// PayloadAttributesV5 (consumed by engine_forkchoiceUpdatedV5).
//
// The execution payload, payload body, blobs bundle, getPayload response,
// and block-access-list containers are unchanged from Amsterdam and are
// reused via the spec/amsterdam package; only the structures that grow a
// new field are redefined here.
package bogota

// SSZ list-size and byte-length limits introduced in Bogota.
const (
	// MaxBytesPerInclusionList caps the total byte length of the RLP
	// encoding of the transactions returned in an inclusion list and the
	// `inclusionListTransactions` field on engine_newPayloadV6 /
	// PayloadAttributesV5 (EIP-7805: `MAX_BYTES_PER_INCLUSION_LIST`).
	MaxBytesPerInclusionList = 1 << 13
)
