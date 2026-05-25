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
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedGetPayloadResponse wraps the per-fork engine_getPayload response
// containers. Paris (V1) returns an ExecutionPayload directly without a
// wrapper container, so it has no entry here.
type VersionedGetPayloadResponse struct {
	Version version.DataVersion

	Shanghai  *shanghai.GetPayloadResponse
	Cancun    *cancun.GetPayloadResponse
	Prague    *prague.GetPayloadResponse
	Osaka     *osaka.GetPayloadResponse
	Amsterdam *amsterdam.GetPayloadResponse
}

// IsEmpty returns true if no response is set for the current version.
func (v *VersionedGetPayloadResponse) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionShanghai:
		return v.Shanghai == nil
	case version.DataVersionCancun:
		return v.Cancun == nil
	case version.DataVersionPrague:
		return v.Prague == nil
	case version.DataVersionOsaka:
		return v.Osaka == nil
	case version.DataVersionAmsterdam:
		return v.Amsterdam == nil
	default:
		return true
	}
}
