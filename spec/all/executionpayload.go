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

package all

import (
	"encoding/json"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// ExecutionPayload is a fork-agnostic execution payload holding the union of
// fields across every fork. The fields populated on a given instance depend
// on Version.
type ExecutionPayload struct {
	Version version.DataVersion

	ParentHash      paris.Hash32
	FeeRecipient    paris.Address
	StateRoot       paris.Hash32
	ReceiptsRoot    paris.Hash32
	LogsBloom       paris.Bloom
	PrevRandao      paris.Hash32
	BlockNumber     uint64
	GasLimit        uint64
	GasUsed         uint64
	Timestamp       uint64
	ExtraData       []byte
	BaseFeePerGas   *uint256.Int
	BlockHash       paris.Hash32
	Transactions    []paris.Transaction
	Withdrawals     []*shanghai.Withdrawal // shanghai+
	BlobGasUsed     uint64                 // cancun+
	ExcessBlobGas   uint64                 // cancun+
	BlockAccessList amsterdam.BlockAccessList
	SlotNumber      uint64 // amsterdam+
}

// viewType returns a typed nil pointer to the fork-specific ExecutionPayload
// for the active Version.
func (e *ExecutionPayload) viewType() (any, error) {
	switch e.Version {
	case version.DataVersionParis:
		return (*paris.ExecutionPayload)(nil), nil
	case version.DataVersionShanghai:
		return (*shanghai.ExecutionPayload)(nil), nil
	case version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka:
		// Prague and Osaka reuse the Cancun execution-payload schema.
		return (*cancun.ExecutionPayload)(nil), nil
	case version.DataVersionAmsterdam:
		return (*amsterdam.ExecutionPayload)(nil), nil
	default:
		return nil, fmt.Errorf("ExecutionPayload: unsupported version %d", e.Version)
	}
}

// populateVersion sets Version.
func (e *ExecutionPayload) populateVersion(v version.DataVersion) {
	e.Version = v
}

// ToView returns a fresh fork-specific ExecutionPayload populated from e.
func (e *ExecutionPayload) ToView() (any, error) {
	return toViewByCopy(e)
}

// FromView populates e from a fork-specific ExecutionPayload.
func (e *ExecutionPayload) FromView(view any) error {
	if e.Version == version.DataVersionUnknown {
		v, err := executionPayloadViewVersion(view)
		if err != nil {
			return err
		}

		e.Version = v
	}

	return copyByName(view, e)
}

// executionPayloadViewVersion maps a view type to its DataVersion. Cancun is
// returned for the shared cancun schema; callers that know the precise fork
// (prague/osaka) should pin Version before FromView.
func executionPayloadViewVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *paris.ExecutionPayload:
		return version.DataVersionParis, nil
	case *shanghai.ExecutionPayload:
		return version.DataVersionShanghai, nil
	case *cancun.ExecutionPayload:
		return version.DataVersionCancun, nil
	case *amsterdam.ExecutionPayload:
		return version.DataVersionAmsterdam, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("ExecutionPayload: unsupported view type %T", view)
	}
}

// ToVersioned converts e into a *spec.VersionedExecutionPayload.
func (e *ExecutionPayload) ToVersioned() (*spec.VersionedExecutionPayload, error) {
	out := &spec.VersionedExecutionPayload{}
	if err := toVersioned(e.Version, e, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates e from src.
func (e *ExecutionPayload) FromVersioned(src *spec.VersionedExecutionPayload) error {
	return fromVersioned(e, src)
}

// MarshalJSON implements json.Marshaler, emitting the fork-specific form.
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	view, err := e.ToView()
	if err != nil {
		return nil, err
	}

	return json.Marshal(view)
}

// UnmarshalJSON implements json.Unmarshaler. Version must be set first so the
// correct fork view is used.
func (e *ExecutionPayload) UnmarshalJSON(input []byte) error {
	inst, err := newViewInstance(e)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(input, inst); err != nil {
		return err
	}

	return copyByName(inst, e)
}

// MarshalSSZ marshals the payload under the active fork view.
func (e *ExecutionPayload) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(e, ds, make([]byte, 0, sizeSSZDyn(e, ds)))
}

// MarshalSSZTo marshals the payload under the active fork view, appending to buf.
func (e *ExecutionPayload) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(e, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the payload under the active fork view.
func (e *ExecutionPayload) SizeSSZ() int {
	return sizeSSZDyn(e, globalDynSSZ())
}

// UnmarshalSSZ decodes the payload under the active fork view. Version must
// be set first.
func (e *ExecutionPayload) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(e, e, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (e *ExecutionPayload) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(e)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (e *ExecutionPayload) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(e, globalDynSSZ(), hh)
}

// String returns the JSON representation.
func (e *ExecutionPayload) String() string {
	out, err := e.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
