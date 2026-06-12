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
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// PayloadAttributes is a fork-agnostic payload-build attributes object
// holding the union of fields across every fork.
type PayloadAttributes struct {
	Version version.DataVersion

	Timestamp             uint64
	PrevRandao            paris.Hash32
	SuggestedFeeRecipient paris.Address
	Withdrawals           []*shanghai.Withdrawal // shanghai+
	ParentBeaconBlockRoot paris.Hash32           // cancun+
	SlotNumber            uint64                 // amsterdam+
	TargetGasLimit        uint64                 // amsterdam+
}

// viewType returns a typed nil pointer to the fork-specific PayloadAttributes
// for the active Version.
func (p *PayloadAttributes) viewType() (any, error) {
	switch p.Version {
	case version.DataVersionParis:
		return (*paris.PayloadAttributes)(nil), nil
	case version.DataVersionShanghai:
		return (*shanghai.PayloadAttributes)(nil), nil
	case version.DataVersionCancun,
		version.DataVersionPrague,
		version.DataVersionOsaka:
		return (*cancun.PayloadAttributes)(nil), nil
	case version.DataVersionAmsterdam:
		return (*amsterdam.PayloadAttributes)(nil), nil
	default:
		return nil, fmt.Errorf("PayloadAttributes: unsupported version %d", p.Version)
	}
}

// populateVersion sets Version.
func (p *PayloadAttributes) populateVersion(v version.DataVersion) {
	p.Version = v
}

// ToView returns a fresh fork-specific PayloadAttributes populated from p.
func (p *PayloadAttributes) ToView() (any, error) {
	return toViewByCopy(p)
}

// FromView populates p from a fork-specific PayloadAttributes.
func (p *PayloadAttributes) FromView(view any) error {
	if p.Version == version.DataVersionUnknown {
		v, err := payloadAttributesViewVersion(view)
		if err != nil {
			return err
		}

		p.Version = v
	}

	return copyByName(view, p)
}

func payloadAttributesViewVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *paris.PayloadAttributes:
		return version.DataVersionParis, nil
	case *shanghai.PayloadAttributes:
		return version.DataVersionShanghai, nil
	case *cancun.PayloadAttributes:
		return version.DataVersionCancun, nil
	case *amsterdam.PayloadAttributes:
		return version.DataVersionAmsterdam, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("PayloadAttributes: unsupported view type %T", view)
	}
}

// ToVersioned converts p into a *spec.VersionedPayloadAttributes.
func (p *PayloadAttributes) ToVersioned() (*spec.VersionedPayloadAttributes, error) {
	out := &spec.VersionedPayloadAttributes{}
	if err := toVersioned(p.Version, p, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates p from src.
func (p *PayloadAttributes) FromVersioned(src *spec.VersionedPayloadAttributes) error {
	return fromVersioned(p, src)
}

// MarshalJSON implements json.Marshaler, emitting the fork-specific form.
func (p *PayloadAttributes) MarshalJSON() ([]byte, error) {
	view, err := p.ToView()
	if err != nil {
		return nil, err
	}

	return json.Marshal(view)
}

// UnmarshalJSON implements json.Unmarshaler. Version must be set first.
func (p *PayloadAttributes) UnmarshalJSON(input []byte) error {
	inst, err := newViewInstance(p)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(input, inst); err != nil {
		return err
	}

	return copyByName(inst, p)
}

// MarshalSSZ marshals the attributes under the active fork view.
func (p *PayloadAttributes) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(p, ds, make([]byte, 0, sizeSSZDyn(p, ds)))
}

// MarshalSSZTo marshals the attributes under the active fork view, appending to buf.
func (p *PayloadAttributes) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(p, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the attributes under the active fork view.
func (p *PayloadAttributes) SizeSSZ() int {
	return sizeSSZDyn(p, globalDynSSZ())
}

// UnmarshalSSZ decodes the attributes under the active fork view. Version
// must be set first.
func (p *PayloadAttributes) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(p, p, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (p *PayloadAttributes) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(p)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (p *PayloadAttributes) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(p, globalDynSSZ(), hh)
}

// String returns the JSON representation.
func (p *PayloadAttributes) String() string {
	out, err := p.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
