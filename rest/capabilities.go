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
)

// ExchangeCapabilities returns the EL's supported capability names via
// GET /engine/v2/capabilities. Unlike the JSON-RPC method, the new endpoint
// is one-way: the supplied `supported` argument is ignored (the CL
// advertises its capabilities via headers / out-of-band metadata in the
// Marius design).
//
// TODO: implement once a CapabilitiesResponse container is finalised in
// spec/identification/ and the structured fields documented in
// refactor.md (per-endpoint size limits etc.) settle.
func (s *Service) ExchangeCapabilities(_ context.Context, _ []string) ([]string, error) {
	return nil, ErrNotImplemented
}
