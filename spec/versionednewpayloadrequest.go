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
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedNewPayloadRequest wraps the per-fork NewPayloadRequest types.
// Osaka has no distinct NewPayload request — it reuses prague's V4 — so
// it is not present here.
type VersionedNewPayloadRequest struct {
	Version version.DataVersion

	Paris     *paris.NewPayloadRequest
	Shanghai  *shanghai.NewPayloadRequest
	Cancun    *cancun.NewPayloadRequest
	Prague    *prague.NewPayloadRequest
	Amsterdam *amsterdam.NewPayloadRequest
}

// IsEmpty returns true if no request is set for the current version.
func (v *VersionedNewPayloadRequest) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionParis:
		return v.Paris == nil
	case version.DataVersionShanghai:
		return v.Shanghai == nil
	case version.DataVersionCancun:
		return v.Cancun == nil
	case version.DataVersionPrague, version.DataVersionOsaka:
		return v.Prague == nil
	case version.DataVersionAmsterdam:
		return v.Amsterdam == nil
	default:
		return true
	}
}
