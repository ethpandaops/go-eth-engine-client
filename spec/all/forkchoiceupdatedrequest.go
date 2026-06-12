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
	"fmt"

	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// ForkchoiceUpdatedRequest is a fork-agnostic engine_forkchoiceUpdated
// request. The ForkchoiceState shape is constant across forks; the optional
// PayloadAttributes is a fork-agnostic union.
//
// CustodyColumns is the amsterdam-only third JSON-RPC parameter to
// engine_forkchoiceUpdatedV4. It is not part of the SSZ transport's
// ForkchoiceUpdatedV4Request container, so it is excluded from SSZ marshaling
// (no fork view defines a matching field) and is only consumed by the
// JSON-RPC client when Version is amsterdam.
type ForkchoiceUpdatedRequest struct {
	Version version.DataVersion

	ForkchoiceState   *paris.ForkchoiceState
	PayloadAttributes *PayloadAttributes
	CustodyColumns    *amsterdam.CustodyColumns
}

func (r *ForkchoiceUpdatedRequest) viewType() (any, error) {
	switch r.Version {
	case version.DataVersionParis:
		return (*paris.ForkchoiceUpdatedRequest)(nil), nil
	case version.DataVersionShanghai:
		return (*shanghai.ForkchoiceUpdatedRequest)(nil), nil
	case version.DataVersionCancun, version.DataVersionPrague, version.DataVersionOsaka:
		return (*cancun.ForkchoiceUpdatedRequest)(nil), nil
	case version.DataVersionAmsterdam:
		return (*amsterdam.ForkchoiceUpdatedRequest)(nil), nil
	default:
		return nil, fmt.Errorf("ForkchoiceUpdatedRequest: unsupported version %d", r.Version)
	}
}

func (r *ForkchoiceUpdatedRequest) populateVersion(v version.DataVersion) {
	r.Version = v

	if r.PayloadAttributes != nil {
		r.PayloadAttributes.populateVersion(v)
	}
}

// ToView returns a fresh fork-specific ForkchoiceUpdatedRequest from r.
func (r *ForkchoiceUpdatedRequest) ToView() (any, error) {
	return toViewByCopy(r)
}

// FromView populates r from a fork-specific ForkchoiceUpdatedRequest.
func (r *ForkchoiceUpdatedRequest) FromView(view any) error {
	if r.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *paris.ForkchoiceUpdatedRequest:
			r.Version = version.DataVersionParis
		case *shanghai.ForkchoiceUpdatedRequest:
			r.Version = version.DataVersionShanghai
		case *cancun.ForkchoiceUpdatedRequest:
			r.Version = version.DataVersionCancun
		case *amsterdam.ForkchoiceUpdatedRequest:
			r.Version = version.DataVersionAmsterdam
		default:
			return fmt.Errorf("ForkchoiceUpdatedRequest: unsupported view type %T", view)
		}
	}

	if err := copyByName(view, r); err != nil {
		return err
	}

	r.populateVersion(r.Version)

	return nil
}

// ToVersioned converts r into a *spec.VersionedForkchoiceUpdatedRequest.
func (r *ForkchoiceUpdatedRequest) ToVersioned() (*spec.VersionedForkchoiceUpdatedRequest, error) {
	out := &spec.VersionedForkchoiceUpdatedRequest{}
	if err := toVersioned(r.Version, r, out); err != nil {
		return nil, err
	}

	// CustodyColumns is not part of any fork view, so the generic
	// reflection conversion skips it; carry it across explicitly.
	out.CustodyColumns = r.CustodyColumns

	return out, nil
}

// FromVersioned populates r from src.
func (r *ForkchoiceUpdatedRequest) FromVersioned(src *spec.VersionedForkchoiceUpdatedRequest) error {
	if err := fromVersioned(r, src); err != nil {
		return err
	}

	r.CustodyColumns = src.CustodyColumns

	return nil
}

// MarshalSSZ marshals the request under the active fork view.
func (r *ForkchoiceUpdatedRequest) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(r, ds, make([]byte, 0, sizeSSZDyn(r, ds)))
}

// MarshalSSZTo marshals the request under the active fork view, appending to buf.
func (r *ForkchoiceUpdatedRequest) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(r, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the request under the active fork view.
func (r *ForkchoiceUpdatedRequest) SizeSSZ() int {
	return sizeSSZDyn(r, globalDynSSZ())
}

// UnmarshalSSZ decodes the request under the active fork view. Version must
// be set first.
func (r *ForkchoiceUpdatedRequest) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(r, r, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (r *ForkchoiceUpdatedRequest) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(r)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (r *ForkchoiceUpdatedRequest) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(r, globalDynSSZ(), hh)
}
