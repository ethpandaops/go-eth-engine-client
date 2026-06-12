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
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// NewPayloadRequest is a fork-agnostic engine_newPayload request. The
// nested ExecutionPayload is itself a fork-agnostic union; ExpectedBlob*,
// ParentBeaconBlockRoot, and ExecutionRequests appear from cancun and
// prague respectively.
type NewPayloadRequest struct {
	Version version.DataVersion

	ExecutionPayload            *ExecutionPayload
	ExpectedBlobVersionedHashes []paris.Hash32            // cancun+
	ParentBeaconBlockRoot       paris.Hash32              // cancun+
	ExecutionRequests           []prague.ExecutionRequest // prague+
}

func (r *NewPayloadRequest) viewType() (any, error) {
	switch r.Version {
	case version.DataVersionParis:
		return (*paris.NewPayloadRequest)(nil), nil
	case version.DataVersionShanghai:
		return (*shanghai.NewPayloadRequest)(nil), nil
	case version.DataVersionCancun:
		return (*cancun.NewPayloadRequest)(nil), nil
	case version.DataVersionPrague, version.DataVersionOsaka:
		return (*prague.NewPayloadRequest)(nil), nil
	case version.DataVersionAmsterdam:
		return (*amsterdam.NewPayloadRequest)(nil), nil
	default:
		return nil, fmt.Errorf("NewPayloadRequest: unsupported version %d", r.Version)
	}
}

func (r *NewPayloadRequest) populateVersion(v version.DataVersion) {
	r.Version = v

	if r.ExecutionPayload != nil {
		r.ExecutionPayload.populateVersion(v)
	}
}

// ToView returns a fresh fork-specific NewPayloadRequest populated from r.
func (r *NewPayloadRequest) ToView() (any, error) {
	return toViewByCopy(r)
}

// FromView populates r from a fork-specific NewPayloadRequest.
func (r *NewPayloadRequest) FromView(view any) error {
	if r.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *paris.NewPayloadRequest:
			r.Version = version.DataVersionParis
		case *shanghai.NewPayloadRequest:
			r.Version = version.DataVersionShanghai
		case *cancun.NewPayloadRequest:
			r.Version = version.DataVersionCancun
		case *prague.NewPayloadRequest:
			r.Version = version.DataVersionPrague
		case *amsterdam.NewPayloadRequest:
			r.Version = version.DataVersionAmsterdam
		default:
			return fmt.Errorf("NewPayloadRequest: unsupported view type %T", view)
		}
	}

	if err := copyByName(view, r); err != nil {
		return err
	}

	r.populateVersion(r.Version)

	return nil
}

// ToVersioned converts r into a *spec.VersionedNewPayloadRequest.
func (r *NewPayloadRequest) ToVersioned() (*spec.VersionedNewPayloadRequest, error) {
	out := &spec.VersionedNewPayloadRequest{}
	if err := toVersioned(r.Version, r, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates r from src.
func (r *NewPayloadRequest) FromVersioned(src *spec.VersionedNewPayloadRequest) error {
	return fromVersioned(r, src)
}

// MarshalSSZ marshals the request under the active fork view.
func (r *NewPayloadRequest) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(r, ds, make([]byte, 0, sizeSSZDyn(r, ds)))
}

// MarshalSSZTo marshals the request under the active fork view, appending to buf.
func (r *NewPayloadRequest) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(r, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the request under the active fork view.
func (r *NewPayloadRequest) SizeSSZ() int {
	return sizeSSZDyn(r, globalDynSSZ())
}

// UnmarshalSSZ decodes the request under the active fork view. Version must
// be set first.
func (r *NewPayloadRequest) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(r, r, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (r *NewPayloadRequest) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(r)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (r *NewPayloadRequest) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(r, globalDynSSZ(), hh)
}
