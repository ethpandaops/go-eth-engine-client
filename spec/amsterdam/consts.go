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

// SSZ list-size and byte-length limits introduced in Amsterdam.
const (
	// CustodyColumnsLength is the byte length of the
	// custodyColumns bit-array sent on engine_forkchoiceUpdatedV4 as a
	// 3rd JSON-RPC parameter (CELLS_PER_EXT_BLOB / 8 = 16 bytes).
	CustodyColumnsLength = 16
	// MaxBlockAccessListBytes mirrors MAX_BYTES_PER_TRANSACTION as used
	// for the RLP-encoded block access list (EIP-7928).
	MaxBlockAccessListBytes = 1 << 30
	// CellsPerExtBlob is the number of cells in an extended blob (EIP-7594).
	CellsPerExtBlob = 128
	// CellLength is the size in bytes of a single blob cell: 64 field
	// elements of 32 bytes each (EIP-7594).
	CellLength = 64 * 32
	// MaxBlobHashesRequest is the maximum number of blob versioned hashes
	// in an engine_getBlobs request.
	MaxBlobHashesRequest = 128
)
