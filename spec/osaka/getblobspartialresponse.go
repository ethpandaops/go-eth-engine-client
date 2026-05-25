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

// GetBlobsPartialResponse is the response from engine_getBlobsV3 — the
// variant that supports partial responses. It corresponds to the SSZ
// container `GetBlobsV3Response`:
// `List[List[BlobAndProofV2, 1], MAX_BLOB_HASHES_REQUEST]`. A missing
// blob is a nil entry in BlobsAndProofs; the canonical optional-list
// (List[T, 1]) form encodes the inner pointer.
type GetBlobsPartialResponse struct {
	BlobsAndProofs []*BlobAndProof `dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" ssz-type:"list,optional-list" json:"blobsAndProofs"`
}
