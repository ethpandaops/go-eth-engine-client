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

// VersionedGetBlobsRequest wraps the per-fork all-or-nothing GetBlobsRequest
// types: cancun's V1 and osaka's V2 (identical shape, distinct methods).
type VersionedGetBlobsRequest struct {
	Version version.DataVersion

	Cancun *cancun.GetBlobsRequest
	Osaka  *osaka.GetBlobsRequest
}

// IsEmpty returns true if no request is set for the current version.
func (v *VersionedGetBlobsRequest) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionCancun:
		return v.Cancun == nil
	case version.DataVersionOsaka:
		return v.Osaka == nil
	default:
		return true
	}
}
