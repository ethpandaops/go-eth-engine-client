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

	"github.com/ethpandaops/go-eth-engine-client/spec/identification"
)

// ClientVersion returns the EL's identity via GET /engine/v2/identity. The
// CL's identity is no longer carried in the request body -- it travels on
// every request as the X-Engine-Client-Version header, configured via
// [WithClientVersionHeader]. The argument is therefore ignored and kept
// only to satisfy the [engine.ClientVersionProvider] interface shared with
// the legacy JSON-RPC transport.
//
// TODO: implement once IdentityResponse (the
// `List[ClientVersion, MAX_CLIENT_VERSIONS]` wrapper) lands in
// spec/identification/.
func (s *Service) ClientVersion(
	_ context.Context,
	_ *identification.ClientVersion,
) ([]*identification.ClientVersion, error) {
	return nil, ErrNotImplemented
}
