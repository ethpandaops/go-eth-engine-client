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

// Package spec exposes fork-agnostic Versioned* wrappers over the per-fork
// Engine API spec types in subpackages (paris, shanghai, cancun, prague,
// osaka, amsterdam). Each Versioned* carries a [version.DataVersion] and a
// pointer per fork, with accessors that switch on the version to return
// fields that exist across forks.
package spec

import (
	"errors"

	"github.com/holiman/uint256"

	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedExecutionPayload contains an execution payload that may be from
// any fork that defines an ExecutionPayload shape: paris, shanghai, cancun
// (also used through prague + osaka), or amsterdam.
type VersionedExecutionPayload struct {
	Version version.DataVersion

	Paris     *paris.ExecutionPayload
	Shanghai  *shanghai.ExecutionPayload
	Cancun    *cancun.ExecutionPayload
	Prague    *cancun.ExecutionPayload
	Osaka     *cancun.ExecutionPayload
	Amsterdam *amsterdam.ExecutionPayload
}

// IsEmpty returns true if no payload is set for the current version.
func (v *VersionedExecutionPayload) IsEmpty() bool {
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
	default:
		return true
	}
}

// ParentHash returns the parent hash of the payload.
func (v *VersionedExecutionPayload) ParentHash() (paris.Hash32, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return paris.Hash32{}, errors.New("no paris payload")
		}

		return v.Paris.ParentHash, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return paris.Hash32{}, errors.New("no shanghai payload")
		}

		return v.Shanghai.ParentHash, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return paris.Hash32{}, errors.New("no cancun payload")
		}

		return v.Cancun.ParentHash, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return paris.Hash32{}, errors.New("no prague payload")
		}

		return v.Prague.ParentHash, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return paris.Hash32{}, errors.New("no osaka payload")
		}

		return v.Osaka.ParentHash, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return paris.Hash32{}, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.ParentHash, nil
	default:
		return paris.Hash32{}, errors.New("unknown version")
	}
}

// BlockHash returns the block hash of the payload.
func (v *VersionedExecutionPayload) BlockHash() (paris.Hash32, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return paris.Hash32{}, errors.New("no paris payload")
		}

		return v.Paris.BlockHash, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return paris.Hash32{}, errors.New("no shanghai payload")
		}

		return v.Shanghai.BlockHash, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return paris.Hash32{}, errors.New("no cancun payload")
		}

		return v.Cancun.BlockHash, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return paris.Hash32{}, errors.New("no prague payload")
		}

		return v.Prague.BlockHash, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return paris.Hash32{}, errors.New("no osaka payload")
		}

		return v.Osaka.BlockHash, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return paris.Hash32{}, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.BlockHash, nil
	default:
		return paris.Hash32{}, errors.New("unknown version")
	}
}

// BlockNumber returns the block number of the payload.
func (v *VersionedExecutionPayload) BlockNumber() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return 0, errors.New("no paris payload")
		}

		return v.Paris.BlockNumber, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return 0, errors.New("no shanghai payload")
		}

		return v.Shanghai.BlockNumber, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return 0, errors.New("no cancun payload")
		}

		return v.Cancun.BlockNumber, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return 0, errors.New("no prague payload")
		}

		return v.Prague.BlockNumber, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return 0, errors.New("no osaka payload")
		}

		return v.Osaka.BlockNumber, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Timestamp returns the timestamp of the payload.
func (v *VersionedExecutionPayload) Timestamp() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return 0, errors.New("no paris payload")
		}

		return v.Paris.Timestamp, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return 0, errors.New("no shanghai payload")
		}

		return v.Shanghai.Timestamp, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return 0, errors.New("no cancun payload")
		}

		return v.Cancun.Timestamp, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return 0, errors.New("no prague payload")
		}

		return v.Prague.Timestamp, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return 0, errors.New("no osaka payload")
		}

		return v.Osaka.Timestamp, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.Timestamp, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BaseFeePerGas returns the base fee per gas of the payload.
