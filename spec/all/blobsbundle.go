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
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// BlobsBundle is a fork-agnostic blobs bundle. The cancun V1 schema carries
// one proof per blob; osaka V2 (reused by amsterdam) carries
// CELLS_PER_EXT_BLOB cell proofs per blob. The Go field shape is identical
// across both — only the SSZ list bounds differ — so a single union covers
// them.
type BlobsBundle struct {
	Version version.DataVersion

	Commitments []cancun.KZGCommitment
	Proofs      []cancun.KZGProof
	Blobs       []cancun.Blob
}

func (b *BlobsBundle) viewType() (any, error) {
	switch b.Version {
	case version.DataVersionCancun, version.DataVersionPrague:
		return (*cancun.BlobsBundle)(nil), nil
	case version.DataVersionOsaka,
		version.DataVersionAmsterdam,
		version.DataVersionBogota:
		return (*osaka.BlobsBundle)(nil), nil
	default:
		return nil, fmt.Errorf("BlobsBundle: unsupported version %d", b.Version)
	}
}

func (b *BlobsBundle) populateVersion(v version.DataVersion) {
	b.Version = v
}

// ToView returns a fresh fork-specific BlobsBundle populated from b.
func (b *BlobsBundle) ToView() (any, error) {
	return toViewByCopy(b)
}

// FromView populates b from a fork-specific BlobsBundle.
func (b *BlobsBundle) FromView(view any) error {
	if b.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *cancun.BlobsBundle:
			b.Version = version.DataVersionCancun
		case *osaka.BlobsBundle:
			b.Version = version.DataVersionOsaka
		default:
			return fmt.Errorf("BlobsBundle: unsupported view type %T", view)
		}
	}

	return copyByName(view, b)
}

// ToVersioned converts b into a *spec.VersionedBlobsBundle.
func (b *BlobsBundle) ToVersioned() (*spec.VersionedBlobsBundle, error) {
	out := &spec.VersionedBlobsBundle{}
	if err := toVersioned(b.Version, b, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates b from src.
func (b *BlobsBundle) FromVersioned(src *spec.VersionedBlobsBundle) error {
	return fromVersioned(b, src)
}

// MarshalJSON implements json.Marshaler.
func (b *BlobsBundle) MarshalJSON() ([]byte, error) {
	view, err := b.ToView()
	if err != nil {
		return nil, err
	}

	return json.Marshal(view)
}

// UnmarshalJSON implements json.Unmarshaler. Version must be set first.
func (b *BlobsBundle) UnmarshalJSON(input []byte) error {
	inst, err := newViewInstance(b)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(input, inst); err != nil {
		return err
	}

	return copyByName(inst, b)
}

// MarshalSSZ marshals the bundle under the active fork view.
func (b *BlobsBundle) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(b, ds, make([]byte, 0, sizeSSZDyn(b, ds)))
}

// MarshalSSZTo marshals the bundle under the active fork view, appending to buf.
func (b *BlobsBundle) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(b, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the bundle under the active fork view.
func (b *BlobsBundle) SizeSSZ() int {
	return sizeSSZDyn(b, globalDynSSZ())
}

// UnmarshalSSZ decodes the bundle under the active fork view. Version must
// be set first.
func (b *BlobsBundle) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(b, b, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (b *BlobsBundle) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(b)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (b *BlobsBundle) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(b, globalDynSSZ(), hh)
}

// String returns the JSON representation.
func (b *BlobsBundle) String() string {
	out, err := b.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
