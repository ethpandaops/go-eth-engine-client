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
	"errors"
	"fmt"
)

// ErrNotImplemented is returned by every provider method while the REST +
// SSZ transport is a scaffold.
var ErrNotImplemented = errors.New("rest: REST + SSZ transport not implemented yet (see execution-apis PR #793)")

// ProblemDetails mirrors the RFC 7807 application/problem+json error body
// the Marius spec uses for hot-path failures. Real implementations will
// populate this from the response body when the EL returns a 4xx/5xx.
type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// Error implements the error interface.
func (p *ProblemDetails) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("engine REST error %d: %s: %s", p.Status, p.Title, p.Detail)
	}

	return fmt.Sprintf("engine REST error %d: %s", p.Status, p.Title)
}
