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

package amsterdam

import (
	"encoding/json"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	"github.com/ethpandaops/go-eth-engine-client/internal/jsonhex"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
)

// ExecutionPayload is the payload structure exchanged by engine_newPayloadV5
// and engine_getPayloadV6's response. It corresponds to the SSZ container
// `ExecutionPayloadV4`: V3 extended with `blockAccessList` and
// `slotNumber`.
type ExecutionPayload struct {
	ParentHash      paris.Hash32           `ssz-size:"32"                                                       json:"parentHash"`
	FeeRecipient    paris.Address          `ssz-size:"20"                                                       json:"feeRecipient"`
	StateRoot       paris.Hash32           `ssz-size:"32"                                                       json:"stateRoot"`
	ReceiptsRoot    paris.Hash32           `ssz-size:"32"                                                       json:"receiptsRoot"`
	LogsBloom       paris.Bloom            `ssz-size:"256"                                                      json:"logsBloom"`
	PrevRandao      paris.Hash32           `ssz-size:"32"                                                       json:"prevRandao"`
	BlockNumber     uint64                 `json:"blockNumber"`
	GasLimit        uint64                 `json:"gasLimit"`
	GasUsed         uint64                 `json:"gasUsed"`
	Timestamp       uint64                 `json:"timestamp"`
	ExtraData       []byte                 `dynssz-max:"MAX_EXTRA_DATA_BYTES"                                   ssz-max:"32"                json:"extraData"`
	BaseFeePerGas   *uint256.Int           `ssz-type:"uint256"                                                  json:"baseFeePerGas"`
	BlockHash       paris.Hash32           `ssz-size:"32"                                                       json:"blockHash"`
	Transactions    []paris.Transaction    `dynssz-max:"MAX_TRANSACTIONS_PER_PAYLOAD,MAX_BYTES_PER_TRANSACTION" ssz-max:"1048576,1073741824" ssz-size:"?,?"            json:"transactions"`
	Withdrawals     []*shanghai.Withdrawal `dynssz-max:"MAX_WITHDRAWALS_PER_PAYLOAD"                            ssz-max:"16"                json:"withdrawals"`
	BlobGasUsed     uint64                 `json:"blobGasUsed"`
	ExcessBlobGas   uint64                 `json:"excessBlobGas"`
	BlockAccessList BlockAccessList        `dynssz-max:"MAX_BYTES_PER_TRANSACTION"                              ssz-max:"1073741824"        json:"blockAccessList"`
	SlotNumber      uint64                 `json:"slotNumber"`
}

type executionPayloadJSON struct {
	ParentHash      paris.Hash32           `json:"parentHash"`
	FeeRecipient    paris.Address          `json:"feeRecipient"`
	StateRoot       paris.Hash32           `json:"stateRoot"`
	ReceiptsRoot    paris.Hash32           `json:"receiptsRoot"`
	LogsBloom       paris.Bloom            `json:"logsBloom"`
	PrevRandao      paris.Hash32           `json:"prevRandao"`
	BlockNumber     jsonhex.QuantityU64    `json:"blockNumber"`
	GasLimit        jsonhex.QuantityU64    `json:"gasLimit"`
	GasUsed         jsonhex.QuantityU64    `json:"gasUsed"`
	Timestamp       jsonhex.QuantityU64    `json:"timestamp"`
	ExtraData       jsonhex.Bytes          `json:"extraData"`
	BaseFeePerGas   *jsonhex.QuantityU256  `json:"baseFeePerGas"`
	BlockHash       paris.Hash32           `json:"blockHash"`
	Transactions    []paris.Transaction    `json:"transactions"`
	Withdrawals     []*shanghai.Withdrawal `json:"withdrawals"`
	BlobGasUsed     jsonhex.QuantityU64    `json:"blobGasUsed"`
	ExcessBlobGas   jsonhex.QuantityU64    `json:"excessBlobGas"`
	BlockAccessList BlockAccessList        `json:"blockAccessList"`
	SlotNumber      jsonhex.QuantityU64    `json:"slotNumber"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	return json.Marshal(&executionPayloadJSON{
		ParentHash:      e.ParentHash,
		FeeRecipient:    e.FeeRecipient,
		StateRoot:       e.StateRoot,
		ReceiptsRoot:    e.ReceiptsRoot,
		LogsBloom:       e.LogsBloom,
		PrevRandao:      e.PrevRandao,
		BlockNumber:     jsonhex.QuantityU64(e.BlockNumber),
		GasLimit:        jsonhex.QuantityU64(e.GasLimit),
		GasUsed:         jsonhex.QuantityU64(e.GasUsed),
		Timestamp:       jsonhex.QuantityU64(e.Timestamp),
		ExtraData:       jsonhex.Bytes(e.ExtraData),
		BaseFeePerGas:   (*jsonhex.QuantityU256)(e.BaseFeePerGas),
		BlockHash:       e.BlockHash,
		Transactions:    e.Transactions,
		Withdrawals:     e.Withdrawals,
		BlobGasUsed:     jsonhex.QuantityU64(e.BlobGasUsed),
		ExcessBlobGas:   jsonhex.QuantityU64(e.ExcessBlobGas),
		BlockAccessList: e.BlockAccessList,
		SlotNumber:      jsonhex.QuantityU64(e.SlotNumber),
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

	if data.Withdrawals == nil {
		return errors.New("ExecutionPayload: withdrawals missing")
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
	e.Withdrawals = data.Withdrawals
	e.BlobGasUsed = uint64(data.BlobGasUsed)
	e.ExcessBlobGas = uint64(data.ExcessBlobGas)
	e.BlockAccessList = data.BlockAccessList
	e.SlotNumber = uint64(data.SlotNumber)

	return nil
}

// String returns a JSON representation.
func (e *ExecutionPayload) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
