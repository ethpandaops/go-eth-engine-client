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

	"github.com/holiman/uint256"
	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/prague"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// GetPayloadResponse is a fork-agnostic engine_getPayload response. The
// shanghai V2 response carries only the payload and its value; cancun adds
// the blobs bundle and shouldOverrideBuilder; prague adds executionRequests.
// The nested ExecutionPayload and BlobsBundle are fork-agnostic unions.
//
// Paris (engine_getPayloadV1) returns a bare ExecutionPayload with no
// response container and so has no fork entry here.
type GetPayloadResponse struct {
	Version version.DataVersion

	ExecutionPayload      *ExecutionPayload
	BlockValue            *uint256.Int
	BlobsBundle           *BlobsBundle              // cancun+
	ShouldOverrideBuilder bool                      // cancun+
	ExecutionRequests     []prague.ExecutionRequest // prague+
}

func (r *GetPayloadResponse) viewType() (any, error) {
	switch r.Version {
	case version.DataVersionShanghai:
		return (*shanghai.GetPayloadResponse)(nil), nil
	case version.DataVersionCancun:
		return (*cancun.GetPayloadResponse)(nil), nil
	case version.DataVersionPrague:
		return (*prague.GetPayloadResponse)(nil), nil
	case version.DataVersionOsaka:
		return (*osaka.GetPayloadResponse)(nil), nil
	case version.DataVersionAmsterdam,
		version.DataVersionBogota:
		// Bogota reuses Amsterdam's V6 response (EIP-7805 does not add
		// fields to the build result).
		return (*amsterdam.GetPayloadResponse)(nil), nil
	default:
		return nil, fmt.Errorf("GetPayloadResponse: unsupported version %d", r.Version)
	}
}

func (r *GetPayloadResponse) populateVersion(v version.DataVersion) {
	r.Version = v

	if r.ExecutionPayload != nil {
		r.ExecutionPayload.populateVersion(v)
	}

	if r.BlobsBundle != nil {
		r.BlobsBundle.populateVersion(v)
	}
}

// ToView returns a fresh fork-specific GetPayloadResponse from r.
func (r *GetPayloadResponse) ToView() (any, error) {
	return toViewByCopy(r)
}

// FromView populates r from a fork-specific GetPayloadResponse.
func (r *GetPayloadResponse) FromView(view any) error {
	if r.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *shanghai.GetPayloadResponse:
			r.Version = version.DataVersionShanghai
		case *cancun.GetPayloadResponse:
			r.Version = version.DataVersionCancun
		case *prague.GetPayloadResponse:
			r.Version = version.DataVersionPrague
		case *osaka.GetPayloadResponse:
			r.Version = version.DataVersionOsaka
		case *amsterdam.GetPayloadResponse:
			r.Version = version.DataVersionAmsterdam
		default:
			return fmt.Errorf("GetPayloadResponse: unsupported view type %T", view)
		}
	}

	if err := copyByName(view, r); err != nil {
		return err
	}

	r.populateVersion(r.Version)

	return nil
}

// ToVersioned converts r into a *spec.VersionedGetPayloadResponse.
func (r *GetPayloadResponse) ToVersioned() (*spec.VersionedGetPayloadResponse, error) {
	out := &spec.VersionedGetPayloadResponse{}
	if err := toVersioned(r.Version, r, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates r from src.
func (r *GetPayloadResponse) FromVersioned(src *spec.VersionedGetPayloadResponse) error {
	return fromVersioned(r, src)
}

// MarshalSSZ marshals the response under the active fork view.
func (r *GetPayloadResponse) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(r, ds, make([]byte, 0, sizeSSZDyn(r, ds)))
}

// MarshalSSZTo marshals the response under the active fork view, appending to buf.
func (r *GetPayloadResponse) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(r, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the response under the active fork view.
func (r *GetPayloadResponse) SizeSSZ() int {
	return sizeSSZDyn(r, globalDynSSZ())
}

// UnmarshalSSZ decodes the response under the active fork view. Version must
// be set first.
func (r *GetPayloadResponse) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(r, r, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (r *GetPayloadResponse) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(r)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (r *GetPayloadResponse) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(r, globalDynSSZ(), hh)
}
