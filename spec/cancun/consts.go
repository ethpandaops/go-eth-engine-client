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

// Byte-length constants for the Cancun KZG / blob types.
const (
	// KZGCommitmentLength is the number of bytes in a KZG commitment.
	KZGCommitmentLength = 48
	// KZGProofLength is the number of bytes in a KZG proof.
	KZGProofLength = 48
	// BytesPerFieldElement is the size in bytes of a single field element.
	BytesPerFieldElement = 32
	// FieldElementsPerBlob is the number of field elements per blob.
	FieldElementsPerBlob = 4096
	// BlobLength is the size in bytes of a single blob.
	BlobLength = BytesPerFieldElement * FieldElementsPerBlob
)

// SSZ list-size limits introduced in Cancun.
const (
	// MaxBlobCommitmentsPerBlock is the maximum number of blob KZG
	// commitments included in a single block.
	MaxBlobCommitmentsPerBlock = 4096
	// MaxBlobHashesRequest is the maximum number of blob versioned hashes
	// in a single engine_getBlobs request.
	MaxBlobHashesRequest = 128
)
