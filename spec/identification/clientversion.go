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

// Package identification holds the Engine API client-identification types
// (engine_getClientVersionV1). They are not fork-specific, so they live in
// their own package rather than a per-fork one.
package identification

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
)

// Byte-length limits for ClientVersion fields (ssz-encoding.md).
const (
	// MaxClientCodeLength is the maximum length of the client code field.
	MaxClientCodeLength = 2
	// MaxClientNameLength is the maximum length of the client name field.
	MaxClientNameLength = 64
	// MaxClientVersionLength is the maximum length of the version field.
	MaxClientVersionLength = 64
	// CommitLength is the length of the commit prefix.
	CommitLength = 4
	// MaxClientVersions is the maximum number of client versions in a
	// GetClientVersionV1Response.
	MaxClientVersions = 4
)

// ClientVersion identifies an Engine API client. It corresponds to the SSZ
// container `ClientVersionV1` and the JSON `ClientVersionV1` schema.
//
// Code, Name, and Version are UTF-8 strings on the wire (SSZ ByteLists);
// they are held as []byte so they map cleanly onto the SSZ encoding. Commit
// is the 4-byte commit prefix, hex-encoded in JSON.
type ClientVersion struct {
	Code    []byte             `dynssz-max:"MAX_CLIENT_CODE_LENGTH"    ssz-max:"2"`
	Name    []byte             `dynssz-max:"MAX_CLIENT_NAME_LENGTH"    ssz-max:"64"`
	Version []byte             `dynssz-max:"MAX_CLIENT_VERSION_LENGTH" ssz-max:"64"`
	Commit  [CommitLength]byte `ssz-size:"4"`
}

type clientVersionJSON struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// MarshalJSON implements json.Marshaler.
func (c *ClientVersion) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&clientVersionJSON{
		Code:    string(c.Code),
		Name:    string(c.Name),
		Version: string(c.Version),
		Commit:  fmt.Sprintf("%#x", c.Commit[:]),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ClientVersion) UnmarshalJSON(input []byte) error {
	var data clientVersionJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "ClientVersion")
	}

	c.Code = []byte(data.Code)
	c.Name = []byte(data.Name)
	c.Version = []byte(data.Version)

	return jsonhex.DecodeFixedText(c.Commit[:], []byte(data.Commit), CommitLength, "commit")
}

// String returns a JSON representation.
func (c *ClientVersion) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
