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

	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// ExecutionPayloadBody is a fork-agnostic execution payload body. The
// shanghai V1 schema is reused by cancun, prague, and osaka; amsterdam adds
// a nullable block access list.
type ExecutionPayloadBody struct {
	Version version.DataVersion

	Transactions    []paris.Transaction
	Withdrawals     []*shanghai.Withdrawal
	BlockAccessList *amsterdam.BlockAccessList // amsterdam+, nullable
}

func (b *ExecutionPayloadBody) viewType() (any, error) {
	switch b.Version {
	case version.DataVersionShanghai,
		version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka:
		return (*shanghai.ExecutionPayloadBody)(nil), nil
	case version.DataVersionAmsterdam,
		version.DataVersionBogota:
		// Bogota reuses Amsterdam's V2 body shape.
		return (*amsterdam.ExecutionPayloadBody)(nil), nil
	default:
		return nil, fmt.Errorf("ExecutionPayloadBody: unsupported version %d", b.Version)
	}
}

func (b *ExecutionPayloadBody) populateVersion(v version.DataVersion) {
	b.Version = v
}

// ToView returns a fresh fork-specific ExecutionPayloadBody populated from b.
func (b *ExecutionPayloadBody) ToView() (any, error) {
	return toViewByCopy(b)
}

// FromView populates b from a fork-specific ExecutionPayloadBody.
func (b *ExecutionPayloadBody) FromView(view any) error {
	if b.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *shanghai.ExecutionPayloadBody:
			b.Version = version.DataVersionShanghai
		case *amsterdam.ExecutionPayloadBody:
			b.Version = version.DataVersionAmsterdam
		default:
			return fmt.Errorf("ExecutionPayloadBody: unsupported view type %T", view)
		}
	}

	return copyByName(view, b)
}

// ToVersioned converts b into a *spec.VersionedExecutionPayloadBody.
func (b *ExecutionPayloadBody) ToVersioned() (*spec.VersionedExecutionPayloadBody, error) {
	out := &spec.VersionedExecutionPayloadBody{}
	if err := toVersioned(b.Version, b, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates b from src.
func (b *ExecutionPayloadBody) FromVersioned(src *spec.VersionedExecutionPayloadBody) error {
	return fromVersioned(b, src)
}

// MarshalJSON implements json.Marshaler.
func (b *ExecutionPayloadBody) MarshalJSON() ([]byte, error) {
	view, err := b.ToView()
	if err != nil {
		return nil, err
	}

	return json.Marshal(view)
}

// UnmarshalJSON implements json.Unmarshaler. Version must be set first.
func (b *ExecutionPayloadBody) UnmarshalJSON(input []byte) error {
	inst, err := newViewInstance(b)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(input, inst); err != nil {
		return err
	}

	return copyByName(inst, b)
}

// MarshalSSZ marshals the body under the active fork view.
func (b *ExecutionPayloadBody) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(b, ds, make([]byte, 0, sizeSSZDyn(b, ds)))
}

// MarshalSSZTo marshals the body under the active fork view, appending to buf.
func (b *ExecutionPayloadBody) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(b, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the body under the active fork view.
func (b *ExecutionPayloadBody) SizeSSZ() int {
	return sizeSSZDyn(b, globalDynSSZ())
}

// UnmarshalSSZ decodes the body under the active fork view. Version must be
// set first.
func (b *ExecutionPayloadBody) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(b, b, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (b *ExecutionPayloadBody) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(b)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (b *ExecutionPayloadBody) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(b, globalDynSSZ(), hh)
}

// String returns the JSON representation.
func (b *ExecutionPayloadBody) String() string {
	out, err := b.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
