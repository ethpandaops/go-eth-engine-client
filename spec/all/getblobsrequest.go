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
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetBlobsRequest is a fork-agnostic all-or-nothing engine_getBlobs request
// (cancun V1 / osaka V2). The request shape is identical across both forks.
type GetBlobsRequest struct {
	Version version.DataVersion

	BlobVersionedHashes []paris.Hash32
}

func (r *GetBlobsRequest) viewType() (any, error) {
	switch r.Version {
	case version.DataVersionCancun:
		return (*cancun.GetBlobsRequest)(nil), nil
	case version.DataVersionOsaka:
		return (*osaka.GetBlobsRequest)(nil), nil
	default:
		return nil, fmt.Errorf("GetBlobsRequest: unsupported version %d", r.Version)
	}
}

func (r *GetBlobsRequest) populateVersion(v version.DataVersion) {
	r.Version = v
}

// ToView returns a fresh fork-specific GetBlobsRequest populated from r.
func (r *GetBlobsRequest) ToView() (any, error) {
	return toViewByCopy(r)
}

// FromView populates r from a fork-specific GetBlobsRequest.
func (r *GetBlobsRequest) FromView(view any) error {
	if r.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *cancun.GetBlobsRequest:
			r.Version = version.DataVersionCancun
		case *osaka.GetBlobsRequest:
			r.Version = version.DataVersionOsaka
		default:
			return fmt.Errorf("GetBlobsRequest: unsupported view type %T", view)
		}
	}

	return copyByName(view, r)
}

// ToVersioned converts r into a *spec.VersionedGetBlobsRequest.
func (r *GetBlobsRequest) ToVersioned() (*spec.VersionedGetBlobsRequest, error) {
	out := &spec.VersionedGetBlobsRequest{}
	if err := toVersioned(r.Version, r, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates r from src.
func (r *GetBlobsRequest) FromVersioned(src *spec.VersionedGetBlobsRequest) error {
	return fromVersioned(r, src)
}

// MarshalSSZ marshals the request under the active fork view.
func (r *GetBlobsRequest) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(r, ds, make([]byte, 0, sizeSSZDyn(r, ds)))
}

// MarshalSSZTo marshals the request under the active fork view, appending to buf.
func (r *GetBlobsRequest) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(r, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the request under the active fork view.
func (r *GetBlobsRequest) SizeSSZ() int {
	return sizeSSZDyn(r, globalDynSSZ())
}

// UnmarshalSSZ decodes the request under the active fork view. Version must
// be set first.
func (r *GetBlobsRequest) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(r, r, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (r *GetBlobsRequest) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(r)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (r *GetBlobsRequest) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(r, globalDynSSZ(), hh)
}
