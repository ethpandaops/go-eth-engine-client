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

package paris

import (
	"encoding/json"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
)

// ExecutionPayload is the payload structure exchanged by engine_newPayloadV1
// and engine_getPayloadV1. It corresponds to the SSZ container
// `ExecutionPayloadV1` defined in the SSZ transport spec.
type ExecutionPayload struct {
	ParentHash    Hash32        `ssz-size:"32"                                                       json:"parentHash"`
	FeeRecipient  Address       `ssz-size:"20"                                                       json:"feeRecipient"`
	StateRoot     Hash32        `ssz-size:"32"                                                       json:"stateRoot"`
	ReceiptsRoot  Hash32        `ssz-size:"32"                                                       json:"receiptsRoot"`
	LogsBloom     Bloom         `ssz-size:"256"                                                      json:"logsBloom"`
	PrevRandao    Hash32        `ssz-size:"32"                                                       json:"prevRandao"`
	BlockNumber   uint64        `json:"blockNumber"`
	GasLimit      uint64        `json:"gasLimit"`
	GasUsed       uint64        `json:"gasUsed"`
	Timestamp     uint64        `json:"timestamp"`
	ExtraData     []byte        `dynssz-max:"MAX_EXTRA_DATA_BYTES"                                   ssz-max:"32"                json:"extraData"`
	BaseFeePerGas *uint256.Int  `ssz-type:"uint256"                                                  json:"baseFeePerGas"`
	BlockHash     Hash32        `ssz-size:"32"                                                       json:"blockHash"`
	Transactions  []Transaction `dynssz-max:"MAX_TRANSACTIONS_PER_PAYLOAD,MAX_BYTES_PER_TRANSACTION" ssz-max:"1048576,1073741824" ssz-size:"?,?"            json:"transactions"`
}

// executionPayloadJSON is the wire-format mirror of ExecutionPayload used to
// drive Engine API JSON marshaling. Numeric fields are emitted as
// 0x-prefixed hex QUANTITY strings; byte fields use the type-level
// MarshalJSON implementations.
type executionPayloadJSON struct {
	ParentHash    Hash32                `json:"parentHash"`
	FeeRecipient  Address               `json:"feeRecipient"`
	StateRoot     Hash32                `json:"stateRoot"`
	ReceiptsRoot  Hash32                `json:"receiptsRoot"`
	LogsBloom     Bloom                 `json:"logsBloom"`
	PrevRandao    Hash32                `json:"prevRandao"`
	BlockNumber   jsonhex.QuantityU64   `json:"blockNumber"`
	GasLimit      jsonhex.QuantityU64   `json:"gasLimit"`
	GasUsed       jsonhex.QuantityU64   `json:"gasUsed"`
	Timestamp     jsonhex.QuantityU64   `json:"timestamp"`
	ExtraData     jsonhex.Bytes         `json:"extraData"`
	BaseFeePerGas *jsonhex.QuantityU256 `json:"baseFeePerGas"`
	BlockHash     Hash32                `json:"blockHash"`
	Transactions  []Transaction         `json:"transactions"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&executionPayloadJSON{
		ParentHash:    e.ParentHash,
		FeeRecipient:  e.FeeRecipient,
		StateRoot:     e.StateRoot,
		ReceiptsRoot:  e.ReceiptsRoot,
		LogsBloom:     e.LogsBloom,
		PrevRandao:    e.PrevRandao,
		BlockNumber:   jsonhex.QuantityU64(e.BlockNumber),
		GasLimit:      jsonhex.QuantityU64(e.GasLimit),
		GasUsed:       jsonhex.QuantityU64(e.GasUsed),
		Timestamp:     jsonhex.QuantityU64(e.Timestamp),
		ExtraData:     jsonhex.Bytes(e.ExtraData),
		BaseFeePerGas: (*jsonhex.QuantityU256)(e.BaseFeePerGas),
		BlockHash:     e.BlockHash,
		Transactions:  e.Transactions,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayload) UnmarshalJSON(input []byte) error {
	var data executionPayloadJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "ExecutionPayload")
	}

	if data.BaseFeePerGas == nil {
		return errors.New("ExecutionPayload: baseFeePerGas missing")
	}

	if data.Transactions == nil {
		return errors.New("ExecutionPayload: transactions missing")
	}

	if len(data.ExtraData) > MaxExtraDataBytes {
		return fmt.Errorf("ExecutionPayload: extraData longer than %d bytes", MaxExtraDataBytes)
	}

	if len(data.Transactions) > MaxTransactionsPerPayload {
		return fmt.Errorf("ExecutionPayload: too many transactions (%d > %d)",
			len(data.Transactions), MaxTransactionsPerPayload)
	}

	e.ParentHash = data.ParentHash
	e.FeeRecipient = data.FeeRecipient
	e.StateRoot = data.StateRoot
	e.ReceiptsRoot = data.ReceiptsRoot
	e.LogsBloom = data.LogsBloom
	e.PrevRandao = data.PrevRandao
	e.BlockNumber = uint64(data.BlockNumber)
	e.GasLimit = uint64(data.GasLimit)
	e.GasUsed = uint64(data.GasUsed)
	e.Timestamp = uint64(data.Timestamp)
	e.ExtraData = []byte(data.ExtraData)
	e.BaseFeePerGas = (*uint256.Int)(data.BaseFeePerGas)
	e.BlockHash = data.BlockHash
	e.Transactions = data.Transactions

	return nil
}

// String returns a JSON representation of the payload.
func (e *ExecutionPayload) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
