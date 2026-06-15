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

package rest

import (
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
)

// ForkchoiceUpdated updates the fork choice and optionally starts a payload
// build via POST /engine/v2/{fork}/forkchoice with an SSZ-encoded
// ForkchoiceUpdate body. The Marius envelope folds the (Amsterdam+) custody
// columns into the SSZ container, so unlike the JSON-RPC path the custody
// set travels alongside the rest of the request, not as a side-channel
// field on [spec.VersionedForkchoiceUpdatedRequest].
//
// TODO: implement once amsterdam.ForkchoiceUpdatedRequest carries
// CustodyColumns as an optional SSZ field and the side-channel field is
// removed from [spec.VersionedForkchoiceUpdatedRequest].
func (s *Service) ForkchoiceUpdated(
	_ context.Context,
	_ *spec.VersionedForkchoiceUpdatedRequest,
) (*paris.ForkchoiceUpdatedResponse, error) {
	return nil, ErrNotImplemented
}

// ForkchoiceUpdatedAgnostic is the fork-agnostic-typed variant.
func (s *Service) ForkchoiceUpdatedAgnostic(
	_ context.Context,
	_ *all.ForkchoiceUpdatedRequest,
) (*paris.ForkchoiceUpdatedResponse, error) {
	return nil, ErrNotImplemented
}
