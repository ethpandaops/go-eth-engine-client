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

// GetBlobsCellsResponse is the response from engine_getBlobsV4. Entries are
// nullable per the JSON spec (a `null` entry indicates a missing blob),
// encoded as `List[List[BlobCellsAndProofsV1, 1], MAX_BLOB_HASHES_REQUEST]`
// in SSZ.
//
// The PR #764 SSZ transport spec does not define this container; the SSZ
// schema here follows the same nullable-list convention as the other
// getBlobs responses.
type GetBlobsCellsResponse struct {
	BlobsAndProofs []*BlobCellsAndProofs `ssz-type:"list,optional-list" dynssz-max:"MAX_BLOB_HASHES_REQUEST" ssz-max:"128" json:"blobsAndProofs"`
}
