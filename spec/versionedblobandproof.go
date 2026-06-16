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

package spec

import (
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedBlobAndProof wraps the per-fork BlobAndProof types. Cancun and
// prague use V1 (single proof per blob); osaka, amsterdam, and bogota use
// V2 (cell proofs).
type VersionedBlobAndProof struct {
	Version version.DataVersion

	Cancun    *cancun.BlobAndProof
	Prague    *cancun.BlobAndProof
	Osaka     *osaka.BlobAndProof
	Amsterdam *osaka.BlobAndProof
	Bogota    *osaka.BlobAndProof
}

// IsEmpty returns true if nothing is set for the current version.
func (v *VersionedBlobAndProof) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionCancun:
		return v.Cancun == nil
	case version.DataVersionPrague:
		return v.Prague == nil
	case version.DataVersionOsaka:
		return v.Osaka == nil
	case version.DataVersionAmsterdam:
		return v.Amsterdam == nil
	case version.DataVersionBogota:
		return v.Bogota == nil
	default:
		return true
	}
}
