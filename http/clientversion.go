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

package http

import (
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec/identification"
)

// ClientVersion exchanges client identity via engine_getClientVersionV1,
// sending the consensus client's version and returning the execution
// client's version information.
func (s *Service) ClientVersion(
	ctx context.Context,
	clientVersion *identification.ClientVersion,
) ([]*identification.ClientVersion, error) {
	if clientVersion == nil {
		clientVersion = &identification.ClientVersion{}
	}

	var out []*identification.ClientVersion
	if err := s.call(ctx, "engine_getClientVersionV1", []any{clientVersion}, &out); err != nil {
		return nil, err
	}

	return out, nil
}
