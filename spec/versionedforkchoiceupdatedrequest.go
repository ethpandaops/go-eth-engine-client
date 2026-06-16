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
	"github.com/ethpandaops/go-eth-engine-client/spec/bogota"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedForkchoiceUpdatedRequest wraps the per-fork
// ForkchoiceUpdatedRequest types. Prague and Osaka reuse cancun's V3.
// Bogota's V5 mirrors amsterdam's V4 but with PayloadAttributesV5 in place
// of V4 (adds `inclusionListTransactions`).
//
// CustodyColumns is the third JSON-RPC parameter introduced in
// engine_forkchoiceUpdatedV4 and inherited by V5. It has no SSZ
// representation in the per-fork request containers, so it is carried
// alongside them here for the JSON-RPC client to consume when Version is
// amsterdam or bogota.
type VersionedForkchoiceUpdatedRequest struct {
	Version version.DataVersion

	Paris          *paris.ForkchoiceUpdatedRequest
	Shanghai       *shanghai.ForkchoiceUpdatedRequest
	Cancun         *cancun.ForkchoiceUpdatedRequest
	Prague         *cancun.ForkchoiceUpdatedRequest
	Osaka          *cancun.ForkchoiceUpdatedRequest
	Amsterdam      *amsterdam.ForkchoiceUpdatedRequest
	Bogota         *bogota.ForkchoiceUpdatedRequest
	CustodyColumns *amsterdam.CustodyColumns
}

// IsEmpty returns true if no request is set for the current version.
func (v *VersionedForkchoiceUpdatedRequest) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionParis:
		return v.Paris == nil
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
	case version.DataVersionBogota:
		return v.Bogota == nil
	default:
		return true
	}
}
