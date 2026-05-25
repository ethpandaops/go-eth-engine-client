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

package osaka

// SSZ list-size limits introduced in Osaka (EIP-7594 cell proofs).
const (
	// CellsPerExtBlob is the number of cells in an extended blob.
	CellsPerExtBlob = 128
	// MaxBlobProofsBundle is the maximum total number of cell proofs in a
	// BlobsBundleV2 (MAX_BLOB_COMMITMENTS_PER_BLOCK * CELLS_PER_EXT_BLOB).
	MaxBlobProofsBundle = 4096 * 128
)