func (v *VersionedExecutionPayload) BaseFeePerGas() (*uint256.Int, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return nil, errors.New("no paris payload")
		}

		return v.Paris.BaseFeePerGas, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return nil, errors.New("no shanghai payload")
		}

		return v.Shanghai.BaseFeePerGas, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return nil, errors.New("no cancun payload")
		}

		return v.Cancun.BaseFeePerGas, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return nil, errors.New("no prague payload")
		}

		return v.Prague.BaseFeePerGas, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return nil, errors.New("no osaka payload")
		}

		return v.Osaka.BaseFeePerGas, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return nil, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.BaseFeePerGas, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Transactions returns the transactions of the payload.
func (v *VersionedExecutionPayload) Transactions() ([]paris.Transaction, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return nil, errors.New("no paris payload")
		}

		return v.Paris.Transactions, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return nil, errors.New("no shanghai payload")
		}

		return v.Shanghai.Transactions, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return nil, errors.New("no cancun payload")
		}

		return v.Cancun.Transactions, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return nil, errors.New("no prague payload")
		}

		return v.Prague.Transactions, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return nil, errors.New("no osaka payload")
		}

		return v.Osaka.Transactions, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return nil, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.Transactions, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals of the payload (shanghai+).
func (v *VersionedExecutionPayload) Withdrawals() ([]*shanghai.Withdrawal, error) {
	switch v.Version {
	case version.DataVersionParis:
		return nil, errors.New("no withdrawals in paris")
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return nil, errors.New("no shanghai payload")
		}

		return v.Shanghai.Withdrawals, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return nil, errors.New("no cancun payload")
		}

		return v.Cancun.Withdrawals, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return nil, errors.New("no prague payload")
		}

		return v.Prague.Withdrawals, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return nil, errors.New("no osaka payload")
		}

		return v.Osaka.Withdrawals, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return nil, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BlobGasUsed returns the blob gas used of the payload (cancun+).
func (v *VersionedExecutionPayload) BlobGasUsed() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis, version.DataVersionShanghai:
		return 0, errors.New("no blob gas before cancun")
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return 0, errors.New("no cancun payload")
		}

		return v.Cancun.BlobGasUsed, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return 0, errors.New("no prague payload")
		}

		return v.Prague.BlobGasUsed, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return 0, errors.New("no osaka payload")
		}

		return v.Osaka.BlobGasUsed, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.BlobGasUsed, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExcessBlobGas returns the excess blob gas of the payload (cancun+).
func (v *VersionedExecutionPayload) ExcessBlobGas() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis, version.DataVersionShanghai:
		return 0, errors.New("no excess blob gas before cancun")
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return 0, errors.New("no cancun payload")
		}

		return v.Cancun.ExcessBlobGas, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return 0, errors.New("no prague payload")
		}

		return v.Prague.ExcessBlobGas, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return 0, errors.New("no osaka payload")
		}

		return v.Osaka.ExcessBlobGas, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.ExcessBlobGas, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BlockAccessList returns the block access list of the payload
// (amsterdam+).
func (v *VersionedExecutionPayload) BlockAccessList() (amsterdam.BlockAccessList, error) {
	switch v.Version {
	case version.DataVersionParis, version.DataVersionShanghai,
		version.DataVersionCancun, version.DataVersionPrague,
		version.DataVersionOsaka:
		return nil, errors.New("no block access list before amsterdam")
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return nil, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.BlockAccessList, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// SlotNumber returns the slot number of the payload (amsterdam+).
func (v *VersionedExecutionPayload) SlotNumber() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis, version.DataVersionShanghai,
		version.DataVersionCancun, version.DataVersionPrague,
		version.DataVersionOsaka:
		return 0, errors.New("no slot number before amsterdam")
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam payload")
		}

		return v.Amsterdam.SlotNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// String returns a JSON representation of the active payload.
func (v *VersionedExecutionPayload) String() string {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return ""
		}

		return v.Paris.String()
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return ""
		}

		return v.Shanghai.String()
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return ""
		}

		return v.Cancun.String()
	case version.DataVersionPrague:
		if v.Prague == nil {
			return ""
		}

		return v.Prague.String()
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return ""
		}

		return v.Osaka.String()
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return ""
		}

		return v.Amsterdam.String()
	default:
		return "unknown version"
	}
}
