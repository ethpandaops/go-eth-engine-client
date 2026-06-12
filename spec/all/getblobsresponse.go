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
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetBlobsResponse is a fork-agnostic all-or-nothing engine_getBlobs
// response (cancun V1 / osaka V2). Entries are fork-agnostic BlobAndProof
// unions; a nil entry denotes a blob the EL did not have.
type GetBlobsResponse struct {
	Version version.DataVersion

	BlobsAndProofs []*BlobAndProof
}

func (r *GetBlobsResponse) viewType() (any, error) {
	switch r.Version {
	case version.DataVersionCancun:
		return (*cancun.GetBlobsResponse)(nil), nil
	case version.DataVersionOsaka:
		return (*osaka.GetBlobsResponse)(nil), nil
	default:
		return nil, fmt.Errorf("GetBlobsResponse: unsupported version %d", r.Version)
	}
}

func (r *GetBlobsResponse) populateVersion(v version.DataVersion) {
	r.Version = v

	for _, bp := range r.BlobsAndProofs {
		if bp != nil {
			bp.populateVersion(v)
		}
	}
}

// ToView returns a fresh fork-specific GetBlobsResponse populated from r.
func (r *GetBlobsResponse) ToView() (any, error) {
	return toViewByCopy(r)
}

// FromView populates r from a fork-specific GetBlobsResponse.
func (r *GetBlobsResponse) FromView(view any) error {
	if r.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *cancun.GetBlobsResponse:
			r.Version = version.DataVersionCancun
		case *osaka.GetBlobsResponse:
			r.Version = version.DataVersionOsaka
		default:
			return fmt.Errorf("GetBlobsResponse: unsupported view type %T", view)
		}
	}

	if err := copyByName(view, r); err != nil {
		return err
	}

	r.populateVersion(r.Version)

	return nil
}

// ToVersioned converts r into a *spec.VersionedGetBlobsResponse.
func (r *GetBlobsResponse) ToVersioned() (*spec.VersionedGetBlobsResponse, error) {
	out := &spec.VersionedGetBlobsResponse{}
	if err := toVersioned(r.Version, r, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates r from src.
func (r *GetBlobsResponse) FromVersioned(src *spec.VersionedGetBlobsResponse) error {
	return fromVersioned(r, src)
}

// MarshalSSZ marshals the response under the active fork view.
func (r *GetBlobsResponse) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(r, ds, make([]byte, 0, sizeSSZDyn(r, ds)))
}

// MarshalSSZTo marshals the response under the active fork view, appending to buf.
func (r *GetBlobsResponse) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(r, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the response under the active fork view.
func (r *GetBlobsResponse) SizeSSZ() int {
	return sizeSSZDyn(r, globalDynSSZ())
}

// UnmarshalSSZ decodes the response under the active fork view. Version must
// be set first.
func (r *GetBlobsResponse) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(r, r, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (r *GetBlobsResponse) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(r)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (r *GetBlobsResponse) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(r, globalDynSSZ(), hh)
}
