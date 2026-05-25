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

package prague

import (
	"encoding/json"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
)

// GetPayloadResponse is the response from engine_getPayloadV4. It
// corresponds to the SSZ container `GetPayloadResponseV4`: V3 extended with
// `executionRequests`. The execution payload remains ExecutionPayloadV3
// (cancun); BlobsBundle is still V1.
type GetPayloadResponse struct {
	ExecutionPayload      *cancun.ExecutionPayload `json:"executionPayload"`
	BlockValue            *uint256.Int             `ssz-type:"uint256" json:"blockValue"`
	BlobsBundle           *cancun.BlobsBundle      `json:"blobsBundle"`
	ShouldOverrideBuilder bool                     `json:"shouldOverrideBuilder"`
	ExecutionRequests     []ExecutionRequest       `dynssz-max:"MAX_EXECUTION_REQUESTS,MAX_BYTES_PER_TRANSACTION" ssz-max:"256,1073741824" ssz-size:"?,?" json:"executionRequests"`
}

type getPayloadResponseJSON struct {
	ExecutionPayload      *cancun.ExecutionPayload `json:"executionPayload"`
	BlockValue            *jsonhex.QuantityU256    `json:"blockValue"`
	BlobsBundle           *cancun.BlobsBundle      `json:"blobsBundle"`
	ShouldOverrideBuilder bool                     `json:"shouldOverrideBuilder"`
	ExecutionRequests     []ExecutionRequest       `json:"executionRequests"`
}

// MarshalJSON implements json.Marshaler.
func (g *GetPayloadResponse) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&getPayloadResponseJSON{
		ExecutionPayload:      g.ExecutionPayload,
		BlockValue:            (*jsonhex.QuantityU256)(g.BlockValue),
		BlobsBundle:           g.BlobsBundle,
		ShouldOverrideBuilder: g.ShouldOverrideBuilder,
		ExecutionRequests:     g.ExecutionRequests,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (g *GetPayloadResponse) UnmarshalJSON(input []byte) error {
	var data getPayloadResponseJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "GetPayloadResponse")
	}

	if data.ExecutionPayload == nil {
		return errors.New("GetPayloadResponse: executionPayload missing")
	}

	if data.BlockValue == nil {
		return errors.New("GetPayloadResponse: blockValue missing")
	}

	if data.BlobsBundle == nil {
		return errors.New("GetPayloadResponse: blobsBundle missing")
	}

	if data.ExecutionRequests == nil {
		return errors.New("GetPayloadResponse: executionRequests missing")
	}

	g.ExecutionPayload = data.ExecutionPayload
	g.BlockValue = (*uint256.Int)(data.BlockValue)
	g.BlobsBundle = data.BlobsBundle
	g.ShouldOverrideBuilder = data.ShouldOverrideBuilder
	g.ExecutionRequests = data.ExecutionRequests

	return nil
}

// String returns a JSON representation.
func (g *GetPayloadResponse) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
